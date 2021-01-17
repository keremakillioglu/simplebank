package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// register the custom validator with gin
	// binding.Validator.Engine() returns a general interface -> convert to validator pointer
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// p1: name of validation tag, p2 : validCyrrency func in validator.go
		v.RegisterValidation("currency", validCurrency)
		// binding:"... currency at account.go createAccountRequest param & transfer.go"
	}

	// route, handlerfunc
	//if we pass multiple parameters: route,middlewares, handlefunc
	router.POST("/accounts", server.createAccount)

	// id is a parameter provided by URI
	router.GET("/accounts/:id", server.getAccount)

	// parameters will be retrieved from querystring
	router.GET("/accounts", server.listAccount)

	// transfer details specified in req body
	router.POST("/transfers", server.createTransfer)

	//if we pass multiple parameters: route,middlewares, handlefunc
	router.POST("/newuser", server.createUser)

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
