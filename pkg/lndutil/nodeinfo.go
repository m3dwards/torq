package lndutil

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// getMissingNodePubKeys creates a string slice with all the PubKey of all nodes
// we have a channel with, but where we do not have any node event records.
func getMissingNodePubKeys(db *sqlx.DB) ([]string, error) {

	// Fetch a list of all unknown channel PubKeys by subtracting a set of all
	// pub_keys from all channels, from a set of pub_keys from all nodes already stored.
	q := `select
		   CASE
			   WHEN existing_nodes is null THEN all_nodes
			   ELSE array_subtract(all_nodes, existing_nodes)
			END missing from
		(select array_agg(distinct pub_key) as all_nodes from channel_event where event_type in (0,1)) as c,
		(select array_agg(distinct pub_key) as existing_nodes from node_event as ce) as e;`

	var res []string
	err := db.QueryRowx(q).Scan(pq.Array(&res))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ImportMissingNodeEvents imports information about all nodes that we have had a channel with.
func ImportMissingNodeEvents(client lnrpc.LightningClient, db *sqlx.DB) error {

	pubKeyList, err := getMissingNodePubKeys(db)
	if err != nil {
		return err
	}

	ctx := context.Background()
	for _, p := range pubKeyList {
		rsp, err := client.GetNodeInfo(ctx, &lnrpc.NodeInfoRequest{PubKey: p, IncludeChannels: false})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.NotFound:
					continue
				default:
					return errors.Wrapf(err, "failed to get alias for node with pubkey %s", p)
				}
			}
		}
		ts := time.Now()
		err = InsertNodeEvent(db, ts, rsp.Node.PubKey, rsp.Node.Alias, rsp.Node.Color,
			rsp.Node.Addresses, rsp.Node.Features)
		if err != nil {
			return err
		}
	}

	return nil
}
