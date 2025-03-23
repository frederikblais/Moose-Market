package models

import "time"

// Stock represents basic stock data
type Stock struct {
    Symbol    string  `json:"symbol"`
    Name      string  `json:"name"`
    Price     float64 `json:"price"`
    Change    float64 `json:"change"`
    ChangePercent float64 `json:"change_percent"`
    Open      float64 `json:"open"`
    High      float64 `json:"high"`
    Low       float64 `json:"low"`
    Volume    int64   `json:"volume"`
    MarketCap float64 `json:"market_cap"`
    PE        float64 `json:"pe"`
    Dividend  float64 `json:"dividend"`
    Exchange  string  `json:"exchange"`
    Timestamp int64   `json:"timestamp"`
}

// CandleData represents OHLC data for charting
type CandleData struct {
    Symbol    string        `json:"symbol"`
    Timeframe string        `json:"timeframe"` // 1m, 5m, 15m, 1h, 1d, 1w
    Candles   []CandleStick `json:"candles"`
}

// CandleStick represents a single candlestick in a chart
type CandleStick struct {
    Time   time.Time `json:"time"`
    Open   float64   `json:"open"`
    High   float64   `json:"high"`
    Low    float64   `json:"low"`
    Close  float64   `json:"close"`
    Volume int64     `json:"volume"`
}

// DrawingObject represents a drawing on a chart
type DrawingObject struct {
    ID       string    `json:"id"`
    Symbol   string    `json:"symbol"`
    Type     string    `json:"type"` // line, rectangle, ellipse, fibonacci, etc.
    Color    string    `json:"color"`
    Points   []Point   `json:"points"`
    Text     string    `json:"text,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

// Point represents a point in a drawing
type Point struct {
    X float64 `json:"x"` // X can be a timestamp or index
    Y float64 `json:"y"` // Y is the price level
}