package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_store "github.com/peternabil/go-api/mocks"
	"github.com/peternabil/go-api/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGetDailyValues(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	spendings := []models.Spending{{
		Date:     time.Now(),
		Total:    100,
		Negative: false,
	}}
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "success",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransactionsDateRangeGroupByDay(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(spendings, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "wrong negative param",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=test",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "wrong end date format",
			param: "?end_date=2023-12-19 16:23-25&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "wrong end date format",
			param: "?start_date=2023-12-19 16:23-25&end_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "db error",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				spendings := []models.Spending{{
					Date:     time.Now(),
					Total:    100,
					Negative: false,
				}}
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransactionsDateRangeGroupByDay(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(spendings, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "error reading token",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, errors.New("error reading token"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStore := mock_store.NewMockStore(mockCtrl)
			tt.buildStubs(mockStore)

			server, _ := NewServer(mockStore, nil)
			recorder := httptest.NewRecorder()

			reader := strings.NewReader("")

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/daily%s", tt.param), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestGetHighestCategory(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	spendings := []models.SpendingCategory{{
		Date:       time.Now(),
		Total:      100,
		Negative:   false,
		Cname:      "test category",
		CategoryID: uuid.New(),
	}}
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "success",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetHighestSpendingCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(spendings, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "wrong negative param",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=test",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "wrong end date format",
			param: "?end_date=2023-12-19 16:23-25&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "wrong end date format",
			param: "?start_date=2023-12-19 16:23-25&end_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "db error",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetHighestSpendingCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(spendings, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStore := mock_store.NewMockStore(mockCtrl)
			tt.buildStubs(mockStore)

			server, _ := NewServer(mockStore, nil)
			recorder := httptest.NewRecorder()

			reader := strings.NewReader("")

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/highest-cat%s", tt.param), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestGetHighestPriority(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	spendings := []models.SpendingPriority{{
		Date:       time.Now(),
		Total:      100,
		Negative:   false,
		Pname:      "test priority",
		PriorityID: uuid.New(),
	}}
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "success",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetHighestSpendingPriority(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(spendings, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "wrong negative param",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=test",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "wrong end date format",
			param: "?end_date=2023-12-19 16:23-25&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "wrong end date format",
			param: "?start_date=2023-12-19 16:23-25&end_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "db error",
			param: "?end_date=2023-12-19T16:23:25.742Z&start_date=2023-11-19T16:21:53.561Z&negative=true",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetHighestSpendingPriority(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(spendings, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStore := mock_store.NewMockStore(mockCtrl)
			tt.buildStubs(mockStore)

			server, _ := NewServer(mockStore, nil)
			recorder := httptest.NewRecorder()

			reader := strings.NewReader("")

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/highest-prio%s", tt.param), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
