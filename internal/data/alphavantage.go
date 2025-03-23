// File: internal/data/alphavantage.go
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	baseURL       = "https://www.alphavantage.co/query"
	envAPIKeyName = "ALPHAVANTAGE_API_KEY"
)

// AlphaVantageClient is a client for the Alpha Vantage API
type AlphaVantageClient struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewAlphaVantageClient creates a new Alpha Vantage API client
func NewAlphaVantageClient() *AlphaVantageClient {
	apiKey := os.Getenv(envAPIKeyName)
	
	return &AlphaVantageClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SearchResult represents a search result from Alpha Vantage
type SearchResult struct {
	Symbol      string `json:"1. symbol"`
	Name        string `json:"2. name"`
	Type        string `json:"3. type"`
	Region      string `json:"4. region"`
	MarketOpen  string `json:"5. marketOpen"`
	MarketClose string `json:"6. marketClose"`
	Timezone    string `json:"7. timezone"`
	Currency    string `json:"8. currency"`
	MatchScore  string `json:"9. matchScore"`
}

// SearchResponse is the response from a search query
type SearchResponse struct {
	BestMatches []SearchResult `json:"bestMatches"`
}

// SearchSymbols searches for stock symbols by keyword
func (c *AlphaVantageClient) SearchSymbols(keywords string) ([]SearchResult, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("alpha vantage API key not set, please set the %s environment variable", envAPIKeyName)
	}

	// Build the request URL
	params := url.Values{}
	params.Add("function", "SYMBOL_SEARCH")
	params.Add("keywords", keywords)
	params.Add("apikey", c.APIKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the request
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Alpha Vantage: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Alpha Vantage response: %w", err)
	}

	// Parse the response
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("error parsing Alpha Vantage response: %w", err)
	}

	return searchResp.BestMatches, nil
}

// QuoteResult represents a global quote result
type QuoteResult struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

// QuoteResponse is the response from a quote query
type QuoteResponse struct {
	GlobalQuote QuoteResult `json:"Global Quote"`
}

// GetQuote gets a stock quote
func (c *AlphaVantageClient) GetQuote(symbol string) (*QuoteResult, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("alpha vantage API key not set, please set the %s environment variable", envAPIKeyName)
	}

	// Build the request URL
	params := url.Values{}
	params.Add("function", "GLOBAL_QUOTE")
	params.Add("symbol", symbol)
	params.Add("apikey", c.APIKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the request
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Alpha Vantage: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Alpha Vantage response: %w", err)
	}

	// Parse the response
	var quoteResp QuoteResponse
	if err := json.Unmarshal(body, &quoteResp); err != nil {
		return nil, fmt.Errorf("error parsing Alpha Vantage response: %w", err)
	}

	return &quoteResp.GlobalQuote, nil
}

// TimeSeriesData represents the data for a specific date
type TimeSeriesData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// DailyTimeSeriesResponse is the response from a daily time series query
type DailyTimeSeriesResponse struct {
	MetaData   map[string]string             `json:"Meta Data"`
	TimeSeries map[string]TimeSeriesData `json:"Time Series (Daily)"`
}

// GetDailyTimeSeries gets daily time series data for a symbol
func (c *AlphaVantageClient) GetDailyTimeSeries(symbol string, compact bool) (map[string]TimeSeriesData, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("alpha vantage API key not set, please set the %s environment variable", envAPIKeyName)
	}

	// Build the request URL
	params := url.Values{}
	params.Add("function", "TIME_SERIES_DAILY")
	params.Add("symbol", symbol)
	
	if compact {
		params.Add("outputsize", "compact") // Last 100 data points
	} else {
		params.Add("outputsize", "full") // Full history
	}
	
	params.Add("apikey", c.APIKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the request
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Alpha Vantage: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Alpha Vantage response: %w", err)
	}

	// Parse the response
	var timeSeriesResp DailyTimeSeriesResponse
	if err := json.Unmarshal(body, &timeSeriesResp); err != nil {
		return nil, fmt.Errorf("error parsing Alpha Vantage response: %w", err)
	}

	return timeSeriesResp.TimeSeries, nil
}

// IntradayTimeSeriesResponse is the response from an intraday time series query
type IntradayTimeSeriesResponse struct {
	MetaData   map[string]string             `json:"Meta Data"`
	TimeSeries map[string]TimeSeriesData `json:"Time Series (5min)"`
}

// GetIntradayTimeSeries gets intraday time series data for a symbol
func (c *AlphaVantageClient) GetIntradayTimeSeries(symbol string, interval string, compact bool) (map[string]TimeSeriesData, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("alpha vantage API key not set, please set the %s environment variable", envAPIKeyName)
	}

	// Build the request URL
	params := url.Values{}
	params.Add("function", "TIME_SERIES_INTRADAY")
	params.Add("symbol", symbol)
	params.Add("interval", interval) // 1min, 5min, 15min, 30min, 60min
	
	if compact {
		params.Add("outputsize", "compact") // Last 100 data points
	} else {
		params.Add("outputsize", "full") // Full history
	}
	
	params.Add("apikey", c.APIKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the request
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Alpha Vantage: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Alpha Vantage response: %w", err)
	}

	// Parse the response
	var timeSeriesResp map[string]interface{}
	if err := json.Unmarshal(body, &timeSeriesResp); err != nil {
		return nil, fmt.Errorf("error parsing Alpha Vantage response: %w", err)
	}

	// Find the time series key, which varies based on interval
	timeSeriesKey := fmt.Sprintf("Time Series (%s)", interval)
	timeSeriesData, ok := timeSeriesResp[timeSeriesKey].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format, time series data not found")
	}

	// Convert to the desired format
	result := make(map[string]TimeSeriesData)
	for date, data := range timeSeriesData {
		if dataMap, ok := data.(map[string]interface{}); ok {
			result[date] = TimeSeriesData{
				Open:   dataMap["1. open"].(string),
				High:   dataMap["2. high"].(string),
				Low:    dataMap["3. low"].(string),
				Close:  dataMap["4. close"].(string),
				Volume: dataMap["5. volume"].(string),
			}
		}
	}

	return result, nil
}