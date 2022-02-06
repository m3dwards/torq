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
		select coalesce(o.chan_id_out, i.chan_id_in) as chan_id,
			   coalesce(amount_out, 0) as amt_out,
			   coalesce(fee_out, 0) as fee_out,
			   coalesce(o.count_out, 0) as count_out,
			   coalesce(amount_in, 0) as amt_in,
			   coalesce(fee_in, 0) as fee_in,
			   coalesce(i.count_in, 0) as count_in
		from (
			select outgoing_channel_id as chan_id_out,
				   floor(sum(outgoing_amount_msat/1000)) as amount_out,
				   floor(sum(fee_msat/1000)) as fee_out,
				   count(time_ns) as count_out
			from forward
			where outgoing_channel_id in ('792773172513013761', '775745035938955265')
			and (time >= '2022-01-15 00:00:00' and time <= '2022-01-31 23:59:59.999999')
			group by outgoing_channel_id) as o
		full join
			 (select incoming_channel_id as chan_id_in,
				   floor(sum(incoming_amount_msat/1000)) as amount_in,
				   floor(sum(fee_msat/1000)) as fee_in,
				   count(time_ns) as count_in
			from forward
			where incoming_channel_id in ('792773172513013761', '775745035938955265')
		    and (time >= '2022-01-15 00:00:00' and time <= '2022-01-31 23:59:59.999999')
			group by incoming_channel_id) as i
		on o.chan_id_out = i.chan_id_in;
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
