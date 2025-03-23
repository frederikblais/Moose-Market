// File: internal/models/profile.go
package models

import "time"

// Profile represents a user profile that can contain multiple accounts
type Profile struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    CreatedAt    time.Time `json:"created_at"`
    LastModified time.Time `json:"last_modified"`
    Accounts     []Account `json:"accounts"`
    Watchlists   []Watchlist `json:"watchlists"`
    Settings     Settings  `json:"settings"`
}

// Account represents a financial account within a profile (TFSA, RRSP, etc)
type Account struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"` // TFSA, RRSP, FHSA, etc.
    Balance     float64   `json:"balance"`
    Positions   []Position `json:"positions"`
    CreatedAt   time.Time `json:"created_at"`
    LastUpdated time.Time `json:"last_updated"`
}

// Position represents a holding in a specific stock within an account
type Position struct {
    StockSymbol  string    `json:"stock_symbol"`
    Quantity     float64   `json:"quantity"`
    AverageCost  float64   `json:"average_cost"`
    Transactions []Transaction `json:"transactions"`
}

// Transaction represents a buy/sell transaction for a position
type Transaction struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"` // buy, sell
    Quantity  float64   `json:"quantity"`
    Price     float64   `json:"price"`
    Date      time.Time `json:"date"`
    Commission float64  `json:"commission"`
    Notes     string    `json:"notes"`
}

// Watchlist represents a collection of stocks to monitor
type Watchlist struct {
    ID        string   `json:"id"`
    Name      string   `json:"name"`
    Symbols   []string `json:"symbols"`
    CreatedAt time.Time `json:"created_at"`
}

// Settings contains user preferences
type Settings struct {
    DarkMode           bool   `json:"dark_mode"`
    Currency           string `json:"currency"`
    RefreshInterval    int    `json:"refresh_interval"` // in seconds
    DefaultWatchlistID string `json:"default_watchlist_id"`
}