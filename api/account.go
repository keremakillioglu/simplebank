package api

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"

	db "github.com/keremakillioglu/simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// balance is zero initially
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

// in gin, everything we do includes a context object
// handlerFunc  involves *Context input
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	// if client provided invalid data
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// param1:httpcode :send 400 bad request, param2: httpobject to send to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//if no errors
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	// return created account in db & error
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violationn":
				// status forbidden (code 403), error message
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return

			}

		}

		// internal error (code 500), error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// successfully created the account
	ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	ID      int64 `uri:"id" binding:"required,min=1"`
	Balance int64 `json:"balance"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	// if client provided invalid data
	if err := ctx.ShouldBindUri(&req); err != nil {
		// param1:httpcode :send 400 bad request, param2: httpobject to send to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// internal error (code 500), error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	// if client provided invalid data
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// limit= pagesize, offset= number of records that db should skip
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := server.store.ListAccounts(ctx, arg)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// internal error (code 500), error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}
