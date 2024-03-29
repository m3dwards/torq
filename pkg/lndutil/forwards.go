package lndutil

import (
	"context"
	"database/sql"
	"github.com/benbjohnson/clock"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
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

func convMicro(ns uint64) time.Time {
	return time.Unix(0, int64(ns)).Round(time.Microsecond).UTC()
}

const querySfwh = `INSERT INTO forward(time, time_ns, fee_msat,
		incoming_channel_id, outgoing_channel_id,
		incoming_amount_msat, outgoing_amount_msat)
	VALUES ($1, $2, $3,$4, $5,$6, $7)
	ON CONFLICT (time, time_ns) DO NOTHING;`

// storeForwardingHistory
func storeForwardingHistory(db *sqlx.DB, fwh []*lnrpc.ForwardingEvent) error {

	if len(fwh) > 0 {
		tx := db.MustBegin()

		for _, event := range fwh {

			if _, err := tx.Exec(querySfwh, convMicro(event.TimestampNs), event.TimestampNs,
				event.FeeMsat, event.ChanIdIn, event.ChanIdOut, event.AmtInMsat,
				event.AmtOutMsat); err != nil {
				return errors.Wrapf(err, "storeForwardingHistory->tx.Exec(%v)",
					querySfwh)
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

	switch err {
	case nil:
		return lastNs, errors.Wrapf(err, "fetchLastForwardTime->row.Scan(%+v)", &lastNs)
	case sql.ErrNoRows:
		return 0, nil
	}

	return lastNs, nil
}

type lightningClientForwardingHistory interface {
	ForwardingHistory(ctx context.Context, in *lnrpc.ForwardingHistoryRequest,
		opts ...grpc.CallOption) (*lnrpc.ForwardingHistoryResponse, error)
}

// fetchForwardingHistory fetches the forwarding history from LND.
func fetchForwardingHistory(ctx context.Context, client lightningClientForwardingHistory,
	lastTimestamp uint64,
	maxEvents int) (
	*lnrpc.ForwardingHistoryResponse, error) {

	fwhReq := &lnrpc.ForwardingHistoryRequest{
		StartTime:    lastTimestamp,
		NumMaxEvents: uint32(maxEvents),
	}
	fwh, err := client.ForwardingHistory(ctx, fwhReq)
	if err != nil {
		return nil, errors.Wrapf(err, "fetchForwardingHistory->ForwardingHistory(%v, %v)", ctx,
			fwhReq)
	}

	return fwh, nil
}

// FwhOptions allows the caller to adjust the number of forwarding events can be requested at a time
// and set a custom time interval between requests.
type FwhOptions struct {
	MaxEvents *int
	Tick      <-chan time.Time
}

// SubscribeForwardingEvents repeatedly requests forwarding history starting after the last
// forwarding stored in the database and stores new forwards.
func SubscribeForwardingEvents(ctx context.Context, client lightningClientForwardingHistory,
	db *sqlx.DB, opt *FwhOptions) error {

	me := MAXEVENTS

	// Check if maxEvents has been set and that it is bellow the hard coded maximum defined by
	// the constant MAXEVENTS.
	if (opt != nil) && ((*opt.MaxEvents > MAXEVENTS) || (*opt.MaxEvents <= 0)) {
		me = *opt.MaxEvents
	}

	// Create the default ticker used to fetch forwards at a set interval
	c := clock.New()
	ticker := c.Tick(10 * time.Second)

	// If a custom ticker is set in the options, override the default ticker.
	if (opt != nil) && (opt.Tick != nil) {
		ticker = opt.Tick
	}
	// Request the forwarding history at the requested interval.
	// NB!: This timer is slowly being shifted because of the time required to
	//fetch and store the response.
	for {
		// Exit if canceled
		select {
		case <-ctx.Done():
			return nil
		case <-ticker:

			// Fetch the nanosecond timestamp of the most recent record we have.
			lastNs, err := fetchLastForwardTime(db)
			lastTimestamp := lastNs / uint64(time.Second)
			if err != nil {
				return errors.Wrapf(err, "SubscribeForwardingEvents->fetchLastForwardTime(%v)", db)
			}

			// Keep fetching until LND returns less than the max number of records requested.
		fetchAll:
			for {
				fwh, err := fetchForwardingHistory(ctx, client, lastTimestamp, me)
				if err != nil {
					return errors.Wrapf(err, "SubscribeForwardingEvents->fetchForwardingHistory(%v, "+
						"%v, %v, %v"+
						")", ctx, client, lastTimestamp, me)
				}

				// Store the forwarding history
				err = storeForwardingHistory(db, fwh.ForwardingEvents)
				if err != nil {
					return errors.Wrapf(err, "SubscribeForwardingEvents->storeForwardingHistory(%v, "+
						"%v)", db, fwh.ForwardingEvents)
				}

				// Stop fetching if there are fewer forwards than max requested
				// (indicates that we have the last forwarding record)
				if len(fwh.ForwardingEvents) < me {
					break fetchAll
				}

			}
		}
	}

	return nil
}
