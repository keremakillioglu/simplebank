package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/keremakillioglu/simplebank/db/sqlc"
)

// balance is zero initially
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// in gin, everything we do includes a context object
// handlerFunc  involves *Context input
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	// if client provided invalid data
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// param1:httpcode :send 400 bad request, param2: httpobject to send to client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	//if no errors
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// return transaction & error
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		// internal error (code 500), error message
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//compare whether sender and receiver have the same type of currency

	ctx.JSON(http.StatusOK, result)

}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {

	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		//if no such account exists
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
