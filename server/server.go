package server

import (
	"context"
	"fmt"
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
	ctx := context.Background()
	errs, ctx := errgroup.WithContext(ctx)

	// HTLC events
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreHtlcEvents(router, db)
		if err != nil {
			return fmt.Errorf("in Start -> SubscribeAndStoreHtlcEvents(): %v", err)
		}
		return nil
	})

	// Channel Events
	errs.Go(func() error {
		err := lndutil.SubscribeAndStoreChannelEvents(client, db)
		if err != nil {
			return fmt.Errorf("in Start -> SubscribeAndStoreChannelEvents(): %v", err)
		}
		return nil
	})

	// Forwarding history
	errs.Go(func() error {

		err := lndutil.SubscribeForwardingEvents(client, db)
		if err != nil {
			return fmt.Errorf("in Start -> SubscribeForwardingEvents(): %v", err)
		}

		return nil
	})

	return errs.Wait()
}

// Fetch static channel state and store it.
