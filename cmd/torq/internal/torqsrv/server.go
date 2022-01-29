package torqsrv

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/lncapital/torq/torqrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
)

//type Server struct {
//	//lncrpc.GetNodeChannelsRequest
//	//lncrpc.GetNodeChannelsResponse
//	//lncrpc.UnimplementedLncrpcServer
//	LndConn grpc.ClientConnInterface
//}

type torqGrpc struct {
	certPath string
	keyPath  string
	host     string
	port     string
	wport    string
	srv      *grpc.Server
	wsrv     *grpcweb.WrappedGrpcServer
	torqrpc.UnimplementedTorqrpcServer
}

func (s torqGrpc) GetForwards(context.Context, *torqrpc.ForwardsRequest) (*torqrpc.
	ForwardResponse, error) {
	//client := lnrpc.NewLightningClient(s.LndConn)
	//ncl, err := state.CreateNodeChannelList(client)
	//
	//if err != nil {
	//	return nil, fmt.Errorf("error CreateNodeChannelList: %v", err)
	//}
	//
	//r := lncrpc.GetNodeChannelsResponse{
	//	Nodes: ncl,
	//}

	return &torqrpc.ForwardResponse{}, nil
}

func NewServer(host, port, wport, cert, key string) (torqGrpc, error) {

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		return torqGrpc{}, errors.Wrapf(err, "->NewServerTLSFromFile(%s, %s)", cert, key)
	}

	opts := grpc.Creds(creds)

	s := grpc.NewServer(opts)

	srv := torqGrpc{
		certPath: cert,
		keyPath:  key,
		host:     host,
		port:     port,
		wport:    wport,
		srv:      s,
		wsrv:     grpcweb.WrapServer(s),
	}

	torqrpc.RegisterTorqrpcServer(srv.srv, &srv)

	return srv, nil
}

func (s *torqGrpc) StartWeb() error {
	// TODO: Replace with log

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		resp.Header().Set("Access-Control-Expose-Headers", "grpc-status, grpc-message")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
			" Accept-Encoding, X-CSRF-Token, XMLHttpRequest, x-user-agent, x-grpc-web, "+
			"grpc-status, grpc-message")

		if s.wsrv.IsGrpcWebRequest(req) || s.wsrv.IsAcceptableGrpcCorsRequest(req) {
			s.wsrv.ServeHTTP(resp, req)
		}
	})

	fmt.Printf("gRPC proxy server listening on: %s:%s \n", s.host, s.wport)
	dns := fmt.Sprintf("%s:%d", s.host, s.wport)
	err := http.ListenAndServeTLS(dns, s.certPath, s.keyPath, nil)

	if err != nil {
		return errors.Wrapf(err, "ListenAndServeTLS(%s, %s, %s, %v)", dns, s.certPath, s.keyPath,
			nil)
	}
	return nil
}

func (s *torqGrpc) StartGrpc() error {

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return errors.Wrapf(err, "net.Listen(\"tcp\", %s", fmt.Sprintf("%s:%d", s.host, s.port))
	}

	// TODO: Replace with log
	fmt.Printf("gRPC server listening on: %s:%s \n", s.host, s.port)
	err = s.srv.Serve(lis)
	if err != nil {
		return errors.Wrapf(err, "srv.Serve(%v)", lis)
	}

	return nil
}
