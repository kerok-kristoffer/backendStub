package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/token"
	"github.com/kerok-kristoffer/formulating/util"
)

type Server struct {
	config      util.Config
	userAccount db.UserAccount
	router      *gin.Engine
	tokenMaker  token.Maker
}

func NewServer(config util.Config, userAccount db.UserAccount) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:      config,
		userAccount: userAccount,
		router:      nil,
		tokenMaker:  maker,
	}

	// adds validator "currency" to api calls according tut #14, not actually used at the moment.
	// unit tests for api from #14 are on tut maker's github
	// could add validator for units (g, l, oz, etc) at a later point, so saving as template
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	userAuthRoutes := router.Group("/users").Use(authMiddleware(server.tokenMaker))
	userAuthRoutes.GET("/:id", server.getUserAccount)
	// todo kerok - implement routes for ingredients/recipes

	// todo kerok - implement separate middleware for admins
	adminRoutes := router.Group("/users").Use(authMiddleware(server.tokenMaker))
	adminRoutes.GET("/", server.listUsers)

	server.router = router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
