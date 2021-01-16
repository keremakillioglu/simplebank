package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	mockdb "github.com/keremakillioglu/simplebank/db/mock"

	"github.com/golang/mock/gomock"

	"github.com/keremakillioglu/simplebank/util"

	db "github.com/keremakillioglu/simplebank/db/sqlc"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	// for table- driven test set to cover all posible scenarios of the GetAccount API
	// store test cases in an anonymous struct
	testCases := []struct {
		name      string
		accountID int64
		// build stubs
		buildStubs func(store *mockdb.MockStore)
		//check the output of the idea
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:      "OK",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			// build stubs
			// gomock expects get account function to be called with any context and have an id equal to accountid created above
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			//check response
			require.Equal(t, http.StatusOK, recorder.Code)
			//check body
			requireBodyMatchAccount(t, recorder.Body, account)
		},
	}, {
		name:      "NOTFOUND",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			// build stubs
			// gomock expects get account function to be called with any context and have an id equal to accountid created above
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			//check response
			require.Equal(t, http.StatusNotFound, recorder.Code)
			// since there is no such id, status notfound will be returned
		},
	},
		{
			name:      "INTERNALERROR",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				// gomock expects get account function to be called with any context and have an id equal to accountid created above
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// since there is no such id, status notfound will be returned
			},
		},

		{
			name:      "INVALIDID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				// gomock expects get account function to be called with any context and have an id equal to accountid created above
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0) //dont call the function because id is invalid
				// dont return anything
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// since there is no such id, status notfound will be returned
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//buildstubs with mockstore
			tc.buildStubs(store)

			//start test server and send a request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// url path of the api we want to call
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			// method, url, body
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// send the request via server.router and record response at recorder
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

}

// checking response body: body is just a bytes buffer and account to compare
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// read data from body
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	// get account from body
	var accountFromResponse db.Account
	err = json.Unmarshal(data, &accountFromResponse)

	require.NoError(t, err)
	require.Equal(t, account, accountFromResponse)
}
