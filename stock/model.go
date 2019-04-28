package stock

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	Id          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Symbol      string             `json:"symbol,omitempty"`
	Name        string             `json:"name,omitempty"`
	MarketCap   float64            `json:"marketCap,omitempty"`
	Sector      string             `json:"sector,omitempty"`
	Industry    string             `json:"industry,omitempty"`
	PricePoints []PricePoint       `json:"pricePoints,omitempty"`
}

type PricePoint struct {
	Date   time.Time `json:"date,omitempty"`
	Symbol string    `json:"symbol,omitempty"`
	Open   float64   `json:"open,omitempty"`
	Close  float64   `json:"close,omitempty"`
	Low    float64   `json:"low,omitempty"`
	High   float64   `json:"high,omitempty"`
	Volume float64   `json:"volume,omitempty"`
}
