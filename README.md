# Echo Error Handler Middleware ![coverage](https://raw.githubusercontent.com/xwinata/echo_error_handler/badges/.badges/main/coverage.svg)

A customizable error handler middleware for the Echo framework. This middleware intercepts errors returned by handlers and returns consistent JSON responses. It supports custom error types with their own status codes.

## Features

- Handles standard `echo.HTTPError`
- Supports custom error types via `WithCustomError`
- Returns JSON responses with appropriate status codes

## Installation
```bash
go get github.com/xwinata/echo_error_handler
```

## Usage
### Define your custom error structure that implement `customError` interface
```
type MyCustomError struct {
	Message string
	Code    int
}

func (e MyCustomError) Error() string {
	return e.Message
}

func (e MyCustomError) ErrorCode() int {
	return e.Code
}
```
#### `customError` interface definition
```
type customError interface {
	Error() string
	ErrorCode() int
}
```
### Register the middleware
```
import (
	"github.com/labstack/echo/v4"
	echoErrorHandler "github.com/xwinata/echo_error_handler"
)

func main() {
	e := echo.New()

	// Create the middleware with your custom error type
	errHandler := echoErrorHandler.New(
		echoErrorHandler.WithCustomError(MyCustomError{}),
	)

	e.Use(errHandler.HandlerFunc)

	// Example route
	e.GET("/test", func(c echo.Context) error {
		return MyCustomError{
			Message: "invalid request",
			Code:    400,
		}
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```
### The behavior
- If the error implements customError, the middleware returns:
```
{
  "Message": "your error message",
  "Code": 400
}
```
- If it's an `*echo.HTTPError`, Echo’s default error response is returned.
- If it’s another error type, a 500 Internal Server Error is returned with the message.
