package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
)

type Server struct {
	userAccount db.UserAccount
	router      *gin.Engine
}

func NewServer(userAccount db.UserAccount) *Server {
	server := &Server{userAccount: userAccount}
	router := gin.Default()

	// adds validator "currency" to api calls according tut #14, not actually used at the moment.
	// unit tests for api from #14 are on tut maker's github
	// could add validator for units (g, l, oz, etc) at a later point, so saving as template
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil
		}
	}

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUserAccount)
	router.GET("/users", server.listUsers)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
