package gapi

import (
	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/pb"
	"github.com/definitely-unique-username/simple_bank/token"
	"github.com/definitely-unique-username/simple_bank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config *util.Config, store db.Store) *Server {
	server := &Server{
		store:      store,
		tokenMaker: token.NewPasetoMaker(config.SymmetricalKey),
		config:     *config,
	}

	return server
}
