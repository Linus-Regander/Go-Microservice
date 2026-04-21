package service

import (
	"context"
	"log"
	"testing"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_DeleteUser(t *testing.T) {
	testCases := map[string]struct {
		expectedErr   error
		repositoryErr error
		userId        string
	}{
		"DeleteUser - success": {
			expectedErr:   nil,
			repositoryErr: nil,
			userId:        "550e8400-e29b-41d4-a716-446655440000",
		},
		"DeleteUser - error - malformed uuid": {
			expectedErr:   ErrMalformedId,
			repositoryErr: nil,
			userId:        "123",
		},
		"DeleteUser - error - repository returned error": {
			expectedErr:   assert.AnError,
			repositoryErr: assert.AnError,
			userId:        "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mr := &MockedRepository{}
			mr.On("DeleteUser", mock.Anything, mock.Anything).Return(tc.repositoryErr)

			s := New(&log.Logger{}, mr)
			require.NotNil(t, s)

			err := s.DeleteUser(context.Background(), tc.userId)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestService_InsertUser(t *testing.T) {
	testCases := map[string]struct {
		expectedErr   error
		repositoryErr error
		request       userModel.UserRequest
	}{
		"InsertUser - success": {
			expectedErr:   nil,
			repositoryErr: nil,
			request: userModel.UserRequest{
				Name:     "Test",
				Username: "Test",
				Role:     string(userModel.Admin),
			},
		},
		"InsertUser - error - repository returned error": {
			expectedErr:   assert.AnError,
			repositoryErr: assert.AnError,
			request: userModel.UserRequest{
				Name:     "Test",
				Username: "Test",
				Role:     string(userModel.Admin),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mr := &MockedRepository{}
			mr.On("InsertUser", mock.Anything, mock.Anything).Return(tc.repositoryErr)

			s := New(&log.Logger{}, mr)
			require.NotNil(t, s)

			err := s.InsertUser(context.Background(), tc.request)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	testCases := map[string]struct {
		expectedErr   error
		repositoryErr error
		request       userModel.UserRequest
	}{
		"UpdateUser - success": {
			expectedErr:   nil,
			repositoryErr: nil,
			request: userModel.UserRequest{
				Name:     "Test",
				Username: "Test",
				Role:     string(userModel.Admin),
			},
		},
		"UpdateUser - error - repository returned error": {
			expectedErr:   assert.AnError,
			repositoryErr: assert.AnError,
			request: userModel.UserRequest{
				Name:     "Test",
				Username: "Test",
				Role:     string(userModel.Admin),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mr := &MockedRepository{}
			mr.On("UpdateUser", mock.Anything, mock.Anything).Return(tc.repositoryErr)

			s := New(&log.Logger{}, mr)
			require.NotNil(t, s)

			err := s.UpdateUser(context.Background(), tc.request)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestService_SelectAllUsers(t *testing.T) {
	testCases := map[string]struct {
		expectedErr      error
		repositoryErr    error
		users            userModel.Users
		userPagination   userModel.Pagination
		expectedResponse userModel.UserResponse
	}{
		"UpdateUser - success": {
			expectedErr:   nil,
			repositoryErr: nil,
			users: userModel.Users{
				{
					Name:     "Test",
					Username: "Test",
					Role:     userModel.Admin,
				},
				{
					Name:     "Testo",
					Username: "Testo",
					Role:     userModel.Intern,
				},
			},
			userPagination: userModel.Pagination{
				Count:  2,
				Page:   1,
				Limit:  30,
				Offset: 0,
			},
			expectedResponse: userModel.UserResponse{
				Users: userModel.Users{
					{
						Name:     "Test",
						Username: "Test",
						Role:     userModel.Admin,
					},
					{
						Name:     "Testo",
						Username: "Testo",
						Role:     userModel.Intern,
					},
				},
				Pagination: userModel.Pagination{
					Count:  2,
					Page:   1,
					Limit:  30,
					Offset: 0,
				},
			},
		},
		"UpdateUser - error - repository returned error": {
			expectedErr:      assert.AnError,
			repositoryErr:    assert.AnError,
			users:            userModel.Users{},
			userPagination:   userModel.Pagination{},
			expectedResponse: userModel.UserResponse{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mr := &MockedRepository{}
			mr.On("SelectAllUsers", mock.Anything, mock.Anything).Return(tc.users, tc.userPagination, tc.repositoryErr)

			s := New(&log.Logger{}, mr)
			require.NotNil(t, s)

			userResponse, err := s.SelectAllUsers(context.Background(), userModel.Params{})
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, userResponse)
		})
	}
}
