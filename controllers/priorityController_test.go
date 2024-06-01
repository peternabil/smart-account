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

func TestListPriorities(t *testing.T) {
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
				priorities := []models.Priority{{
					ID:          uuid.New(),
					Name:        "Priority A",
					Description: "A test Priority",
					UserID:      user.UID,
				}}
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetPriorities(gomock.Any()).Times(1).Return(priorities, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "error in priorities",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetPriorities(gomock.Any()).Times(1).Return(nil, errors.New("error in categories"))
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

			request, err := http.NewRequest("GET", "/smart-account/api/v1/priority?page=1&page_size=10", reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestGetPriorty(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	priorityID := uuid.New()
	wrongpriorityID := uuid.New()
	badpriorityID := "bad-id"
	priority := models.Priority{
		ID:          priorityID,
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
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(priority, nil)
			},
			params: priorityID.String(),
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
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(priority, errors.New("error in categories"))
			},
			params: wrongpriorityID.String(),
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
			params: badpriorityID,
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

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/priority/%s", tt.params), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestCreatePriorty(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	priority := models.Priority{
		Name:        "New Priority",
		Description: "A new test Priority",
		Level:       5,
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
				"Name":        priority.Name,
				"Description": priority.Description,
				"Level":       priority.Level,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					CreatePriority(gomock.Any()).Times(1).Return(priority, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name: "db error",
			body: gin.H{
				"name":        priority.Name,
				"description": priority.Description,
				"Level":       priority.Level,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					CreatePriority(gomock.Any()).Times(1).Return(priority, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 400, recorder.Code)
			},
		},
		{
			name: "invalid input",
			body: gin.H{
				"description": priority.Description,
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

			request, err := http.NewRequest("POST", "/smart-account/api/v1/priority", bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestEditPriorty(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	priorityID := uuid.New()
	updatedPriority := models.Priority{
		Name:        "Updated Priority",
		Description: "An updated test Priority",
		Level:       7,
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
			params: priorityID.String(),
			body: gin.H{
				"name":        updatedPriority.Name,
				"description": updatedPriority.Description,
				"Level":       updatedPriority.Level,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(updatedPriority, nil)
				store.EXPECT().
					EditPriority(gomock.Any()).Times(1).Return(updatedPriority, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "db error in get category",
			params: priorityID.String(),
			body: gin.H{
				"name":        updatedPriority.Name,
				"description": updatedPriority.Description,
				"Level":       updatedPriority.Level,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(updatedPriority, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "db error in edit category",
			params: priorityID.String(),
			body: gin.H{
				"name":        updatedPriority.Name,
				"description": updatedPriority.Description,
				"Level":       updatedPriority.Level,
			},
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(updatedPriority, nil)
				store.EXPECT().
					EditPriority(gomock.Any()).Times(1).Return(updatedPriority, errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "invalid input",
			params: priorityID.String(),
			body: gin.H{
				"description": updatedPriority.Description,
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

			request, err := http.NewRequest("PUT", fmt.Sprintf("/smart-account/api/v1/priority/%s", tt.params), bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestDeletePriorty(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	priorityID := uuid.New()
	testCases := []struct {
		name          string
		params        string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "success",
			params: priorityID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					DeletePriority(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name:   "not found",
			params: priorityID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					DeletePriority(gomock.Any()).Times(1).Return(errors.New("priority not found"))
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

			request, err := http.NewRequest("DELETE", fmt.Sprintf("/smart-account/api/v1/priority/%s", tt.params), nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
