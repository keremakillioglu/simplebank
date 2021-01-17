package api

import (
	"net/http"
	"time"

	"github.com/keremakillioglu/simplebank/util"

	"github.com/lib/pq"

	db "github.com/keremakillioglu/simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// in gin, everything we do includes a context object
// handlerFunc  involves *Context input
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	// if client provided invalid data
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// param1:httpcode :send 400 bad request, param2: httpobject to send to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//if no errors
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// return created account in db & error
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violationn":
				// status forbidden (code 403), error message
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		// internal error (code 500), error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	// successfully created the account
	ctx.JSON(http.StatusOK, response)
}
