package lndutil

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"time"
)

type dbForwardEvent struct {

	// The microseconds' version of TimestampNs, used by TimescaleDB
	Time time.Time `db:"time"`

	// The number of nanoseconds elapsed since January 1, 1970 UTC when this
	// circuit was completed.
	TimeNs uint64 `db:"time_ns"`

	// The incoming channel ID that carried the HTLC that created the circuit.
	ChanIdIn uint64 `db:"incoming_channel_id"`

	// The outgoing channel ID that carried the preimage that completed the
	// circuit.
	ChanIdOut uint64 `db:"outgoing_channel_id"`

	// The total fee (in milli-satoshis) that this payment circuit carried.
	FeeMsat uint64 `db:"fee_msat"`

	// The total amount (in milli-satoshis) of the incoming HTLC that created
	// half the circuit.
	AmtInMsat uint64 `db:"incoming_amount_msat"`

	// The total amount (in milli-satoshis) of the outgoing HTLC that created
	// the second half of the circuit.
	AmtOutMsat uint64 `db:"outgoing_amount_msat"`
}

// storeForwardingHistory
func storeForwardingHistory(db *sqlx.DB, fwh []*lnrpc.ForwardingEvent) error {

	if len(fwh) > 0 {
		tx := db.MustBegin()

		for _, event := range fwh {

			dbEvent := dbForwardEvent{
				Time:       time.Unix(0, int64(event.TimestampNs)).Round(time.Microsecond),
				TimeNs:     event.TimestampNs,
				ChanIdIn:   event.ChanIdIn,
				ChanIdOut:  event.ChanIdOut,
				FeeMsat:    event.FeeMsat,
				AmtInMsat:  event.AmtInMsat,
				AmtOutMsat: event.AmtOutMsat,
			}

			sql := `
			INSERT INTO forward(time, time_ns, fee_msat,
				incoming_channel_id, outgoing_channel_id,
				incoming_amount_msat, outgoing_amount_msat)
			VALUES (:time, :time_ns, :fee_msat,
				:incoming_channel_id, :outgoing_channel_id,
				:incoming_amount_msat, :outgoing_amount_msat);`

			if _, err := tx.NamedExec(sql, dbEvent); err != nil {
				return err
			}
		}
		tx.Commit()
	}

	return nil
}

// MAXEVENTS is used to set the maximum events in ForwardingHistoryRequest.
// It's also used to check if we need to request more.
const MAXEVENTS int = 50000

// fetchLastForwardTime fetches the latest recorded forward, if none is set already.
// This should only run once when a server starts.
func fetchLastForwardTime(db *sqlx.DB) (uint64, error) {

	var lastNs uint64

	row := db.QueryRow("SELECT time_ns FROM forward ORDER BY time_ns DESC LIMIT 1;")
	err := row.Scan(&lastNs)
	if err == sql.ErrNoRows {
		return 0, nil
	}
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

// SubscribeForwardingEvents repeatedly requests forwarding history starting after the last
// forwarding stored in the database and stores new forwards.
func SubscribeForwardingEvents(client lnrpc.LightningClient, db *sqlx.DB) error {

	// Request the forwarding history at the requested interval.
	// NB!: This timer is slowly being shifted because of the time required to
	//fetch and store the response.
	for range time.Tick(30 * time.Second) {

		// Fetch the nanosecond timestamp of the most recent record we have.
		lastNs, err := fetchLastForwardTime(db)
		if err != nil {
			return err
		}

		// Keep fetching until LND returns less than the max number of records requested.
		for {
			fwh, err := fetchForwardingHistory(lastNs, client)
			if err != nil {
				return err
			}

			// Store the forwarding history
			err = storeForwardingHistory(db, fwh.ForwardingEvents)
			if err != nil {
				return err
			}

			// Stop fetching if there are fewer forwards than max requested
			// (indicates that we have the last forwarding record)
			if len(fwh.ForwardingEvents) < MAXEVENTS {
				break
			}
		}
	}

	return nil
}
