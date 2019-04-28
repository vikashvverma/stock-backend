package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/vikashvverma/stock-backend/config"
	"github.com/vikashvverma/stock-backend/factory"
	"github.com/vikashvverma/stock-backend/handler"
	"github.com/vikashvverma/stock-backend/healthcheck"
)

// Router returns the router for all the API handler.
func Router(f factory.Factory, c *config.Config, l *logrus.Logger) *mux.Router {
	l.Out = c.LogFile()
	l.Level = logrus.Level(c.LogLevel())

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", healthcheck.Self).Methods(http.MethodGet)
	router.HandleFunc("/stock/{name}", handler.Find(f.Trader(), f, l)).Methods(http.MethodGet)
	router.HandleFunc("/stock/{from}/{to}", handler.FindList(f.Trader(), f, l)).Queries("ticker", "{ticker}").Methods(http.MethodGet)
	router.HandleFunc("/stock/top/{from}/{to}", handler.Top(f.Trader(), f, l)).Methods(http.MethodGet)

	return router
}
