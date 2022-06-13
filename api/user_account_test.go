package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"
	mockdb "github.com/kerok-kristoffer/backendStub/db/mock"
	db "github.com/kerok-kristoffer/backendStub/db/sqlc"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserAccountAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userId        int64
		buildStubs    func(userAccount *mockdb.MockUserAccount)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID,
			buildStubs: func(account *mockdb.MockUserAccount) {
				account.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userId: user.ID,
			buildStubs: func(account *mockdb.MockUserAccount) {
				account.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userId: user.ID,
			buildStubs: func(account *mockdb.MockUserAccount) {
				account.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidId",
			userId: 0,
			buildStubs: func(account *mockdb.MockUserAccount) {
				account.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// TODO: add more cases
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			controller.Finish()

			account := mockdb.NewMockUserAccount(controller)
			tc.buildStubs(account)

			server := NewServer(account)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d", tc.userId)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)

}

func randomUser() db.User {
	f := faker.New()
	return db.User{
		ID:       f.Int64(),
		FullName: f.Person().Name(),
		Hash:     f.Hash().MD5(),
	}

}
