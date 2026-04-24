package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	//
	// Main API paths.
	//

	// mainPath represents the main api path.
	mainPath = "/service"

	//
	// User API paths.
	//

	// userByIdPath represents the path for delete user endpoint.
	userByIdPath = "/user/{id}"

	// userPath represents the path for update and insert user endpoints.
	userPath = "/user"

	// usersPath represents the path for select users endpoint.
	usersPath = "/users"
)

type (
	// Handler holds the interface to handler functions.
	Handler interface {
		UserHandler
	}

	// UserHandler holds the interface to user handler functions.
	UserHandler interface {
		DeleteUser() http.HandlerFunc
		InsertUser() http.HandlerFunc
		UpdateUser() http.HandlerFunc
		SelectUsers() http.HandlerFunc
	}

	// Router holds information about an API router.
	Router struct {
		Handler   Handler
		chiRouter *chi.Mux
	}
)

// New returns a new router instance.
func New(handler Handler) *Router {
	return &Router{
		Handler:   handler,
		chiRouter: chi.NewRouter(),
	}
}

// SetupChi sets up the Chi api with endpoints.
func (r *Router) SetupChi() (string, http.Handler) {
	r.chiRouter.Delete(userByIdPath, r.Handler.DeleteUser())
	r.chiRouter.Put(userPath, r.Handler.UpdateUser())
	r.chiRouter.Post(userPath, r.Handler.InsertUser())
	r.chiRouter.Get(usersPath, r.Handler.SelectUsers())

	return mainPath, r.chiRouter
}
