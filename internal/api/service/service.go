package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/google/uuid"
)

// ErrMalformedId represents an error for when an id is malformed, not of UUID format.
var ErrMalformedId = errors.New("malformed id")

type (
	// Repository represents the interface for repository functions.
	Repository interface {
		UserRepository
	}

	// UserRepository represents the interface for a user repository.
	UserRepository interface {
		DeleteUser(ctx context.Context, userId string) error
		InsertUser(ctx context.Context, user userModel.User) error
		UpdateUser(ctx context.Context, user userModel.User) error
		SelectAllUsers(ctx context.Context, userParams userModel.Params) (userModel.Users, userModel.Pagination, error)
	}

	// Service holds information about an API service.
	Service struct {
		log        *log.Logger
		repository Repository
	}
)

// New returns a new Service.
func New(logger *log.Logger, repository Repository) *Service {
	return &Service{
		log:        logger,
		repository: repository,
	}
}

// DeleteUser deletes a user from repository.
func (s *Service) DeleteUser(ctx context.Context, userId string) error {
	var err error

	if strings.TrimSpace(userId) == "" {
		return ErrMalformedId
	}

	if _, err = uuid.Parse(userId); err != nil {
		return ErrMalformedId
	}

	return s.repository.DeleteUser(ctx, userId)
}

// InsertUser inserts a user to repository.
func (s *Service) InsertUser(ctx context.Context, userRequest userModel.UserRequest) error {
	u := userModel.User{
		ID:        uuid.New().String(),
		Name:      userRequest.Name,
		Username:  userRequest.Username,
		Role:      userModel.Role(userRequest.Role),
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	return s.repository.InsertUser(ctx, u)
}

// UpdateUser updates user in repository.
func (s *Service) UpdateUser(ctx context.Context, userRequest userModel.UserRequest) error {
	u := userModel.User{
		Name:      userRequest.Name,
		Username:  userRequest.Username,
		Role:      userModel.Role(userRequest.Role),
		UpdatedAt: time.Now(),
	}

	return s.repository.UpdateUser(ctx, u)
}

// SelectAllUsers selects all users in repository, with params for filtering.
func (s *Service) SelectAllUsers(ctx context.Context, params userModel.Params) (userModel.UserResponse, error) {
	var userResponse userModel.UserResponse

	users, pagination, err := s.repository.SelectAllUsers(ctx, params)
	if err != nil {
		return userResponse, err
	}

	userResponse.Users = users
	userResponse.Pagination = pagination

	return userResponse, nil
}
