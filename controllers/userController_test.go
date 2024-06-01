package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_store "github.com/peternabil/go-api/mocks"
	"github.com/peternabil/go-api/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserIndex(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	users := []models.User{
		{
			UID:       uuid.New(),
			Email:     "user1@test.com",
			FirstName: "Test",
			LastName:  "User1",
		},
		{
			UID:       uuid.New(),
			Email:     "user2@test.com",
			FirstName: "Test",
			LastName:  "User2",
		},
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUsers().Times(1).Return(users, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "internal server error",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUsers().Times(1).Return([]models.User{}, errors.New("error fetching users"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			request, err := http.NewRequest("GET", "/smart-account/api/v1/users", nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestUserFind(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user1@test.com",
		FirstName: "Test",
		LastName:  "User1",
		Password:  string(encryptedPass),
	}

	testCases := []struct {
		name          string
		param         string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "success",
			param: user.UID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUser(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "invalid uuid",
			param: "invalid-uuid",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUser(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "user not found",
			param: user.UID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUser(gomock.Any()).Times(1).Return(models.User{}, errors.New("user not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/users/%s", tt.param), nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestSignUp(t *testing.T) {
	user := models.User{
		Email:     "user1@gmail.com",
		FirstName: "Test",
		LastName:  "User1",
		Password:  "Password@123!",
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			body: gin.H{
				"Email":     user.Email,
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"Password":  user.Password,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					SignUp(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "server error",
			body: gin.H{
				"Email":     user.Email,
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"Password":  user.Password,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					SignUp(gomock.Any()).Times(1).Return(user, errors.New("server error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "wrong params",
			body: gin.H{
				"Email":     "u@gmail.com",
				"FirstName": "Te",
				"Password":  "Pass",
			},
			buildStubs: func(store *mock_store.MockStore) {
				// store.EXPECT().
				// 	SignUp(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "fake email",
			body: gin.H{
				"Email":     "peter@test.com",
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"Password":  user.Password,
			},
			buildStubs: func(store *mock_store.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "invalid email",
			body: gin.H{
				"Email":     "peter.com",
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"Password":  user.Password,
			},
			buildStubs: func(store *mock_store.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "password validation error",
			body: gin.H{
				"Email":     user.Email,
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"Password":  "easypasswordtoguess",
			},
			buildStubs: func(store *mock_store.MockStore) {
				// store.EXPECT().
				// SignUp(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/smart-account/api/auth/signup", bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestLogin(t *testing.T) {
	password := "Password123!"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		Email:    "user1@gmail.com",
		Password: string(encryptedPass),
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			body: gin.H{
				"Email":    user.Email,
				"Password": password,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					FindUser(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					CreateToken(gomock.Any()).Times(1).Return("token", nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "wrong input",
			body: gin.H{
				"Password": password,
			},
			buildStubs: func(store *mock_store.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "incorrect password",
			body: gin.H{
				"Email":    user.Email,
				"Password": "wrong-password",
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					FindUser(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "user not found",
			body: gin.H{
				"Email":    "nonexistent@gmail.com",
				"Password": password,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					FindUser(gomock.Any()).Times(1).Return(models.User{}, errors.New("user not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "error in token create",
			body: gin.H{
				"Email":    user.Email,
				"Password": password,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					FindUser(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					CreateToken(gomock.Any()).Times(1).Return("token", errors.New("error in token creation"))
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

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/smart-account/api/auth/login", bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
