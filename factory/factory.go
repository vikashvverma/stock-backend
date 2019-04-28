package factory

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vikashvverma/stock-backend/config"
	"github.com/vikashvverma/stock-backend/stock"
)

var dmDB sync.Once

// Factory represents factory for the service.
type Factory interface {
	Client() *mongo.Client
	Trader() stock.Trader
}

type factory struct {
	config  *config.Config
	logger  *logrus.Logger
	db      *sql.DB
	client  *mongo.Client
	seating map[int]int
}

// NewFactory returns a factory object.
func NewFactory(c *config.Config, l *logrus.Logger) Factory {
	return &factory{
		config: c,
		logger: l,
	}
}

// Client returns a new database connection.
func (f *factory) Client() *mongo.Client {
	var dbError error
	dmDB.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(f.config.DBConnection()))

		f.client = client
		dbError = err
	})

	if dbError != nil {
		f.logger.WithError(dbError).Fatalf("Could not establish connection to the DB: %s", dbError)
	}

	return f.client
}

// Trader returns a new stock.Trader instance
func (f *factory) Trader() stock.Trader {
	return stock.New(f.Client())
}
