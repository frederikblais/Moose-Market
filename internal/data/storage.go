package data

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/frederikblais/Moose-Market/internal/models"
)

const dataFile = "moosemarket_data.json"

// SaveStocks saves stock data to a JSON file.
func SaveStocks(stocks []models.Stock) error {
    bytes, err := json.MarshalIndent(stocks, "", "  ")
    if err != nil {
        return err
    }

    // For cross-platform, ensure a relative or user-home path:
    // Here we just store it in the current working directory
    return ioutil.WriteFile(dataFile, bytes, 0644)
}

// LoadStocks loads stock data from a JSON file, if it exists.
func LoadStocks() ([]models.Stock, error) {
    var stocks []models.Stock

    if _, err := os.Stat(dataFile); os.IsNotExist(err) {
        // If no file yet, return an empty slice
        return stocks, nil
    }

    bytes, err := ioutil.ReadFile(filepath.Clean(dataFile))
    if err != nil {
        return stocks, err
    }

    err = json.Unmarshal(bytes, &stocks)
    return stocks, err
}
