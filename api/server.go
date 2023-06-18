package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/nehalshaquib/my-bank/db/sqlc"
)

// Server serves http request for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new http server and setups routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//api routes
	router.GET("/", server.Hello)
	router.POST("/account", server.CreateAccount)
	router.GET("/accounts", server.ListAccounts)
	router.GET("/account/:id", server.GetAccount)
	router.DELETE("/account/:id", server.DeleteAccount)

	server.router = router

	return server
}

func (server *Server) Start() error {
	return server.router.Run("localhost:8080")
}

func (server *Server) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "This is mybank server")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
