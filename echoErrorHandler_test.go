package echo_error_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type dummyError struct {
	LeMessage string `json:"le_message"`
	KodeError int    `json:"kode_error"`
}

func (cE dummyError) Error() string {
	return cE.LeMessage
}

func (cE dummyError) ErrorCode() int {
	return cE.KodeError
}

func TestErrorHandler(t *testing.T) {
	errorHandler := New(WithCustomError(dummyError{}))

	e := echo.New()
	e.Use(errorHandler.HandlerFunc)
	e.GET("/bad-request", func(c echo.Context) (err error) {
		err = &dummyError{
			KodeError: http.StatusBadRequest,
			LeMessage: "Error:Field validation",
		}

		return err
	})
	e.GET("/echo-error", func(c echo.Context) (err error) {
		newErr := errors.New("echo error occured")
		err = &echo.HTTPError{
			Internal: newErr,
			Code:     http.StatusBadGateway,
			Message:  newErr.Error(),
		}

		return err
	})
	e.GET("/random-error", func(c echo.Context) (err error) {
		err = errors.New("random error")

		return err
	})
	e.GET("/ok", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, "ok")
	})

	t.Run("TestErrorHandler: pass next", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var body string
		json.Unmarshal(rec.Body.Bytes(), &body)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, body, "ok")
	})

	t.Run("TestErrorHandler: returns standard error for recognized errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bad-request", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var body map[string]any
		json.Unmarshal(rec.Body.Bytes(), &body)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, fmt.Sprint(body["kode_error"]), fmt.Sprint(http.StatusBadRequest))
		assert.Equal(t, body["le_message"], "Error:Field validation")

		req = httptest.NewRequest(http.MethodGet, "/echo-error", nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		json.Unmarshal(rec.Body.Bytes(), &body)

		assert.Equal(t, http.StatusBadGateway, rec.Code)
		assert.Equal(t, fmt.Sprint(body["message"]), "echo error occured")
	})

	t.Run("TestErrorHandler: returns standard echo http error for unrecognized errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/random-error", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var body map[string]any
		json.Unmarshal(rec.Body.Bytes(), &body)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, body["message"], "random error")
	})
}
