package api

import (
	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	router := gin.Default()
	server := &Server{store: store, router: router}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	accountsGroup := router.Group("/accounts")
	accountsGroup.GET("/", server.getAccounts)
	accountsGroup.POST("/", server.createAccount)
	accountsGroup.GET("/:id", server.getAccount)

	usersGroup := router.Group("/users")
	usersGroup.POST("/", server.createUser)

	transfersGroup := router.Group("/transfers")
	transfersGroup.POST("/", server.createTransfer)

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
