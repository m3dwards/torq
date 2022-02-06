package lndutil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lightningnetwork/lnd/lnrpc"
	"io"
	"time"
)

func InsertNodeEvent(db *sqlx.DB, ts time.Time, pubKey string, alias string, color string,
	na []*lnrpc.NodeAddress, f map[uint32]*lnrpc.Feature) error {

	// Create json byte object from node address map
	najb, err := json.Marshal(na)
	if err != nil {
		return errors.Wrapf(err, "insertNodeEvent -> json.Marshal(%v)", na)
	}

	// Create json byte object from features list
	fjb, err := json.Marshal(f)
	if err != nil {
		return errors.Wrapf(err, "insertNodeEvent -> json.Marshal(%v)", f)
	}

	db.Exec(`INSERT INTO node_event (timestamp, pub_key, alias, color, node_addresses, features) 
    	VALUES ($1, $2, $3, $4, $5, $6)`, ts, pubKey, alias, color, najb, fjb)

	return nil
}

//
func storeChanGraphUpdate(db *sqlx.DB, update *lnrpc.GraphTopologyUpdate) error {
	//fmt.Println("----- New update -----")

	// Get a list of relevant nodes.

	for _, nu := range update.NodeUpdates {
		fmt.Printf("\nNode Update:\n%v\n", nu)
		//err := insertNodeEvent(db, nu)
		//if err != nil {
		//	return err
		//}
	}

	// Get a list of relevant channels

	//for _, cu := range update.ChannelUpdates {
	//	fmt.Printf("\nChannel Update Update:\n%v\n", cu)
	//}

	return nil
}

// pubKeyList is used to store which node and channel updates to store. We only want to store
// updates that are relevant to our channels and their nodes.
var pubKeyList []string

// InitPeerList fetches all public keys from the list of all channels. This is used to
// filter out noise from the graph updates.
func InitPeerList(db *sqlx.DB) error {
	q := `select array_agg(distinct pub_key) as all_nodes from channel_event where event_type in (0, 1);`
	err := db.QueryRow(q).Scan(pq.Array(&pubKeyList))
	if err != nil {
		return errors.Wrapf(err, "InitPeerList -> db.QueryRow(%s).Scan(pq.Array(%v))", q, pubKeyList)
	}
	return nil
}

// UpdatePeerList is meant to run as a gorouting. It adds new public keys to the pubKeyList
// and removes existing ones.
func UpdatePeerList(c chan string) {

	var pubKey string
	var present bool

	for {
		// Wait for new peers to enter
		pubKey = <-c
		// Add it to the peer list, if not already present.
		for _, p := range pubKeyList {
			if p == pubKey {
				present = true
				break
			}
		}

		// If not present add it to the public key  list
		if present == false {
			pubKeyList = append(pubKeyList, pubKey)
		}

		// Reset to false in order to allow the next public key to be added.
		present = false
	}
}

// isRelevant is used to check if any public key is in the pubKeyList.
func isRelevant(pubKeys ...string) bool {

	for _, p := range pubKeyList {

		// Check if any of the provided public keys equals the current public key.
		for _, np := range pubKeys {
			if p == np {
				// If found, no reason to search further, immediatly return true.
				return true
			}
		}

	}

	return false
}

// SubscribeAndStoreChannelGraph Subscribes to channel updates
func SubscribeAndStoreChannelGraph(ctx context.Context, client lnrpc.LightningClient, db *sqlx.DB) error {

	req := lnrpc.GraphTopologySubscription{}
	stream, err := client.SubscribeChannelGraph(ctx, &req)

	if err != nil {
		return errors.Wrapf(err, "SubscribeAndStoreChannelGraph -> client.SubscribeChannelGraph(%v, %v)", ctx, req)
	}

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		gpu, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "SubscribeChannelEvents -> stream.Recv()")
		}

		//storeChanGraphUpdate(db, gpu)
		for _, nu := range gpu.NodeUpdates {
			if isRelevant(nu.IdentityKey) {
				ts := time.Now()
				err := InsertNodeEvent(db, ts, nu.IdentityKey, nu.Alias, nu.Color,
					nu.NodeAddresses, nu.Features)
				if err != nil {
					return err
				}
			}
		}

		//for _, cu := range gpu.ChannelUpdates {
		//	if isRelevant(cu.AdvertisingNode, cu.ConnectingNode) {
		//
		//		cu.RoutingPolicy
		//		ts := time.Now()
		//		err := InsertNodeEvent(db, ts, nu.IdentityKey, nu.Alias, nu.Color,
		//			nu.NodeAddresses, nu.Features)
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}

		//gpu.ChannelUpdates

		//err = storeChannelEvent(db, gpu)
		//if err != nil {
		//	return errors.Wrapf(err, "storeChannelEvent(%v, %v)", db, client)
		//}

	}

	return nil
}
