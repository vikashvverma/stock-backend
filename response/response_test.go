package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	s := Response{Success: true}
	w := httptest.NewRecorder()

	err := s.Send(w)
	require.NoError(t, err, "Expected no error writing JSON response")
	require.Nil(t, err, "Expected err to be nil")

	result := w.Result()
	var response Response
	err = json.NewDecoder(result.Body).Decode(&response)
	require.NoError(t, err, "Expected no error reading JSON response")
	require.Nil(t, err, "Expected err to be nil")

	assert.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, s, response)
}

func TestServerError(t *testing.T) {
	e := Response{Errors: &Error{Reason: "unexpected error happened"}}

	w := httptest.NewRecorder()
	err := e.ServerError(w)
	require.NoError(t, err, "Expected no err writing JSON response")
	require.Nil(t, err, "Expected err to be nil")

	result := w.Result()
	var response Response
	err = json.NewDecoder(result.Body).Decode(&response)
	require.NoError(t, err, "Expected no error reading response body")
	require.Nil(t, err, "Expected err to be nil")

	assert.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, e, response)
}

func TestClientError(t *testing.T) {
	e := Response{Errors: &Error{Reason: "invalid request"}}
	w := httptest.NewRecorder()

	err := e.ClientError(w)
	require.NoError(t, err, "Expected no error writing JSON response")
	require.Nil(t, err, "Expected err to be nil")

	result := w.Result()
	var response Response
	err = json.NewDecoder(result.Body).Decode(&response)
	require.NoError(t, err, "Expected no error reading response body")
	require.Nil(t, err, "Expected err to be nil")

	assert.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, e, response)
}
