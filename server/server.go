package server

import (
	"context"
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
func Start(conn *grpc.ClientConn, db *sqlx.DB) error {

	router := routerrpc.NewRouterClient(conn)
	client := lnrpc.NewLightningClient(conn)

	// Create an error group to catch errors from go routines.
	// TODO: Improve this by using the context to propogate the error,
	//   shutting down the if one of the subscribe go routines fail.
	//   https://www.fullstory.com/blog/why-errgroup-withcontext-in-golang-server-handlers/
	// TODO: Also consider using the same context used by the gRPC connection from Golang and the
	//   gRPC server of Torq
	ctx := context.Background()
	errs, ctx := errgroup.WithContext(ctx)

	// HTLC events
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreHtlcEvents(router, db)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeAndStoreHtlcEvents(%v, %v)", router, db)
		}
		return nil
	})

	// Channel Events
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreChannelEvents(client, db)
		if err != nil {
			return errors.Wrapf(err, "Start->SubscribeAndStoreChannelEvents(%v, %v)", router, db)
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
