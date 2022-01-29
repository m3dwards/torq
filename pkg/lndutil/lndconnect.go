package lndutil

import (
	"fmt"
	"io/ioutil"

	"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

// Connect connects to LND using gRPC.
func Connect(host, tlsCertPath, macaroonPath string) (*grpc.ClientConn, error) {

	tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		return nil, fmt.Errorf("cannot get lnd tls credentials: %v", err)
	}

	macaroonBytes, err := ioutil.ReadFile(macaroonPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read macaroon file: %v", err)
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		return nil, fmt.Errorf("cannot unmarshal macaroon: %v", err)
	}

	macCred, err := macaroons.NewMacaroonCredential(mac)
	if err != nil {
		return nil, fmt.Errorf("cannot create macaroon credentials: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(macCred),
	}

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot dial to lnd: %v", err)
	}

	return conn, nil
}
