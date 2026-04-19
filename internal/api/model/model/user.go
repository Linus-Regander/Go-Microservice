package model

const (
	Admin  Role = "admin"
	Intern Role = "intern"
)

type (
	// User holds information of a user.
	User struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Role     Role   `json:"role"`
	}

	// Users represents a slice of users.
	Users []User

	// Role represents the role of a user.
	Role string
)
