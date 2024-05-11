package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
				// category := models.Category{
				// 	ID:          uuid.New(),
				// 	Name:        "Category A",
				// 	Description: "A test Category",
				// 	UserID:      user.UID,
				// }
				// priority := models.Priority{
				// 	ID:          uuid.New(),
				// 	Name:        "Priority A",
				// 	Description: "A test Priority",
				// 	Level:       5,
				// 	UserID:      user.UID,
				// }
				// transactions := []models.Transaction{{
				// 	ID:          uuid.New(),
				// 	Title:       "transaction title",
				// 	CategoryID:  category.ID,
				// 	Category:    category,
				// 	Priority:    priority,
				// 	Amount:      100,
				// 	Negative:    true,
				// 	Description: "a test transaction",
				// 	PriorityID:  priority.ID,
				// 	UserID:      user.UID,
				// }}
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return()
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransactions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "error in transactions",
			buildStubs: func(store *mock_store.MockStore) {
				// category := models.Category{
				// 	ID:          uuid.New(),
				// 	Name:        "Category A",
				// 	Description: "A test Category",
				// 	UserID:      user.UID,
				// }
				// priority := models.Priority{
				// 	ID:          uuid.New(),
				// 	Name:        "Priority A",
				// 	Description: "A test Priority",
				// 	Level:       5,
				// 	UserID:      user.UID,
				// }
				// transactions := []models.Transaction{{
				// 	ID:          uuid.New(),
				// 	Title:       "transaction title",
				// 	CategoryID:  category.ID,
				// 	Category:    category,
				// 	Priority:    priority,
				// 	Amount:      100,
				// 	Negative:    true,
				// 	Description: "a test transaction",
				// 	PriorityID:  priority.ID,
				// 	UserID:      user.UID,
				// }}
				store.EXPECT().
					ReadToken(gomock.Any()).Times(1).Return()
				store.EXPECT().
					GetUserFromToken(gomock.Any()).Times(1).Return(user)
				store.EXPECT().
					GetTransactions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("error in transactions"))
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
