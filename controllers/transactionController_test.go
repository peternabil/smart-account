package controllers

import (
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

func TestListTransactions(t *testing.T) {
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
				category := models.Category{
					ID:          uuid.New(),
					Name:        "Category A",
					Description: "A test Category",
					UserID:      user.UID,
				}
				priority := models.Priority{
					ID:          uuid.New(),
					Name:        "Priority A",
					Description: "A test Priority",
					Level:       5,
					UserID:      user.UID,
				}
				transactions := []models.Transaction{{
					ID:          uuid.New(),
					Title:       "transaction title",
					CategoryID:  category.ID,
					Category:    category,
					Priority:    priority,
					Amount:      100,
					Negative:    true,
					Description: "a test transaction",
					PriorityID:  priority.ID,
					UserID:      user.UID,
				}}
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransactions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(transactions, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "error in transactions",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransactions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error in transactions"))
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

			request, err := http.NewRequest("GET", "/smart-account/api/v1/transaction?page=1&page_size=10", reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestGetTransaction(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}
	transactionId := uuid.New()
	wrongTransactionId := uuid.New()
	badTransactionId := "bad-id"
	category := models.Category{
		ID:          uuid.New(),
		Name:        "Category A",
		Description: "A test Category",
		UserID:      user.UID,
	}
	priority := models.Priority{
		ID:          uuid.New(),
		Name:        "Priority A",
		Description: "A test Priority",
		Level:       5,
		UserID:      user.UID,
	}
	transaction := models.Transaction{
		ID:          transactionId,
		Title:       "transaction title",
		CategoryID:  category.ID,
		Category:    category,
		Priority:    priority,
		Amount:      100,
		Negative:    true,
		Description: "a test transaction",
		PriorityID:  priority.ID,
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
					GetTransaction(gomock.Any()).Times(1).Return(transaction, nil)
			},
			params: transactionId.String(),
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
					GetTransaction(gomock.Any()).Times(1).Return(transaction, errors.New("error in transactions"))
			},
			params: wrongTransactionId.String(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "bad uuid",
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			params: badTransactionId,
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

			request, err := http.NewRequest("GET", fmt.Sprintf("/smart-account/api/v1/transaction/%s", tt.params), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestTransactionCreate(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}

	transactionID := uuid.New()
	categoryID := uuid.New()
	priorityID := uuid.New()
	validBody := map[string]interface{}{
		"Title":       "test transaction",
		"Category":    categoryID.String(),
		"Amount":      100,
		"Negative":    false,
		"Description": "test description",
		"Priority":    priorityID.String(),
	}
	inValidBody := map[string]interface{}{
		"Title": 1,
	}

	category := models.Category{
		ID:          categoryID,
		Name:        "Category A",
		Description: "A test Category",
		UserID:      user.UID,
	}

	priority := models.Priority{
		ID:          priorityID,
		Name:        "Priority A",
		Description: "A test Priority",
		Level:       5,
		UserID:      user.UID,
	}

	transaction := models.Transaction{
		ID:          transactionID,
		Title:       validBody["Title"].(string),
		CategoryID:  category.ID,
		PriorityID:  priority.ID,
		Amount:      validBody["Amount"].(int),
		Negative:    validBody["Negative"].(bool),
		Description: validBody["Description"].(string),
		UserID:      user.UID,
	}

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: validBody,
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(category, nil)
				store.EXPECT().
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(priority, nil)
				store.EXPECT().
					CreateTransaction(gomock.Any()).Times(1).Return(transaction, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "CategoryNotFound",
			body: validBody,
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(models.Category{}, errors.New("category not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.JSONEq(t, `{"error": "category not found"}`, recorder.Body.String())
			},
		},
		{
			name: "PriorityNotFound",
			body: validBody,
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(category, nil)
				store.EXPECT().
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(models.Priority{}, errors.New("priority not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.JSONEq(t, `{"error": "priority not found"}`, recorder.Body.String())
			},
		},
		{
			name: "InvalidRequestBody",
			body: inValidBody,
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: validBody,
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).Times(1).Return(category, nil)
				store.EXPECT().
					GetPriority(gomock.Any(), gomock.Any()).Times(1).Return(priority, nil)
				store.EXPECT().
					CreateTransaction(gomock.Any()).Times(1).Return(models.Transaction{}, errors.New("could not create transaction"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				require.JSONEq(t, `{"error": "Could not create transaction"}`, recorder.Body.String())
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

			var bodyReader *strings.Reader
			if tt.body != nil {
				bodyBytes, err := json.Marshal(tt.body)
				require.NoError(t, err)
				bodyReader = strings.NewReader(string(bodyBytes))
			} else {
				bodyReader = strings.NewReader("")
			}

			url := "/smart-account/api/v1/transaction"
			request, err := http.NewRequest(http.MethodPost, url, bodyReader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestTransactionEdit(t *testing.T) {
	password := "Password123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
		Password:  string(encryptedPass),
	}

	categoryID := uuid.New()
	priorityID := uuid.New()
	transactionID := uuid.New()

	validBody := map[string]interface{}{
		"Title":       "edited transaction",
		"Category":    categoryID.String(),
		"Amount":      150,
		"Negative":    false,
		"Description": "edited description",
		"Priority":    priorityID.String(),
	}
	inValidBody := map[string]interface{}{
		"Title": 1,
	}

	category := models.Category{
		ID:          categoryID,
		Name:        "Category A",
		Description: "A test Category",
		UserID:      user.UID,
	}

	priority := models.Priority{
		ID:          priorityID,
		Name:        "Priority A",
		Description: "A test Priority",
		Level:       5,
		UserID:      user.UID,
	}

	existingTransaction := models.Transaction{
		ID:          transactionID,
		Title:       "existing transaction",
		CategoryID:  categoryID,
		PriorityID:  priorityID,
		Amount:      100,
		Negative:    false,
		Description: "existing description",
		UserID:      user.UID,
	}

	updatedTransaction := existingTransaction
	updatedTransaction.Title = validBody["Title"].(string)
	updatedTransaction.CategoryID = category.ID
	updatedTransaction.PriorityID = priority.ID
	updatedTransaction.Amount = validBody["Amount"].(int)
	updatedTransaction.Negative = validBody["Negative"].(bool)
	updatedTransaction.Description = validBody["Description"].(string)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		transactionID string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			body:          validBody,
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransaction(gomock.Any()).Times(1).Return(existingTransaction, nil)
				store.EXPECT().
					GetCategory(category.ID, gomock.Any()).Times(1).Return(category, nil)
				store.EXPECT().
					GetPriority(priority.ID, gomock.Any()).Times(1).Return(priority, nil)
				store.EXPECT().
					EditTransaction(gomock.Any()).Times(1).Return(updatedTransaction, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "TransactionNotFound",
			body:          validBody,
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransaction(gomock.Any()).Times(1).Return(models.Transaction{}, errors.New("transaction not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "CategoryNotFound",
			body:          validBody,
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransaction(gomock.Any()).Times(1).Return(existingTransaction, nil)
				store.EXPECT().
					GetCategory(category.ID, gomock.Any()).Times(1).Return(models.Category{}, errors.New("category not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "PriorityNotFound",
			body:          validBody,
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransaction(gomock.Any()).Times(1).Return(existingTransaction, nil)
				store.EXPECT().
					GetCategory(category.ID, gomock.Any()).Times(1).Return(category, nil)
				store.EXPECT().
					GetPriority(priority.ID, gomock.Any()).Times(1).Return(models.Priority{}, errors.New("priority not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "InvalidRequestBody",
			body:          inValidBody,
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "InternalServerError",
			body:          validBody,
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransaction(gomock.Any()).Times(1).Return(existingTransaction, nil)
				store.EXPECT().
					GetCategory(category.ID, gomock.Any()).Times(1).Return(category, nil)
				store.EXPECT().
					GetPriority(priority.ID, gomock.Any()).Times(1).Return(priority, nil)
				store.EXPECT().
					EditTransaction(gomock.Any()).Times(1).Return(models.Transaction{}, errors.New("could not edit transaction"))
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

			var bodyReader *strings.Reader
			if tt.body != nil {
				bodyBytes, err := json.Marshal(tt.body)
				require.NoError(t, err)
				bodyReader = strings.NewReader(string(bodyBytes))
			} else {
				bodyReader = strings.NewReader("")
			}

			url := fmt.Sprintf("/smart-account/api/v1/transaction/%s", tt.transactionID)
			request, err := http.NewRequest(http.MethodPut, url, bodyReader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestTransactionDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := models.User{
		UID:       uuid.New(),
		Email:     "user@test.com",
		FirstName: "test",
		LastName:  "user",
	}

	transactionID := uuid.New()

	testCases := []struct {
		name          string
		transactionID string
		buildStubs    func(store *mock_store.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					DeleteTransaction(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "TransactionNotFound",
			transactionID: transactionID.String(),
			buildStubs: func(store *mock_store.MockStore) {
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().
					DeleteTransaction(gomock.Any()).Times(1).Return(errors.New("transaction not found"))
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

			reader := strings.NewReader("")

			request, err := http.NewRequest("DELETE", fmt.Sprintf("/smart-account/api/v1/transaction/%s", tt.transactionID), reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}
