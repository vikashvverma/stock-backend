package auth

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/vikashvverma/stock-backend/response"
)

type Authenticator interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type requestAuthenticator struct {
	Logger        *logrus.Logger
	APIKey        string
	AllowedRoutes []string
}

func New(l *logrus.Logger, apiKey string, ar []string) Authenticator {
	return &requestAuthenticator{Logger: l, APIKey: apiKey, AllowedRoutes: ar}
}

func (ra *requestAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	for _, route := range ra.AllowedRoutes {
		if route == r.RequestURI {
			next(w, r)
			return
		}
	}

	authToken := r.Header.Get("API-KEY")
	if authToken != ra.APIKey {
		ra.Logger.Errorf("Authenticator: unauthorized")
		response.Error{Reason: "forbidden"}.Forbidden(w)
		return
	}

	next(w, r)
}
