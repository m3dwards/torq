package torqsrv

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lncapital/torq/torqrpc"
	"time"
)

func (s torqGrpc) GetChannelFlow(ctx context.Context, cfr *torqrpc.ChannelFlowRequest) (*torqrpc.ChannelFlow, error) {

	qry := `
		select
			coalesce(ROUND(sum(incoming_amount_msat/1000) FILTER (WHERE incoming_channel_id in (?))), 0) as AmtIn,
			coalesce(ROUND(sum(fee_msat/1000) FILTER (WHERE incoming_channel_id in (?))), 0)  as FeeIn,
			coalesce(count(time_ns) FILTER (WHERE incoming_channel_id in (?)), 0) as CountIn,
			coalesce(ROUND(sum(outgoing_amount_msat/1000) FILTER (WHERE outgoing_channel_id in (?))), 0) as AmtOut,
			coalesce(ROUND(sum(fee_msat/1000) FILTER (WHERE outgoing_channel_id in (?))), 0)  as FeeOut,
			coalesce(count(time_ns) FILTER (WHERE outgoing_channel_id in (?)), 0) as CountOut
		from forward
		where (time >= ? and time <= ?);
		`
	// Query the forward table for aggregated forwarding history
	q, args, err := sqlx.In(qry,
		cfr.ChanIds, cfr.ChanIds, cfr.ChanIds, cfr.ChanIds, cfr.ChanIds, cfr.ChanIds,
		time.Unix(cfr.FromTime, 0), time.Unix(cfr.ToTime, 0))
	if err != nil {
		return nil, errors.Wrap(err, "In")
	}

	q = s.db.Rebind(q)
	cf := torqrpc.ChannelFlow{}
	err = s.db.QueryRowx(q, args...).Scan(
		&cf.AmtIn,
		&cf.FeeIn,
		&cf.CountIn,
		&cf.AmtOut,
		&cf.FeeOut,
		&cf.CountOut,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	// Return the filtering parameters
	cf.ChanIds = cfr.ChanIds
	cf.FromTime = cfr.FromTime
	cf.ToTime = cfr.ToTime

	return &cf, nil
}

func (s torqGrpc) GetForwards(context.Context, *torqrpc.ForwardsRequest) (*torqrpc.
	Forwards, error) {

	return &torqrpc.Forwards{}, nil
}
