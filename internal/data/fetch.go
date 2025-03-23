package data

import (
	"math"
    "math/rand"
    "time"
    "fmt"
    "strings"

    "github.com/frederikblais/Moose-Market/internal/models"
)

// Mock databases for testing
var (
    symbolsNYSE = []string{"AAPL", "GOOGL", "MSFT", "TSLA", "AMZN", "V", "JNJ", "WMT", "PG", "JPM"}
    symbolsTSX  = []string{"RY", "TD", "BNS", "ENB", "CNR", "BCE", "CM", "BMO", "SU", "CP"}
    
    companyNames = map[string]string{
        "AAPL": "Apple Inc.",
        "GOOGL": "Alphabet Inc.",
        "MSFT": "Microsoft Corporation",
        "TSLA": "Tesla, Inc.",
        "AMZN": "Amazon.com, Inc.",
        "V": "Visa Inc.",
        "JNJ": "Johnson & Johnson",
        "WMT": "Walmart Inc.",
        "PG": "Procter & Gamble Co.",
        "JPM": "JPMorgan Chase & Co.",
        "RY": "Royal Bank of Canada",
        "TD": "Toronto-Dominion Bank",
        "BNS": "Bank of Nova Scotia",
        "ENB": "Enbridge Inc.",
        "CNR": "Canadian National Railway",
        "BCE": "BCE Inc.",
        "CM": "Canadian Imperial Bank of Commerce",
        "BMO": "Bank of Montreal",
        "SU": "Suncor Energy Inc.",
        "CP": "Canadian Pacific Railway Limited",
    }
    
    exchanges = map[string]string{
        "AAPL": "NASDAQ",
        "GOOGL": "NASDAQ",
        "MSFT": "NASDAQ",
        "TSLA": "NASDAQ",
        "AMZN": "NASDAQ",
        "V": "NYSE",
        "JNJ": "NYSE",
        "WMT": "NYSE",
        "PG": "NYSE",
        "JPM": "NYSE",
        "RY": "TSX",
        "TD": "TSX",
        "BNS": "TSX",
        "ENB": "TSX",
        "CNR": "TSX",
        "BCE": "TSX",
        "CM": "TSX",
        "BMO": "TSX",
        "SU": "TSX",
        "CP": "TSX",
    }
)

// GetAllSymbols returns all available stock symbols
func GetAllSymbols() []string {
    allSymbols := append(symbolsNYSE, symbolsTSX...)
    return allSymbols
}

// SearchSymbols searches for symbols that match the query
func SearchSymbols(query string) []models.Stock {
    query = strings.ToUpper(query)
    var results []models.Stock
    
    for _, symbol := range GetAllSymbols() {
        if strings.Contains(symbol, query) || strings.Contains(strings.ToUpper(companyNames[symbol]), query) {
            stock, _ := GetStockBySymbol(symbol)
            if stock != nil {
                results = append(results, *stock)
            }
        }
    }
    
    return results
}

// FetchMarketData returns market data for all stocks
func FetchMarketData() ([]models.Stock, error) {
    rand.Seed(time.Now().UnixNano())

    var stocks []models.Stock
    allSymbols := GetAllSymbols()
    
    for _, symbol := range allSymbols {
        stock, err := GetStockBySymbol(symbol)
        if err == nil {
            stocks = append(stocks, *stock)
        }
    }
    
    return stocks, nil
}

// GetStockBySymbol fetches data for a single stock by symbol
func GetStockBySymbol(symbol string) (*models.Stock, error) {
    rand.Seed(time.Now().UnixNano())
    
    // Check if symbol exists in our mock database
    companyName, exists := companyNames[symbol]
    if !exists {
        return nil, fmt.Errorf("symbol %s not found", symbol)
    }
    
    // Generate mock stock data
    basePrice := 100.0 + rand.Float64()*400.0
    change := (rand.Float64() - 0.5) * 10.0
    changePercent := (change / basePrice) * 100
    
    // Create a random daily range that makes sense
    open := basePrice - change*rand.Float64()
    high := open * (1 + rand.Float64()*0.05)
    low := open * (1 - rand.Float64()*0.05)
    
    // Ensure high is actually the highest
    if high < open {
        high, open = open, high
    }
    
    // Ensure low is actually the lowest
    if low > open {
        low, open = open, low
    }
    
    // Generate realistic volume based on price
    volume := int64(rand.Float64() * 10000000)
    
    // Calculate market cap (shares outstanding * price)
    sharesOutstanding := rand.Float64() * 10000000000
    marketCap := sharesOutstanding * basePrice
    
    // Generate a reasonable P/E ratio
    pe := 15.0 + rand.Float64()*25.0
    
    // Generate a small dividend for some stocks
    dividend := 0.0
    if rand.Float64() > 0.3 { // 70% chance of having a dividend
        dividend = basePrice * 0.01 * rand.Float64() * 3.0 // 0-3% yield
    }
    
    stock := &models.Stock{
        Symbol:        symbol,
        Name:          companyName,
        Price:         basePrice,
        Change:        change,
        ChangePercent: changePercent,
        Open:          open,
        High:          high,
        Low:           low,
        Volume:        volume,
        MarketCap:     marketCap,
        PE:            pe,
        Dividend:      dividend,
        Exchange:      exchanges[symbol],
        Timestamp:     time.Now().Unix(),
    }
    
    return stock, nil
}

// GetCandleData generates mock candle data for a symbol
func GetCandleData(symbol, timeframe string, count int) (*models.CandleData, error) {
    // Check if we have this symbol
    if _, exists := companyNames[symbol]; !exists {
        return nil, fmt.Errorf("symbol %s not found", symbol)
    }
    
    // Get initial stock price to use as a reference
    stock, _ := GetStockBySymbol(symbol)
    basePrice := stock.Price
    
    // Set up time intervals based on timeframe
    var interval time.Duration
    switch timeframe {
    case "1m":
        interval = time.Minute
    case "5m":
        interval = 5 * time.Minute
    case "15m":
        interval = 15 * time.Minute
    case "1h":
        interval = time.Hour
    case "1d":
        interval = 24 * time.Hour
    case "1w":
        interval = 7 * 24 * time.Hour
    default:
        interval = time.Hour
    }
    
    // Generate candles
    var candles []models.CandleStick
    
    // Start time (going back 'count' intervals from now)
    endTime := time.Now()
    startTime := endTime.Add(-interval * time.Duration(count))
    currentTime := startTime
    
    prevClose := basePrice * (0.9 + rand.Float64()*0.2) // Start around the base price
    
    for currentTime.Before(endTime) {
        // Calculate volatility based on timeframe
        var volatility float64
        switch timeframe {
        case "1m", "5m":
            volatility = 0.005
        case "15m", "1h":
            volatility = 0.01
        case "1d":
            volatility = 0.02
        case "1w":
            volatility = 0.04
        default:
            volatility = 0.02
        }
        
        // Generate OHLC data with some randomness but maintain a trend
        // maxMove := prevClose * volatility
        open := prevClose * (1 + (rand.Float64()-0.5)*0.01) // Small random change from previous close
        change := rand.Float64()*2 - 1 // Random value between -1 and 1
        close := open * (1 + change*volatility)
        
        high := math.Max(open, close) * (1 + rand.Float64()*volatility)
        low := math.Min(open, close) * (1 - rand.Float64()*volatility)
        
        // Generate volume (higher during market hours)
        var volume int64
        hour := currentTime.Hour()
        if hour >= 9 && hour <= 16 { // Market hours
            volume = int64(500000 + rand.Float64()*4500000)
        } else {
            volume = int64(50000 + rand.Float64()*450000)
        }
        
        candle := models.CandleStick{
            Time:   currentTime,
            Open:   open,
            High:   high,
            Low:    low,
            Close:  close,
            Volume: volume,
        }
        
        candles = append(candles, candle)
        prevClose = close
        currentTime = currentTime.Add(interval)
    }
    
    candleData := &models.CandleData{
        Symbol:    symbol,
        Timeframe: timeframe,
        Candles:   candles,
    }
    
    return candleData, nil
}