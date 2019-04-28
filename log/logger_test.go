package log

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	logger, hook := test.NewNullLogger()
	requestResponseLogger := New(logger)
	nrw := negroni.NewResponseWriter(httptest.NewRecorder())
	r := httptest.NewRequest(http.MethodGet, "/foo/bar/baz", nil)
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/foo/bar/baz", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}
	requestResponseLogger.ServeHTTP(nrw, r, handler)
	assert.Len(t, hook.Entries, 2)

	reqFields := logrus.Fields{
		"Client":     r.RemoteAddr,
		"Method":     http.MethodGet,
		"URL":        "/foo/bar/baz",
		"Referrer":   "",
		"User-Agent": "",
	}

	entry0 := hook.Entries[0]
	entry1 := hook.Entries[1]
	assert.Equal(t, reqFields, entry0.Data)
	assert.Equal(t, "GET", entry1.Data["Method"])
	assert.Equal(t, "/foo/bar/baz", entry1.Data["URL"])
	assert.Equal(t, 200, entry1.Data["StatusCode"])
	assert.Equal(t, "Request", hook.Entries[0].Message)
	assert.Equal(t, "Response", hook.Entries[1].Message)
}
