package main

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
)

func main() {

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "raspibolt",
		Pass:         "turnkey-format-KEG-geodetic",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the transaction
	var h = chainhash.Hash{}
	chainhash.Decode(&h, "d1698baee9d848f2584729df65a2b474776f1d8900ad7f46e522e03a2818331c")
	result, err := client.GetRawTransactionVerbose(&h)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("\nTime: %s \ntxId: %s", time.Unix(result.Blocktime, 0).UTC(), result.Txid)

}
