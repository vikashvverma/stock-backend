package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/vikashvverma/stock-backend/constants"
	"github.com/vikashvverma/stock-backend/factory"
	"github.com/vikashvverma/stock-backend/response"
	"github.com/vikashvverma/stock-backend/stock"
)

// Find represents find API handler.
func Find(t stock.Trader, f factory.Factory, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok {
			l.Errorf("Find: could not read 'name' from path params")
			response.Response{Errors: &response.Error{Reason: "path params not valid"}}.ClientError(w)
			return
		}

		stock, err := t.Find(name)
		if err != nil {
			l.WithError(err).Errorf("Find: error getting price points")
			response.Response{Errors: &response.Error{Reason: "could not find anything"}}.ServerError(w)
			return
		}

		response.Response{
			Success: true,
			Result:  stock,
		}.Send(w)

	}
}

// Top represents top API handler.
func Top(t stock.Trader, f factory.Factory, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fromDateString, ok := vars["from"]
		if !ok {
			l.Errorf("Top: could not read `from` Date from path params")
			response.Response{Errors: &response.Error{Reason: "path params not valid"}}.ClientError(w)
			return
		}

		fromDate, err := time.Parse(fmt.Sprintf("%s-%s-%s", constants.StdZeroDay, constants.StdZeroMonth, constants.StdLongYear), fromDateString)
		if err != nil {
			l.WithError(err).Errorf("Top: invalid `from` date: %s", fromDateString)
			response.Response{Errors: &response.Error{Reason: fmt.Sprintf("invalid from date: %s", fromDateString)}}.ClientError(w)
			return
		}

		toDateString, ok := vars["to"]
		if !ok {
			l.Errorf("Top: could not read `to` Date from path params")
			response.Response{Errors: &response.Error{Reason: "path params not valid"}}.ClientError(w)
			return
		}

		toDate, err := time.Parse(fmt.Sprintf("%s-%s-%s", constants.StdZeroDay, constants.StdZeroMonth, constants.StdLongYear), toDateString)
		if err != nil {
			l.WithError(err).Errorf("Top: invalid `to` date: %s", toDateString)
			response.Response{Errors: &response.Error{Reason: fmt.Sprintf("invalid `to` date: %s", toDateString)}}.ClientError(w)
			return
		}

		topStock, err := t.Top(fromDate, toDate, true)
		if err != nil {
			l.WithError(err).Errorf("Top: error getting top stocks")
			response.Response{Errors: &response.Error{Reason: "could not find anything"}}.ServerError(w)
			return
		}

		bottomStock, err := t.Top(fromDate, toDate, false)
		if err != nil {
			l.WithError(err).Errorf("Top: error getting bottom stocks")
			response.Response{Errors: &response.Error{Reason: "could not find anything"}}.ServerError(w)
			return
		}

		response.Response{
			Success: true,
			Result:  map[string]interface{}{"best": topStock, "least": bottomStock},
		}.Send(w)

	}
}

// FindList represents List API handler.
func FindList(t stock.Trader, f factory.Factory, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fromDateString, ok := vars["from"]
		if !ok {
			l.Errorf("FindList: could not read `from` Date from path params")
			response.Response{Errors: &response.Error{Reason: "path params not valid"}}.ClientError(w)
			return
		}

		fromDate, err := time.Parse(fmt.Sprintf("%s-%s-%s", constants.StdZeroDay, constants.StdZeroMonth, constants.StdLongYear), fromDateString)
		if err != nil {
			l.WithError(err).Errorf("FindList: invalid `from` date: %s", fromDateString)
			response.Response{Errors: &response.Error{Reason: fmt.Sprintf("invalid from date: %s", fromDateString)}}.ClientError(w)
			return
		}

		toDateString, ok := vars["to"]
		if !ok {
			l.Errorf("FindList: could not read `to` Date from path params")
			response.Response{Errors: &response.Error{Reason: "path params not valid"}}.ClientError(w)
			return
		}

		toDate, err := time.Parse(fmt.Sprintf("%s-%s-%s", constants.StdZeroDay, constants.StdZeroMonth, constants.StdLongYear), toDateString)
		if err != nil {
			l.WithError(err).Errorf("FindList: invalid `to` date: %s", toDateString)
			response.Response{Errors: &response.Error{Reason: fmt.Sprintf("invalid `to` date: %s", toDateString)}}.ClientError(w)
			return
		}

		queryParams := r.URL.Query()
		tickers := queryParams["ticker"]

		stocks, err := t.FindAll(tickers, fromDate, toDate)
		if err != nil {
			l.WithError(err).Errorf("FindList: error getting price points")
			response.Response{Errors: &response.Error{Reason: "could not find anything"}}.ServerError(w)
			return
		}

		response.Response{
			Success: true,
			Result:  stocks,
		}.Send(w)

	}
}
