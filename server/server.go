package server

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lncapital/torq/pkg/lndutil"
	"google.golang.org/grpc"
)

// Start runs the background server. It subscribes to events, gossip and
// fetches data as needed and stores it in the database.
// It is meant to run as a background task / daemon and is the bases for all
// of Torqs data collection
func Start(conn *grpc.ClientConn, db *sqlx.DB) error {

	router := routerrpc.NewRouterClient(conn)
	client := lnrpc.NewLightningClient(conn)

	err := lndutil.SubscribeAndStoreHtlcEvents(router, db)
	if err != nil {
		return fmt.Errorf("in Start -> SubscribeAndStoreHtlcEvents(): %v", err)
	}

	err = lndutil.SubscribeAndStoreChannelEvents(client, db)
	if err != nil {
		return fmt.Errorf("in Start -> SubscribeAndStoreChannelEvents(): %v", err)
	}

	return nil
}

// Fetch static channel state and store it.
