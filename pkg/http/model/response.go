package http

type (
	// ErrorResponse holds the structure of a custom erroneous HTTP response.
	ErrorResponse[T any] struct {
		Response   T      `json:"response"`
		Error      string `json:"error"`
		StatusCode int    `json:"statusCode"`
	}

	// SuccessfulResponse holds the structure of a custom successful HTTP response.
	SuccessfulResponse[T any] struct {
		Response   T   `json:"response"`
		StatusCode int `json:"statusCode"`
	}
)
