package data

import (
    "math/rand"
    "time"

    "github.com/frederikblais/Moose-Market/internal/models"
)

var symbols = []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN", "RY", "TD", "BNS"}

func FetchMarketData() ([]models.Stock, error) {
    rand.Seed(time.Now().UnixNano())

    var stocks []models.Stock
    for _, symbol := range symbols {
        price := 100.0 + rand.Float64()*50.0
        change := (rand.Float64() - 0.5) * 2.0

        stocks = append(stocks, models.Stock{
            Symbol:    symbol,
            Name:      "Company " + symbol,
            Price:     price,
            Change:    change,
            Timestamp: time.Now().Unix(),
        })
    }
    return stocks, nil
};