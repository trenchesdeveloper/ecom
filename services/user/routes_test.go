package user

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocks "github.com/trenchesdeveloper/go-ecom/mocks/types"
	"github.com/trenchesdeveloper/go-ecom/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_handleRegister(t *testing.T) {
	userStore := &mocks.UserStore{}
	handler := NewHandler(userStore)

	t.Run("should return error if payload is invalid", func(t *testing.T) {
		invalidPayload := types.RegisterInput{
			FirstName:       "John",
			LastName:        "Doe",
			Email:           "john.doe.com",
			Password:        "password",
			ConfirmPassword: "password",
		}

		marshalled, err := json.Marshal(invalidPayload)

		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		require.NoError(t, err)

		res := httptest.NewRecorder()

		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)

	})

	t.Run("should correctly register a user", func(t *testing.T) {
		validPayload := types.RegisterInput{
			FirstName:       "John",
			LastName:        "Doe",
			Email:           "johnnrr@gmail.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		userStore.On("GetUserByEmail", validPayload.Email).Return(
			nil, nil,
		)

		userStore.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

		marshalled, err := json.Marshal(validPayload)

		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		require.NoError(t, err)

		res := httptest.NewRecorder()

		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(res, req)
		//require.Equal(t, "user not found", res.Body.String())
		require.Equal(t, http.StatusCreated, res.Code)
	})
}
