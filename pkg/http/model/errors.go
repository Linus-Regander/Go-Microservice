package http

import "errors"

var (
	// ErrMissingAccessToken represents an error for when an access token is missing.
	ErrMissingAccessToken = errors.New("missing access token")

	// ErrMissingURLParam represents an error for when a url param is missing.
	ErrMissingURLParam = errors.New("missing url param")

	// ErrMissingRequestBody represents an error for when request body is missing.
	ErrMissingRequestBody = errors.New("missing request body")
)
