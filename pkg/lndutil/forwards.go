package lndutil

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"time"
)

// storeForwardingHistory
//func storeForwardingHistory(db *sqlx.DB, fwh []*lnrpc.ForwardingEvent) error {
//
//	if len(fwh) > 0 {
//		tx := db.MustBegin()
//		sql := "INSERT INTO forwards(time, ) VALUES (:id, " +
//			":company_id, :mobileuser_id)"
//		if _, err := tx.NamedExec(sql); err != nil {
//			fmt.Printf("Failed updating sales transactions for mobileuser %s \n", m.ID)
//			log.Println(err)
//		}
//		tx.Commit()
//	}
//
//	jb, err := json.Marshal(h)
//	if err != nil {
//		return fmt.Errorf("storeForwardingHistory -> json.Marshal(%v): %v", h, err)
//	}
//
//	stm := `INSERT INTO htlc_event (time, event_type, outgoing_channel_id, incoming_channel_id,
//		event) VALUES($1, $2, $3, $4, $5)`
//
//	timestampMs := time.Unix(0, int64(h.TimestampNs)).Round(time.Microsecond)
//	_, err = db.Exec(stm, timestampMs, h.EventType, h.OutgoingChannelId, h.IncomingChannelId, jb)
//	if err != nil {
//		return fmt.Errorf(`storeForwardingHistory -> db.Exec(%s, %v, %v, %v, %v, %v): %v`,
//			stm, timestampMs, h.EventType, h.OutgoingChannelId, h.IncomingChannelId, jb, err)
//	}
//
//	return nil
//}

// MAXEVENTS is used to set the maximum events in ForwardingHistoryRequest.
// It's also used to check if we need to request more.
const MAXEVENTS int = 50000

// fetchLastForwardTime fetches the latest recorded forward, if none is set already.
// This should only run once when a server starts.
func fetchLastForwardTime(db *sqlx.DB) (uint64, error) {

	var lastNs uint64

	row := db.QueryRow("SELECT time_ns FROM forward ORDER BY time_ns LIMIT 1")
	err := row.Scan(&lastNs)
	if err != nil {
		return lastNs, fmt.Errorf("fetchLastForwardTime(): %v", err)
	}

	return lastNs, nil
}

// fetchForwardingHistory fetches the forwarding history from LND.
func fetchForwardingHistory(lastNs uint64, client lnrpc.LightningClient) (
	*lnrpc.ForwardingHistoryResponse, error) {

	ctx := context.Background()
	fwh, err := client.ForwardingHistory(ctx, &lnrpc.ForwardingHistoryRequest{StartTime: lastNs,
		NumMaxEvents: uint32(MAXEVENTS)})
	if err != nil {
		return nil, fmt.Errorf("fetchForwardingHistory -> ForwardingHistory(): %v", err)
	}

	return fwh, nil
}

// SubscribeForwardingUpdates repeatedly requests forwarding history starting after the last
// forwarding stored in the database and stores new forwards.
func SubscribeForwardingUpdates(client lnrpc.LightningClient, db *sqlx.DB) error {

	// Fetch the nanosecond timestamp of the most recent forward we already have.
	lastNs, err := fetchLastForwardTime(db)
	if err != nil {
		return err
	}

	// Request the forwarding history at the requested interval.
	// NB!: This timer is slowly being shifted because of the time required to fetch and store the
	// response.
	for range time.Tick(10 * time.Second) {

		// Keep fetching until LND returns less than the max number of records requested.
		for {
			fwh, err := fetchForwardingHistory(lastNs, client)
			if err != nil {
				return err
			}

			// Store the forwarding history

			// Stop fetching if there are fewer forwards than max requested
			// (indicates that we have the last forwarding record)
			if len(fwh.ForwardingEvents) < MAXEVENTS {
				break
			}
		}
	}

	//fmt.Println("last offset: ", fwh.LastOffsetIndex)
	//fmt.Println("last offset: ", fwh.ForwardingEvents[0].TimestampNs)
	//fmt.Println("last offset: ", fwh.ForwardingEvents[len(fwh.ForwardingEvents)-1].TimestampNs)

	//for _, fw := range fwh.ForwardingEvents {
	//
	//}

	//err = storeForwardingHistory(db, fwh)
	//if err != nil {
	//	return fmt.Errorf("storeForwardingHistory(): %v", err)
	//}

	return nil
}
