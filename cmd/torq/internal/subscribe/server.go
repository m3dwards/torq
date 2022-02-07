package subscribe

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lncapital/torq/pkg/lndutil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Start runs the background server. It subscribes to events, gossip and
// fetches data as needed and stores it in the database.
// It is meant to run as a background task / daemon and is the bases for all
// of Torqs data collection
func Start(ctx context.Context, conn *grpc.ClientConn, db *sqlx.DB) error {

	router := routerrpc.NewRouterClient(conn)
	client := lnrpc.NewLightningClient(conn)

	// Create an error group to catch errors from go routines.
	// TODO: Improve this by using the context to propogate the error,
	//   shutting down the if one of the subscribe go routines fail.
	//   https://www.fullstory.com/blog/why-errgroup-withcontext-in-golang-server-handlers/
	// TODO: Also consider using the same context used by the gRPC connection from Golang and the
	//   gRPC server of Torq
	errs, ctx := errgroup.WithContext(ctx)

	// Import Open channels
	err := lndutil.ImportChannelList(lnrpc.ChannelEventUpdate_OPEN_CHANNEL, db, client)
	if err != nil {
		fmt.Println(err)
		return errors.Wrapf(err, "Start -> importChannelList(%s, %v, %v)",
			lnrpc.ChannelEventUpdate_OPEN_CHANNEL, db, client)
	}

	// Import Closed channels
	err = lndutil.ImportChannelList(lnrpc.ChannelEventUpdate_CLOSED_CHANNEL, db, client)
	if err != nil {
		fmt.Println(err)
		return errors.Wrapf(err, "Start -> importChannelList(%s, %v, %v)",
			lnrpc.ChannelEventUpdate_CLOSED_CHANNEL, db, client)
	}

	// Import Node info (based on channels)
	err = lndutil.ImportMissingNodeEvents(client, db)
	if err != nil {
		fmt.Println(err)
		return errors.Wrapf(err, "Start -> ImportMissingNodeEvents(%v, %v)", client, db)
	}

	// Get the public key of our node
	ni, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return errors.Wrapf(err, "start -> client.GetNodeInfo(ctx, &lnrpc.NodeInfoRequest{})")
	}
	lndutil.InitOurNodesList([]string{ni.IdentityPubkey})

	pubKeyChan := make(chan string)
	chanPointChan := make(chan string)

	// Initialize the peer list
	err = lndutil.InitPeerList(db)
	if err != nil {
		return errors.Wrapf(err, "start -> InitPeerList(%v)", db)
	}

	// Initialize the channel id list
	err = lndutil.InitChanIdList(db)
	if err != nil {
		return errors.Wrapf(err, "start -> InitChanIdList(%v)", db)
	}

	// Start listening for updates
	go lndutil.UpdatePeerList(pubKeyChan)
	go lndutil.UpdateChanIdList(chanPointChan)

	// Transactions
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreTransactions(ctx, client, db)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeAndStoreTransactions(%v, %v, %v)", ctx, client, db)
		}
		return nil
	})

	// Graph (Node updates, fee updates etc.)
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreChannelGraph(ctx, client, db)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeAndStoreChannelGraph(%v, %v, %v)", ctx, client, db)
		}
		return nil
	})

	// HTLC events
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreHtlcEvents(ctx, router, db)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeAndStoreHtlcEvents(%v, %v, %v)", ctx, router, db)
		}
		return nil
	})

	// Channel Events
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreChannelEvents(ctx, client, db, pubKeyChan, chanPointChan)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeAndStoreChannelEvents(%v, %v, %v)", ctx, router, db)
		}
		return nil
	})

	// Forwarding history
	errs.Go(func() error {

		err := lndutil.SubscribeForwardingEvents(ctx, client, db, nil)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeForwardingEvents(%v, %v, %v, %v)", ctx,
				client, db, nil)
		}

		return nil
	})

	return errs.Wait()
}

// Fetch static channel state and store it.
