package user

import "time"

const (
	//
	// Database Constants.
	//
	UsersTableName = "users"

	UserIdFieldName = "id"

	Admin  Role = "admin"
	Intern Role = "intern"
)

type (
	//
	// Request models.
	//

	// UserRequest holds information of a user request.
	UserRequest struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}

	//
	// Response models.
	//

	// UserResponse holds information about a users response.
	UserResponse struct {
		Users      Users      `json:"users"`
		Pagination Pagination `json:"pagination"`
	}

	// Pagination holds pagination information.
	Pagination struct {
		Count  int `json:"count"`
		Page   int `json:"page"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	//
	// Data models.
	//

	// User holds information of a user.
	User struct {
		ID        string    `json:"id" db:"id"`
		Username  string    `json:"username" db:"username"`
		Name      string    `json:"name" db:"name"`
		Role      Role      `json:"role" db:"role"`
		CreatedAt time.Time `json:"createdAt" db:"created_at"`
		UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	}

	// Users represents a slice of users.
	Users []User

	// Role represents the role of a user.
	Role string
)
