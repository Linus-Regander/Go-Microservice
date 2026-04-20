package handler

import (
	"context"
	"log"
	"net/http"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/Linus-Regander/Go-Microservice/pkg/http/methods"
)

type (
	// Service represents the interface for service functions.
	Service interface {
		UserService
	}

	// UserService represents the interface for a user service.
	UserService interface {
		DeleteUser(ctx context.Context, userId string) error
		InsertUser(ctx context.Context, userRequest userModel.UserRequest) error
		UpdateUser(ctx context.Context, userRequest userModel.UserRequest) error
		SelectAllUsers(ctx context.Context, params userModel.UserParams) (userModel.UserResponse, error)
	}

	// Handler holds information about an API handler.
	Handler struct {
		log         *log.Logger
		userService UserService
	}
)

// New returns a new handler.
func New(logger *log.Logger, userService UserService) *Handler {
	return &Handler{
		log:         logger,
		userService: userService,
	}
}

// DeleteUser represents the handler for deleting a user.
func (h *Handler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := methods.Authorized(ctx); err != nil {
			methods.ResponseJSON[any](w, err, http.StatusUnauthorized, nil)

			return
		}

		userId, err := methods.UrlParam(r)
		if err != nil { // should never occur since chi handles 404 if id is missing.
			methods.ResponseJSON[any](w, err, http.StatusBadRequest, nil)

			return
		}

		if err = h.userService.DeleteUser(ctx, userId); err != nil {
			methods.ResponseJSON[any](w, err, http.StatusInternalServerError, nil)

			return
		}

		methods.ResponseJSON[userModel.User](w, nil, http.StatusNoContent, userModel.User{})
	}
}

// InsertUser represents the handler for inserting a user.
func (h *Handler) InsertUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := methods.Authorized(ctx)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusUnauthorized, nil)

			return
		}

		userRequest, err := methods.ParseRequestBody[userModel.UserRequest](r)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusBadRequest, nil)

			return
		}

		if err = h.userService.InsertUser(ctx, *userRequest); err != nil {
			methods.ResponseJSON[any](w, err, http.StatusInternalServerError, nil)

			return
		}

		methods.ResponseJSON[any](w, nil, http.StatusCreated, nil)
	}
}

// UpdateUser represents the handler for updating a user.
func (h *Handler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := methods.Authorized(ctx)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusUnauthorized, nil)

			return
		}

		requestBody, err := methods.ParseRequestBody[userModel.UserRequest](r)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusBadRequest, nil)

			return
		}

		if err = h.userService.UpdateUser(ctx, *requestBody); err != nil {
			methods.ResponseJSON[any](w, err, http.StatusInternalServerError, nil)

			return
		}

		methods.ResponseJSON[any](w, nil, http.StatusNoContent, nil)
	}
}

// SelectUsers represents the handler for selecting users.
func (h *Handler) SelectUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userResponse userModel.UserResponse

		ctx := r.Context()

		err := methods.Authorized(ctx)
		if err != nil {
			methods.ResponseJSON[userModel.UserResponse](w, err, http.StatusUnauthorized, userResponse)

			return
		}

		if userResponse, err = h.userService.SelectAllUsers(ctx, userModel.UserParams{}); err != nil {
			methods.ResponseJSON[userModel.UserResponse](w, err, http.StatusInternalServerError, userResponse)
		}

		methods.ResponseJSON[userModel.UserResponse](w, nil, http.StatusOK, userResponse)
	}
}
