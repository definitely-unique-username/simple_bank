package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/definitely-unique-username/simple_bank/api"
	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/gapi"
	"github.com/definitely-unique-username/simple_bank/pb"
	"github.com/definitely-unique-username/simple_bank/util"
)

func main() {
	config, err := util.LoadConfig("./", ".env")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGatewayServer(config util.Config, store db.Store) {
	server := gapi.NewServer(config, store)

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.HTTPSevrverAddress)

	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("HTTP gateway server started at %s", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server := gapi.NewServer(config, store)
	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("gRPC server started at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start grpc server:")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server := api.NewServer(&config, store)

	if err := server.Start(config.HTTPSevrverAddress); err != nil {
		log.Fatal("Cannot start server", err)
	}
}
