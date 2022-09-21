package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kerok-kristoffer/formulating/db/mock"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/token"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetIngredientsAPI(t *testing.T) {
	user, err := randomUser()

	require.NoError(t, err)

	ingredients := []db.Ingredient{
		{Name: "Wind"},
		{Name: "Fire"},
	}

	ingredientsRequest := listIngredientsRequest{
		PageId:   1,
		PageSize: 10,
	}

	ingredientsByUserParams := ingredientByUserParams(user, ingredientsRequest)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(account *mockdb.MockUserAccount, user db.User)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(account *mockdb.MockUserAccount, user db.User) {
				account.EXPECT().
					GetUserByUserName(gomock.Any(), user.UserName).
					Times(1).
					Return(user, nil)
				account.EXPECT().
					ListIngredientsByUserId(gomock.Any(), gomock.Eq(ingredientsByUserParams)).
					Times(1).
					Return(ingredients, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchIngredients(t, ingredients, recorder.Body)
			},
		}, {
			name: "Empty",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(account *mockdb.MockUserAccount, user db.User) {
				account.EXPECT().
					GetUserByUserName(gomock.Any(), user.UserName).
					Times(1).
					Return(user, nil)
				account.EXPECT().
					ListIngredientsByUserId(gomock.Any(), gomock.Eq(ingredientsByUserParams)).
					Times(1).
					Return([]db.Ingredient{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchIngredients(t, []db.Ingredient(nil), recorder.Body)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			account := mockdb.NewMockUserAccount(controller)
			tc.buildStubs(account, user)

			server := NewTestServer(t, account)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/ingredients?page_id=1&page_size=10")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func requireBodyMatchIngredients(t *testing.T, ingredients []db.Ingredient, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotIngredients []db.Ingredient

	err = json.Unmarshal(data, &gotIngredients)
	require.NoError(t, err)
	require.Equal(t, ingredients, gotIngredients)
}
