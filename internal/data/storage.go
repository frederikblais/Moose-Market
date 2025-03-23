package data

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/frederikblais/Moose-Market/internal/models"
)

// Storage paths
const (
    dataDirectory = "moosemarket_data"
    profilesDir   = "profiles"
    stocksFile    = "stocks.json"
    candlesDir    = "candles"
    drawingsDir   = "drawings"
)

var (
    // Ensure thread safety for storage operations
    storageMutex sync.RWMutex
    
    // Store active profile in memory
    activeProfile *models.Profile
)

// Initialize makes sure all necessary directories exist
func Initialize() error {
    dirs := []string{
        dataDirectory,
        filepath.Join(dataDirectory, profilesDir),
        filepath.Join(dataDirectory, candlesDir),
        filepath.Join(dataDirectory, drawingsDir),
    }

    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
    }

    return nil
}

// Profile Operations

// GetProfiles returns a list of all available profiles
func GetProfiles() ([]models.Profile, error) {
    storageMutex.RLock()
    defer storageMutex.RUnlock()

    var profiles []models.Profile

    // Read the profiles directory
    profilesPath := filepath.Join(dataDirectory, profilesDir)
    files, err := os.ReadDir(profilesPath)
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
            continue
        }

        profilePath := filepath.Join(profilesPath, file.Name())
        profileData, err := os.ReadFile(profilePath)
        if err != nil {
            continue // Skip files we can't read
        }

        var profile models.Profile
        if err := json.Unmarshal(profileData, &profile); err != nil {
            continue // Skip invalid JSON
        }

        profiles = append(profiles, profile)
    }

    return profiles, nil
}

// GetProfileByID loads a profile by its ID
func GetProfileByID(id string) (*models.Profile, error) {
    storageMutex.RLock()
    defer storageMutex.RUnlock()

    profilePath := filepath.Join(dataDirectory, profilesDir, id+".json")
    profileData, err := os.ReadFile(profilePath)
    if err != nil {
        return nil, err
    }

    var profile models.Profile
    if err := json.Unmarshal(profileData, &profile); err != nil {
        return nil, err
    }

    return &profile, nil
}

// SaveProfile saves a profile to disk
func SaveProfile(profile *models.Profile) error {
    storageMutex.Lock()
    defer storageMutex.Unlock()

    profileData, err := json.MarshalIndent(profile, "", "  ")
    if err != nil {
        return err
    }

    profilePath := filepath.Join(dataDirectory, profilesDir, profile.ID+".json")
    return os.WriteFile(profilePath, profileData, 0644)
}

// CreateProfile creates a new profile
func CreateProfile(name string) (*models.Profile, error) {
    id := fmt.Sprintf("profile_%d", time.Now().Unix())
    now := time.Now()

    profile := &models.Profile{
        ID:           id,
        Name:         name,
        CreatedAt:    now,
        LastModified: now,
        Accounts:     []models.Account{},
        Watchlists: []models.Watchlist{
            {
                ID:        fmt.Sprintf("watchlist_%d", now.Unix()),
                Name:      "Default",
                Symbols:   []string{},
                CreatedAt: now,
            },
        },
        Settings: models.Settings{
            DarkMode:        true,
            Currency:        "CAD",
            RefreshInterval: 60,
        },
    }

    if err := SaveProfile(profile); err != nil {
        return nil, err
    }

    return profile, nil
}

// DeleteProfile removes a profile
func DeleteProfile(id string) error {
    storageMutex.Lock()
    defer storageMutex.Unlock()

    profilePath := filepath.Join(dataDirectory, profilesDir, id+".json")
    return os.Remove(profilePath)
}

// SetActiveProfile sets the currently active profile
func SetActiveProfile(profile *models.Profile) {
    storageMutex.Lock()
    defer storageMutex.Unlock()
    activeProfile = profile
}

// GetActiveProfile returns the currently active profile
func GetActiveProfile() *models.Profile {
    storageMutex.RLock()
    defer storageMutex.RUnlock()
    return activeProfile
}

// Stock Operations

// SaveStocks saves stock data to a JSON file
func SaveStocks(stocks []models.Stock) error {
    storageMutex.Lock()
    defer storageMutex.Unlock()

    bytes, err := json.MarshalIndent(stocks, "", "  ")
    if err != nil {
        return err
    }

    stocksPath := filepath.Join(dataDirectory, stocksFile)
    return os.WriteFile(stocksPath, bytes, 0644)
}

// LoadStocks loads stock data from a JSON file
func LoadStocks() ([]models.Stock, error) {
    storageMutex.RLock()
    defer storageMutex.RUnlock()

    var stocks []models.Stock

    stocksPath := filepath.Join(dataDirectory, stocksFile)
    if _, err := os.Stat(stocksPath); os.IsNotExist(err) {
        // If no file yet, return an empty slice
        return stocks, nil
    }

    bytes, err := os.ReadFile(stocksPath)
    if err != nil {
        return stocks, err
    }

    err = json.Unmarshal(bytes, &stocks)
    return stocks, err
}

// Candle Data Operations

// SaveCandleData saves candle data for a specific symbol and timeframe
func SaveCandleData(candles models.CandleData) error {
    storageMutex.Lock()
    defer storageMutex.Unlock()

    fileName := fmt.Sprintf("%s_%s.json", candles.Symbol, candles.Timeframe)
    filePath := filepath.Join(dataDirectory, candlesDir, fileName)

    data, err := json.MarshalIndent(candles, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filePath, data, 0644)
}

// LoadCandleData loads candle data for a specific symbol and timeframe
func LoadCandleData(symbol, timeframe string) (*models.CandleData, error) {
    storageMutex.RLock()
    defer storageMutex.RUnlock()

    fileName := fmt.Sprintf("%s_%s.json", symbol, timeframe)
    filePath := filepath.Join(dataDirectory, candlesDir, fileName)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil, errors.New("candle data not found")
    }

    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var candleData models.CandleData
    if err := json.Unmarshal(data, &candleData); err != nil {
        return nil, err
    }

    return &candleData, nil
}

// Drawing Operations

// SaveDrawings saves all drawings for a symbol
func SaveDrawings(symbol string, drawings []models.DrawingObject) error {
    storageMutex.Lock()
    defer storageMutex.Unlock()

    fileName := fmt.Sprintf("%s_drawings.json", symbol)
    filePath := filepath.Join(dataDirectory, drawingsDir, fileName)

    data, err := json.MarshalIndent(drawings, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filePath, data, 0644)
}

// LoadDrawings loads all drawings for a symbol
func LoadDrawings(symbol string) ([]models.DrawingObject, error) {
    storageMutex.RLock()
    defer storageMutex.RUnlock()

    fileName := fmt.Sprintf("%s_drawings.json", symbol)
    filePath := filepath.Join(dataDirectory, drawingsDir, fileName)

    var drawings []models.DrawingObject

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return drawings, nil // Return empty slice if no drawings yet
    }

    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(data, &drawings); err != nil {
        return nil, err
    }

    return drawings, nil
}