package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vikashvverma/stock-backend/config"
	"github.com/vikashvverma/stock-backend/constants"
	"github.com/vikashvverma/stock-backend/factory"
	"github.com/vikashvverma/stock-backend/stock"
)

func main() {
	var c *config.Config
	var err error

	useFlags := !(len(os.Args) > 2 && os.Args[1] == "-config")
	if useFlags {
		c, err = config.FromFlags(os.Args)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		c, err = config.FromFile(os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
	}

	l := logrus.New()
	f := factory.NewFactory(c, l)

	client := f.Client()

	stocks, err := parseStock(c.Stock())
	if err != nil {
		l.Fatalf("unable to parse stock csv file: %s", err)
	}

	collection := client.Database(constants.Database).Collection(constants.Collection)

	ctx := context.Background()
	csvFile, err := os.Open(c.Data())
	if err != nil {
		log.Fatalf("error opening stock file: %s", err)
	}

	insert(stocks, csvFile, collection, ctx)

}

func insert(stocks map[string]stock.Stock, csvFile *os.File, collection *mongo.Collection, ctx context.Context) {
	// Read File into a Variable
	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		logrus.Fatalf("Error reading file: %s", csvFile.Name())
	}

	// Loop through lines & turn into object
	for i, line := range lines {
		if i == 0 { // skip header
			continue
		}
		date, _ := time.Parse(
			fmt.Sprintf(
				"%s-%s-%s",
				constants.StdLongYear,
				constants.StdZeroMonth,
				constants.StdZeroDay,
			),
			strings.Split(line[0], " ")[0])

		d := stock.PricePoint{
			Date:   date,
			Symbol: line[1],
			Open:   parseToFloat(line[2]),
			Close:  parseToFloat(line[3]),
			Low:    parseToFloat(line[4]),
			High:   parseToFloat(line[5]),
			Volume: parseToFloat(line[6]),
		}

		if _, ok := stocks[line[1]]; !ok {
			stocks[line[1]] = stock.Stock{Symbol: line[1]}
		}

		st := stocks[line[1]]
		st.PricePoints = append(st.PricePoints, d)

		stocks[line[1]] = st
	}
	for _, v := range stocks {
		_, err := collection.InsertOne(ctx, v, &options.InsertOneOptions{})
		if err != nil {
			fmt.Printf("insert failed for the symbol: %s", v.Symbol)
			continue
		}

		fmt.Printf("Inserted: %s\n", v.Symbol)

	}

}

func parseToFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logrus.Fatalf("parseToFloat: unable to parse %v to float: %s", str, err)
	}
	return f
}

func parseStock(filePath string) (map[string]stock.Stock, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("parseStock: error reading file: %s", err)
	}

	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parseStock: error reading all lines: %v", err)
	}

	stocks := make(map[string]stock.Stock, 0)

	for i, line := range lines {
		if i == 0 { //skip header
			continue
		}

		marketCap, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("parseStock: unable to parse market cap at line %d: %s", i+1, err)
		}

		s := stock.Stock{
			Symbol:      line[0],
			Name:        line[1],
			MarketCap:   marketCap,
			Sector:      line[3],
			Industry:    line[4],
			PricePoints: []stock.PricePoint{},
		}
		stocks[line[0]] = s
	}

	return stocks, nil
}
