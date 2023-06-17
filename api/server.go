package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kerok-kristoffer/backendStub/db/sqlc"
	"github.com/kerok-kristoffer/backendStub/token"
	"github.com/kerok-kristoffer/backendStub/util"
)

type Server struct {
	config      util.Config
	userAccount db.UserAccount
	ingredients db.Ingredient
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
	gin.SetMode(gin.ReleaseMode)
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew", server.renewAccessToken)

	userAuthRoutes := router.Group("/users").Use(authMiddleware(server.tokenMaker))
	userAuthRoutes.GET("/:id", server.getUserAccount)

	userAuthRoutes.GET("/sub", server.getSubscriptions)
	userAuthRoutes.POST("/sub", server.applySubscription)

	userAuthRoutes.POST("/functions", server.addIngredientFunction)
	userAuthRoutes.GET("/functions", server.listIngredientFunctions)
	userAuthRoutes.POST("/ingredients", server.addIngredient)
	userAuthRoutes.GET("/ingredients", server.listIngredients)
	userAuthRoutes.GET("/ingredients/count", server.getIngredientCount)
	userAuthRoutes.POST("/ingredients/:id", server.updateIngredient)
	userAuthRoutes.DELETE("ingredients/:id", server.deleteIngredient)

	userAuthRoutes.POST("/formulas", server.addFormula)
	userAuthRoutes.POST("/formulas/:id", server.updateFormula)
	userAuthRoutes.GET("/formulas", server.listFormulas)
	userAuthRoutes.DELETE("/formulas/:id", server.deleteFormula)

	// todo kerok - implement separate middleware for admins and add checks on subscription level
	adminRoutes := router.Group("/users").Use(authMiddleware(server.tokenMaker))
	adminRoutes.GET("/", server.listUsers)

	server.router = router
}

func corsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		config, err := util.LoadConfig("../")
		if err != nil {
			log.Fatalln("Cannot load config:", err)
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.JSON(204, nil)
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
