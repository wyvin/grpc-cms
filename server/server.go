package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"grpc-content/pkg/util"
	pb "grpc-content/proto"
)

var (
	Port        string
	CertName    string
	CertPemPath string
	KeyPemPath  string
	EndPoint    string
)

func Run() (err error) {
	EndPoint = ":" + Port
	tlsConfig := util.GetTLSConfig(CertPemPath, KeyPemPath)

	conn, err := net.Listen("tcp", EndPoint)
	if err != nil {
		log.Printf("TCP Listen err:%v\n", err)
	}
	srv := newServer(conn, tlsConfig)

	log.Printf("gRPC and https listen on: %s\n", Port)

	if err = srv.Serve(tls.NewListener(conn, tlsConfig)); err != nil {
		log.Printf("ListenAndServe: %v\n", err)
	}

	return err
}

func newServer(conn net.Listener, tlsConfig *tls.Config) *http.Server {
	var opts []grpc.ServerOption

	// grpc server
	creds, err := credentials.NewServerTLSFromFile(CertPemPath, KeyPemPath)
	if err != nil {
		log.Printf("Failed to create server TLS credentials %v", err)
	}

	opts = append(opts, grpc.Creds(creds))
	grpcServer := grpc.NewServer(opts...)

	// register grpc pb
	pb.RegisterHelloWorldServer(grpcServer, NewHelloService())
	pb.RegisterAdServer(grpcServer, NewAdService())

	// gw server
	ctx := context.Background()
	dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertName)
	if err != nil {
		log.Printf("Failed to create client TLS credentials %v", err)
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := runtime.NewServeMux()
	// gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	// register grpc-gateway pb
	if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
		log.Printf("Failed to register gw hello server: %v\n", err)
	}
	if err := pb.RegisterAdHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
		log.Printf("Failed to register gw ad server: %v\n", err)
	}

	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	return &http.Server{
		Addr:      EndPoint,
		Handler:   util.GrpcHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}
}
