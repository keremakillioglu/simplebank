package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute and run db queries in transactions
// note that indivdiual queries are handled with Queries struct in db.go
// however, each query makes an operation on a specific table; so queries struct doesnt support transactions
// here, we are extending its functionality by embedding it to Store struct
type Store struct {
	// composition is used, rather than inheritance; as its suggested in GO
	// all individual queries and functions will be available to Store
	// and we will be able to implement transactions by adding more functions
	*Queries
	//  sql.DB object is used because its required to create a new db transaction
	db *sql.DB
}

//NewStore creates a new store
func NewStore(db *sql.DB) *Store {

	return &Store{
		db:      db,
		Queries: New(db), //defined in db.go by sqlc
	}
}

// take a context and callback function as an input, start a new db transaction
// create a new Queries object with that transaction, and call the callback function on the queries
// finally commit or rollback the transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// txKey is used as a context key for db transactions
// use it to get txName
// key with type empty struct, empty bracket means empty object with that type
// var txKey = struct{}{}

//TransferTx performs a money transfer from one account to another
// It creates a transfer record, add account entries and update account balances within sinle tx
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// move money out of account1
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// move money into account2
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return err
	})
	return result, err

}
