package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/keremakillioglu/simplebank/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	// not *dbStore instead bc store is an interface
	store  db.Store
	router *gin.Engine
}

// NewServer creates  a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// route, handlerfunc
	//if we pass multiple parameters: route,middlewares, handlefunc
	router.POST("/accounts", server.createAccount)

	// id is a parameter provided by URI
	router.GET("/accounts/:id", server.getAccount)

	// parameters will be retrieved from querystring
	router.GET("/accounts", server.listAccount)

	// add routes to router
	server.router = router
	return server
}

// we had to get access to store object to set new account to database

// Start runs the HTTP Server on a specific address
// take an address as an input and return an error
func (server *Server) Start(address string) error {
	// use the action provided by gin, sinc erouter part is private it cannot be accessed from outside of this api package
	return server.router.Run(address)
}

// returns key-value pair
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
