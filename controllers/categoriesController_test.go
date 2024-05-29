package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_store "github.com/peternabil/go-api/mocks"
	"github.com/peternabil/go-api/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestListCategories(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	testCases := []struct {
		name          string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			buildStubs: func(store *mock_store.MockStore) {
				categories := []models.Category{{
					ID:          uuid.New(),
					Name:        "Category A",
					Description: "A test Category",
					UserID:      user.UID,
				}}
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategories(gomock.Any()).Times(1).Return(categories, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "error in categories",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategories(gomock.Any()).Times(1).Return(nil, errors.New("error in categories"))
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

			reader := strings.NewReader("")

			request, err := http.NewRequest("GET", "/smart-account/api/v1/category?page=1&page_size=10", reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestGetCategory(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	categoryId := uuid.New()
	wrongCategoryId := uuid.New()
	badCategoryId := "bad-id"
	category := models.Category{
		ID:          categoryId,
		Name:        "Category A",
		Description: "A test Category",
		UserID:      user.UID,
	}
	testCases := []struct {
		name          string
		params        string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(category, nil)
			},
			params: categoryId.String(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "does not exist",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(category, errors.New("error in categories"))
			},
			params: wrongCategoryId.String(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "bad uuid",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
			},
			params: badCategoryId,
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

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/category/%s", tt.params), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestCreateCategory(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	category := models.Category{
		Name:        "New Category",
		Description: "A new test category",
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
				"name":        category.Name,
				"description": category.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					CreateCategory(gomock.Any()).Times(1).Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name: "db error",
			body: gin.H{
				"name":        category.Name,
				"description": category.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					CreateCategory(gomock.Any()).Times(1).Return(category, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 400, recorder.Code)
			},
		},
		{
			name: "invalid input",
			body: gin.H{
				"description": category.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
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

			request, err := http.NewRequest("POST", "/smart-account/api/v1/category", bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestEditCategory(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	categoryId := uuid.New()
	updatedCategory := models.Category{
		Name:        "Updated Category",
		Description: "An updated test category",
	}
	testCases := []struct {
		name          string
		params        string
		body          gin.H
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			params: categoryId.String(),
			body: gin.H{
				"name":        updatedCategory.Name,
				"description": updatedCategory.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(updatedCategory, nil)
				store.EXPECT().
					EditCategory(gomock.Any()).Times(1).Return(updatedCategory, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "db error in get category",
			params: categoryId.String(),
			body: gin.H{
				"name":        updatedCategory.Name,
				"description": updatedCategory.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(updatedCategory, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "db error in edit category",
			params: categoryId.String(),
			body: gin.H{
				"name":        updatedCategory.Name,
				"description": updatedCategory.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(updatedCategory, nil)
				store.EXPECT().
					EditCategory(gomock.Any()).Times(1).Return(updatedCategory, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "invalid input",
			params: categoryId.String(),
			body: gin.H{
				"description": updatedCategory.Description,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
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

			request, err := http.NewRequest("PUT", fmt.Sprintf("/smart-account/api/v1/category/%s", tt.params), bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	categoryId := uuid.New()
	testCases := []struct {
		name          string
		params        string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			params: categoryId.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					DeleteCategory(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name:   "not found",
			params: categoryId.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					DeleteCategory(gomock.Any()).Times(1).Return(errors.New("category not found"))
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

			request, err := http.NewRequest("DELETE", fmt.Sprintf("/smart-account/api/v1/category/%s", tt.params), nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
