package methods

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	internalModel "github.com/Linus-Regander/Go-Microservice/pkg/http/model"

	"github.com/go-chi/chi/v5"
)

// ResponesJSON returns a http response in JSON format.
func ResponseJSON[T any](w http.ResponseWriter, err error, statusCode int, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err != nil {
		if jErr := json.NewEncoder(w).Encode(internalModel.ErrorResponse[T]{
			Response:   response,
			Error:      err.Error(),
			StatusCode: statusCode,
		}); jErr != nil {
			//
		}

		return
	}

	if jErr := json.NewEncoder(w).Encode(internalModel.SuccessfulResponse[T]{
		Response:   response,
		StatusCode: statusCode,
	}); jErr != nil {
		//
	}
}

// Authorized validates if context contains authorized information.
func Authorized(ctx context.Context) error {
	// dummy access token, but proves the expected behavior.
	if ctx.Value("accessToken") == nil {
		return internalModel.ErrMissingAccessToken
	}

	return nil
}

// UrlParam returns a url param.
func UrlParam(r *http.Request) (string, error) {
	userId := chi.URLParam(r, user.IdParamKey)

	if strings.TrimSpace(userId) == "" {
		return "", errors.New(fmt.Sprintf("%w: %s", internalModel.ErrMissingURLParam.Error(), userId))
	}

	return userId, nil
}

// ParseRequestBody parses a request body of generic type T.
func ParseRequestBody[T any](r *http.Request) (*T, error) {
	if r.Body == nil || r.Body == http.NoBody {
		return nil, internalModel.ErrMissingRequestBody
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			return
		}
	}(r.Body)

	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}
