package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
)

type Server struct {
	userAccount *db.UserAccount
	router      *gin.Engine
}

func NewServer(userAccount *db.UserAccount) *Server {
	server := &Server{userAccount: userAccount}
	router := gin.Default()

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
