package api

import (
	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/token"
	"github.com/definitely-unique-username/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config *util.Config, store db.Store) *Server {
	router := gin.Default()
	server := &Server{
		store:      store,
		router:     router,
		tokenMaker: token.NewPasetoMaker(config.SymmetricalKey),
		config:     *config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	accountsGroup := router.Group("/accounts", authMiddleware(server.tokenMaker))
	accountsGroup.GET("/", server.getAccounts)
	accountsGroup.POST("/", server.createAccount)
	accountsGroup.GET("/:id", server.getAccount)

	usersGroup := router.Group("/users")
	usersGroup.POST("/", server.createUser)
	usersGroup.POST("/login", server.loginUser)

	transfersGroup := router.Group("/transfers", authMiddleware(server.tokenMaker))
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
