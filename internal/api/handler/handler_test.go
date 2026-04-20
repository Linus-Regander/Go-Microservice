package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	internalModel "github.com/Linus-Regander/Go-Microservice/pkg/http/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_DeleteUser(t *testing.T) {
	var testUserId = "550e8400-e29b-41d4-a716-446655440000"

	testCases := map[string]struct {
		ctx                func() context.Context
		serviceErr         error
		requestURL         string
		expectedStatusCode int
		expectedResponse   any
	}{
		"DeleteUser - success": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL:         fmt.Sprintf("/user/%s", testUserId),
			serviceErr:         nil,
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
		"DeleteUser - error - unauthorized token": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "invalidToken", "") // invalid mocked access token.
			},
			requestURL:         fmt.Sprintf("/user/%s", testUserId),
			serviceErr:         nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]any{
				"error":      internalModel.ErrMissingAccessToken.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusUnauthorized),
			},
		},
		"DeleteUser - error - service returned error": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL:         fmt.Sprintf("/user/%s", testUserId),
			serviceErr:         assert.AnError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]any{
				"error":      assert.AnError.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusInternalServerError),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := tc.ctx()

			ms := &MockedService{}
			ms.On("DeleteUser", mock.Anything, mock.Anything).Return(tc.serviceErr)

			handler := New(&log.Logger{}, ms)

			r := chi.NewRouter()
			r.Delete("/user/{id}", handler.DeleteUser())

			recorder := httptest.NewRecorder()

			req, reqErr := http.NewRequestWithContext(ctx, http.MethodDelete, tc.requestURL, nil)
			require.NoError(t, reqErr)

			r.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer func() {
				require.NoError(t, res.Body.Close())
			}()

			var result interface{}

			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			if res.StatusCode != http.StatusNoContent {
				require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
				require.Equal(t, tc.expectedResponse, result)
			}
		})
	}
}

func TestHandler_InsertUser(t *testing.T) {
	testCases := map[string]struct {
		ctx                func() context.Context
		serviceErr         error
		requestURL         string
		requestBody        []byte
		expectedStatusCode int
		expectedResponse   any
	}{
		"InsertUser - success": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal(userModel.UserRequest{
					Username: "test",
					Name:     "Test Testsson",
					Role:     string(userModel.Admin),
				})

				return body
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   nil,
		},
		"InsertUser - error - missing body": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				return nil
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"error":      internalModel.ErrMissingRequestBody.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusBadRequest),
			},
		},

		"InsertUser - error - malformed body": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal("invalid")

				return body
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"error":      "json: cannot unmarshal string into Go value of type user.UserRequest",
				"response":   nil,
				"statusCode": float64(http.StatusBadRequest),
			},
		},
		"InsertUser - error - unauthorized token": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "invalidToken", "") // invalid mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal(userModel.UserRequest{
					Username: "test",
					Name:     "Test Testsson",
					Role:     string(userModel.Admin),
				})

				return body
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]any{
				"error":      internalModel.ErrMissingAccessToken.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusUnauthorized),
			},
		},
		"InsertUser - error - service returned error": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal(userModel.UserRequest{
					Username: "test",
					Name:     "Test Testsson",
					Role:     string(userModel.Admin),
				})

				return body
			}(),
			serviceErr:         assert.AnError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]any{
				"error":      assert.AnError.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusInternalServerError),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := tc.ctx()

			ms := &MockedService{}
			ms.On("InsertUser", mock.Anything, mock.Anything).Return(tc.serviceErr)

			handler := New(&log.Logger{}, ms)

			r := chi.NewRouter()
			r.Post("/user", handler.InsertUser())

			recorder := httptest.NewRecorder()

			req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, tc.requestURL, bytes.NewReader(tc.requestBody))
			require.NoError(t, reqErr)

			r.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer func() {
				require.NoError(t, res.Body.Close())
			}()

			var result interface{}

			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			if res.StatusCode != http.StatusCreated {
				require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
				require.Equal(t, tc.expectedResponse, result)
			}
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	testCases := map[string]struct {
		ctx                func() context.Context
		serviceErr         error
		requestURL         string
		requestBody        []byte
		expectedStatusCode int
		expectedResponse   any
	}{
		"UpdateUser - success": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal(userModel.UserRequest{
					Username: "test",
					Name:     "Test Testsson",
					Role:     string(userModel.Admin),
				})

				return body
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
		"UpdateUser - error - missing body": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				return nil
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"error":      internalModel.ErrMissingRequestBody.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusBadRequest),
			},
		},

		"UpdateUser - error - malformed body": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal("invalid")

				return body
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"error":      "json: cannot unmarshal string into Go value of type user.UserRequest",
				"response":   nil,
				"statusCode": float64(http.StatusBadRequest),
			},
		},
		"UpdateUser - error - unauthorized token": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "invalidToken", "") // invalid mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal(userModel.UserRequest{
					Username: "test",
					Name:     "Test Testsson",
					Role:     string(userModel.Admin),
				})

				return body
			}(),
			serviceErr:         nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]any{
				"error":      internalModel.ErrMissingAccessToken.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusUnauthorized),
			},
		},
		"UpdateUser - error - service returned error": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL: fmt.Sprintf("/user"),
			requestBody: func() []byte {
				body, _ := json.Marshal(userModel.UserRequest{
					Username: "test",
					Name:     "Test Testsson",
					Role:     string(userModel.Admin),
				})

				return body
			}(),
			serviceErr:         assert.AnError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]any{
				"error":      assert.AnError.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusInternalServerError),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := tc.ctx()

			ms := &MockedService{}
			ms.On("UpdateUser", mock.Anything, mock.Anything).Return(tc.serviceErr)

			handler := New(&log.Logger{}, ms)

			r := chi.NewRouter()
			r.Put("/user", handler.UpdateUser())

			recorder := httptest.NewRecorder()

			req, reqErr := http.NewRequestWithContext(ctx, http.MethodPut, tc.requestURL, bytes.NewReader(tc.requestBody))
			require.NoError(t, reqErr)

			r.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer func() {
				require.NoError(t, res.Body.Close())
			}()

			var result interface{}

			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			if res.StatusCode != http.StatusNoContent {
				require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
				require.Equal(t, tc.expectedResponse, result)
			}
		})
	}
}

func TestHandler_SelectUsers(t *testing.T) {
	testCases := map[string]struct {
		ctx                func() context.Context
		serviceErr         error
		requestURL         string
		expectedStatusCode int
		expectedResponse   any
	}{
		"SelectedUsers- success": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL:         fmt.Sprintf("/users"),
			serviceErr:         nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse: userModel.UserResponse{
				Users: userModel.Users{
					{
						ID:       "550e8400-e29b-41d4-a716-446655440000",
						Username: "test",
						Name:     "Test Testsson",
						Role:     userModel.Admin,
					},
					{
						ID:       "560e8400-e29b-41d4-a716-446655440000",
						Username: "test2",
						Name:     "Test2 Testsson",
						Role:     userModel.Intern,
					},
				},
				Pagination: userModel.Pagination{
					Count:  2,
					Page:   1,
					Offset: 0,
					Limit:  30,
				},
			},
		},
		"SelectedUsers - error - unauthorized token": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "invalidToken", "") // invalid mocked access token.
			},
			requestURL:         fmt.Sprintf("/users"),
			serviceErr:         nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]any{
				"error":      internalModel.ErrMissingAccessToken.Error(),
				"response":   map[string]interface{}{},
				"statusCode": float64(http.StatusUnauthorized),
			},
		},
		"SelectUsers - error - service returned error": {
			ctx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, "accessToken", "") // mocked access token.
			},
			requestURL:         fmt.Sprintf("/users"),
			serviceErr:         assert.AnError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]any{
				"error":      assert.AnError.Error(),
				"response":   nil,
				"statusCode": float64(http.StatusInternalServerError),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := tc.ctx()

			ms := &MockedService{}
			ms.On("SelectAllUsers", mock.Anything, mock.Anything).Return(tc.expectedResponse, tc.serviceErr)

			h := New(&log.Logger{}, ms)
			require.NotNil(t, h)

			r := chi.NewRouter()
			r.Get("/users", h.SelectUsers())

			recorder := httptest.NewRecorder()

			req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, tc.requestURL, nil)
			require.NoError(t, reqErr)

			r.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer func() {
				require.NoError(t, res.Body.Close())
			}()

			var result interface{}

			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			if res.StatusCode != http.StatusNoContent {
				require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
				require.Equal(t, tc.expectedResponse, result)
			}
		})
	}
}
