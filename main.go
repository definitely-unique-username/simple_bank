package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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

	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server := gapi.NewServer(&config, store)
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
