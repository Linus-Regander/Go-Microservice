package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/Linus-Regander/Go-Microservice/internal/api/model/model"
	"github.com/Linus-Regander/Go-Microservice/pkg/http/methods"
)

type (
	// UserService represents the interface for user service methods.
	UserService interface {
		Delete(ctx context.Context, userId string) error
		Insert(ctx context.Context, user model.User) error
		Update(ctx context.Context, user model.User) error
		SelectAll(ctx context.Context, params model.UserParams) (model.Users, error)
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

		err := methods.Authorized(ctx)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusUnauthorized, nil)

			return
		}

		userId, err := methods.UrlParam(r)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusBadRequest, nil)

			return
		}

		if err = h.userService.Delete(ctx, userId); err != nil {
			methods.ResponseJSON[any](w, err, http.StatusInternalServerError, nil)

			return
		}

		methods.ResponseJSON[model.User](w, nil, http.StatusNoContent, model.User{})
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

		userRequest, err := methods.ParseRequestBody[model.User](r)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusBadRequest, nil)

			return
		}

		if err = h.userService.Insert(ctx, *userRequest); err != nil {
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

		requestBody, err := methods.ParseRequestBody[model.User](r)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusBadRequest, nil)

			return
		}

		if err = h.userService.Update(ctx, *requestBody); err != nil {
			methods.ResponseJSON[any](w, err, http.StatusInternalServerError, nil)

			return
		}

		methods.ResponseJSON[any](w, nil, http.StatusNoContent, nil)
	}
}

// SelectUsers represents the handler for selecting users.
func (h *Handler) SelectUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := methods.Authorized(ctx)
		if err != nil {
			methods.ResponseJSON[any](w, err, http.StatusUnauthorized, nil)

			return
		}

		var users model.Users

		if users, err = h.userService.SelectAll(ctx, model.UserParams{}); err != nil {
			methods.ResponseJSON[model.Users](w, err, http.StatusInternalServerError, users)
		}

		methods.ResponseJSON[model.Users](w, nil, http.StatusOK, users)
	}
}
