package apierrors

import (
	"net/http"

	"github.com/go-chi/render"
)

// PublicError interface used to determine if error can be publicly displayed on
// a GQL response
type PublicError interface {
	// PublicError returns the public facing error message
	PublicError() string

	// ErrorCode returns the application error code
	ErrorCode() int
}

// PublicAPIError interface to determine if an error can be publicly displayed
// on an API response
type PublicAPIError interface {
	PublicError

	// HTTPStatusCode returns the HTTP code to be used in the response
	HTTPStatusCode() int
}

// PublicErrorDetails interface to determine if an error has extra details to
// to show in error
type PublicErrorDetails interface {
	Details() []string
}

// ResponseError generic error struct used to construct error responses.
type ResponseError struct {
	error          error    `json:"-"`
	httpStatusCode int      `json:"-"`
	ErrorCode      int      `json:"code,omitempty"`
	ErrorMessage   string   `json:"error,omitempty"`
	Details        []string `json:"details"`
}

// NewResponseError returns a new generic API error
func NewResponseError(msg string, errorCode, statusCode int, details []string, err error) *ResponseError {
	if details == nil {
		details = []string{}
	}
	return &ResponseError{
		ErrorMessage:   msg,
		ErrorCode:      errorCode,
		httpStatusCode: statusCode,
		Details:        details,
		error:          err,
	}
}

// NewInternalError returns a new generic internal error
func NewInternalError(err error) *ResponseError {
	return NewResponseError("internal error", CodeInternal, http.StatusInternalServerError, nil, err)
}

// Render conform to render.Renderer interface
func (e *ResponseError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.httpStatusCode)
	return nil
}
