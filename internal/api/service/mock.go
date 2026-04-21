package service

import (
	"context"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/stretchr/testify/mock"
)

var _ Repository = &MockedRepository{}

// MockedRepository represents a mocked repository.
type MockedRepository struct {
	mock.Mock
}

func (mr *MockedRepository) DeleteUser(ctx context.Context, userId string) error {
	ret := mr.Called(ctx, userId)

	return ret.Error(0)
}

func (mr *MockedRepository) InsertUser(ctx context.Context, user userModel.User) error {
	ret := mr.Called(ctx, user)

	return ret.Error(0)
}

func (mr *MockedRepository) UpdateUser(ctx context.Context, user userModel.User) error {
	ret := mr.Called(ctx, user)

	return ret.Error(0)
}

func (mr *MockedRepository) SelectAllUsers(ctx context.Context, userParams userModel.Params) (userModel.Users, userModel.Pagination, error) {
	ret := mr.Called(ctx, userParams)

	return ret.Get(0).(userModel.Users), ret.Get(1).(userModel.Pagination), ret.Error(2)
}
