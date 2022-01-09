package lndutil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"io"
	"time"
)

// storeChannelEvent extracts the timestamp, channel ID and PubKey from the
// ChannelEvent and converts the original struct to json.
// Then it's stored in the database in the channel_event table.
func storeChannelEvent(db *sqlx.DB, ce *lnrpc.ChannelEventUpdate) error {

	jb, err := json.Marshal(ce)
	if err != nil {
		return fmt.Errorf("storeChannelEvent -> json.Marshal(%v): %v", ce, err)
	}

	timestampMs := time.Now()

	var ChanID uint64
	var ChannelPoint string
	var PubKey string

	switch ce.Type {
	case lnrpc.ChannelEventUpdate_OPEN_CHANNEL:
		c := ce.GetOpenChannel()
		ChanID = c.ChanId
		ChannelPoint = c.ChannelPoint
		PubKey = c.RemotePubkey
	case lnrpc.ChannelEventUpdate_CLOSED_CHANNEL:
		c := ce.GetClosedChannel()
		ChanID = c.ChanId
		ChannelPoint = c.ChannelPoint
		PubKey = c.RemotePubkey
	case lnrpc.ChannelEventUpdate_FULLY_RESOLVED_CHANNEL:
		c := ce.GetFullyResolvedChannel()
		ChannelPoint = fmt.Sprintf("%b:%d", c.GetFundingTxidBytes(), c.GetOutputIndex())
	case lnrpc.ChannelEventUpdate_ACTIVE_CHANNEL:
		c := ce.GetActiveChannel()
		ChannelPoint = fmt.Sprintf("%b:%d", c.GetFundingTxidBytes(), c.GetOutputIndex())
	case lnrpc.ChannelEventUpdate_INACTIVE_CHANNEL:
		c := ce.GetInactiveChannel()
		ChannelPoint = fmt.Sprintf("%b:%d", c.GetFundingTxidBytes(), c.GetOutputIndex())
	case lnrpc.ChannelEventUpdate_PENDING_OPEN_CHANNEL:
		c := ce.GetPendingOpenChannel()
		ChannelPoint = fmt.Sprintf("%b:%d", c.GetTxid(), c.GetOutputIndex())
	default:
		// TODO: Need to improve error handling and logging in the case of unknown event.
		//  Simply storing the event without any link to a channel.
	}

	stm := `INSERT INTO channel_event (time, event_type, chan_id, chan_point, pub_key, 
event) VALUES($1, $2, $3, $4, $5, $6)`

	_, err = db.Exec(stm, timestampMs, ce.Type.String(), ChanID, ChannelPoint, PubKey, jb)
	if err != nil {
		return fmt.Errorf(`storeChannelEvent -> db.Exec(%s, %v, %s, %v, %s, %s, %v): %v`,
			stm, timestampMs, ce.Type.String(), ChanID, ChannelPoint, PubKey, jb, err)
	}

	return nil
}

// SubscribeAndStoreChannelEvents Subscribes to channel events from LND and stores them in the
// database as a time series
func SubscribeAndStoreChannelEvents(client lnrpc.LightningClient, db *sqlx.DB) error {

	ctx := context.Background()
	stream, err := client.SubscribeChannelEvents(ctx, &lnrpc.ChannelEventSubscription{})
	if err != nil {
		return fmt.Errorf("SubscribeAndStoreChannelEvents -> SubscribeChannelEvents(): %v", err)
	}

	for {

		chanEvent, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("%v.ListFeatures(_) = _, %v", stream, err)
		}

		err = storeChannelEvent(db, chanEvent)
		if err != nil {
			return fmt.Errorf("StreamHTLC(): %v", err)
		}

	}

	return nil
}
