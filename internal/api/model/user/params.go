package user

const (
	//
	// URL Params.
	//

	IdParamKey = "id"
)

// Params holds the structure of user params.
type Params struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Username string `json:"username" db:"username"`
	Role     string `json:"role" db:"role"`
}
