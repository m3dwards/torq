package lndutil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"io"
)

// storeHTLCEvent extracts the timestamp and channel IDs from the HtlcEvent and converts the original struct to json.
// Then it's stored in the database in the htlc table.
func storeHTLCEvent(db *sqlx.DB, h *routerrpc.HtlcEvent) error {

	jb, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("storeHTLCEvent -> json.Marshal(%v): %v", h, err)
	}

	stm := `INSERT INTO htlc (time, outgoing_channel_id, incoming_channel_id, event) VALUES($1)`

	timestampMs := h.TimestampNs / 1000.0
	_, err = db.Exec(stm, timestampMs, h.OutgoingChannelId, h.IncomingChannelId, jb)
	if err != nil {
		return fmt.Errorf(`storeHTLCEvent -> db.Exec(%s, %v, %v, %v, %v): %v`,
			stm, timestampMs, h.OutgoingChannelId, h.IncomingChannelId, jb, err)
	}

	return nil
}

// SubscribeAndStoreHtlcEvents subscribes to HTLC events from LND and stores them in the database as time series.
// NB: LND has marked HTLC event streaming as experimental. Delivery is not guaranteed, so dataset might not be complete
// HTLC events is primarily used to diagnose how good a channel / node is. And if the channel allocation should change.
func SubscribeAndStoreHtlcEvents(router routerrpc.RouterClient, db *sqlx.DB) error {

	ctx := context.Background()
	htlcStream, err := router.SubscribeHtlcEvents(ctx, &routerrpc.SubscribeHtlcEventsRequest{})
	if err != nil {
		return fmt.Errorf("SubscribeAndStoreHtlcEvents -> SubscribeHtlcEvents(): %v", err)
	}

	if err != nil {
		return fmt.Errorf("SubscribeAndStoreHtlcEvents -> createChanIdMap(client): %v", err)
	}

	for {

		htlcEvent, err := htlcStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("%v.ListFeatures(_) = _, %v", htlcStream, err)
		}

		err = storeHTLCEvent(db, htlcEvent)
		if err != nil {
			return fmt.Errorf("StreamHTLC(): %v", err)
		}

	}

	return nil
}
