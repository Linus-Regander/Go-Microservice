package model

const (
	//
	// URL Params.
	//

	UserIdParamKey = "userId"
)

// UserParams holds the structure of user params.
type UserParams struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
