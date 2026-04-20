package handler

import (
	"context"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/stretchr/testify/mock"
)

var _ Service = &MockedService{}

// MockedService represents a mocked service.
type MockedService struct {
	mock.Mock
}

func (ms *MockedService) DeleteUser(ctx context.Context, userId string) error {
	ret := ms.Called(ctx, userId)

	return ret.Error(0)
}

func (ms *MockedService) InsertUser(ctx context.Context, userRequest userModel.UserRequest) error {
	ret := ms.Called(ctx, userRequest)

	return ret.Error(0)
}

func (ms *MockedService) UpdateUser(ctx context.Context, userRequest userModel.UserRequest) error {
	ret := ms.Called(ctx, userRequest)

	return ret.Error(0)
}

func (ms *MockedService) SelectAllUsers(ctx context.Context, params userModel.UserParams) (userModel.UserResponse, error) {
	ret := ms.Called(ctx, params)

	return ret.Get(0).(userModel.UserResponse), ret.Error(1)
}
