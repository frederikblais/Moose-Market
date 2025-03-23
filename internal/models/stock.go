package models

type Stock struct {
    Symbol    string  `json:"symbol"`
    Name      string  `json:"name"`
    Price     float64 `json:"price"`
    Change    float64 `json:"change"`
    Timestamp int64   `json:"timestamp"`
}