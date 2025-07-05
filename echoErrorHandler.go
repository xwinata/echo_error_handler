// Package echo_error_handler provides a customizable error handler middleware
// for the Echo web framework. It supports handling both Echo HTTP errors and
// user-defined custom errors with structured JSON responses.
package echo_error_handler

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// customError defines an interface for user-defined errors that include
// an HTTP status code.
type customError interface {
	Error() string
	ErrorCode() int
}

// options holds the configurable options for the error handler.
type options struct {
	customError customError
}

// optionFunc represents a functional option for configuring the error handler.
type optionFunc func(*options)

// WithCustomError allows the user to pass a custom error implementation
// that satisfies the customError interface. The middleware will return
// this error with its associated status code when matched.
func WithCustomError(err customError) optionFunc {
	return func(o *options) {
		o.customError = err
	}
}

// ErrorHandler is the main struct for the error handler middleware.
type ErrorHandler struct {
	options *options
}

// New creates a new ErrorHandler instance with optional configuration.
// It accepts functional options like WithCustomError.
func New(opts ...optionFunc) *ErrorHandler {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}

	return &ErrorHandler{
		options: options,
	}
}

// HandlerFunc returns an Echo middleware handler that catches errors
// returned from route handlers. It checks for user-defined custom errors,
// Echo HTTP errors, and defaults to 500 Internal Server Error.
func (h *ErrorHandler) HandlerFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			var echoErr *echo.HTTPError
			var customError customError

			// Handle custom error type
			if h.options.customError != nil {
				if errors.As(err, &customError) {
					return c.JSON(customError.ErrorCode(), customError)
				}
			}

			// Handle standard Echo HTTP error
			if errors.As(err, &echoErr) {
				return c.JSON(echoErr.Code, echoErr)
			}

			// Fallback to 500 Internal Server Error
			return c.JSON(echo.ErrInternalServerError.Code, echo.HTTPError{
				Internal: err,
				Message:  err.Error(),
			})
		}
		return err
	}
}
