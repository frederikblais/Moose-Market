// File: internal/ui/components/watchlist_minimal.go
package components

import (
    "fmt"
    // Remove unused import: time

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// MinimalWatchlistContainer is a simplified implementation to avoid the panic
type MinimalWatchlistContainer struct {
    container        *fyne.Container
    currentWatchlist *models.Watchlist
    onSelectStock    func(string)
    onWatchlistChanged func(string)
}

// CreateMinimalWatchlistContainer creates a simplified watchlist component
func CreateMinimalWatchlistContainer(onSelectStock func(string), onWatchlistChanged func(string)) *MinimalWatchlistContainer {
    // Create a label
    label := widget.NewLabel("Watchlist")
    label.TextStyle = fyne.TextStyle{Bold: true}
    
    // Create a button
    button := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
        // Empty function for now - will be set later
    })
    
    // Create a basic container
    content := container.NewVBox(
        container.NewBorder(nil, nil, label, button, nil),
        widget.NewLabel("No stocks in watchlist"),
    )
    
    return &MinimalWatchlistContainer{
        container: content,
        onSelectStock: onSelectStock,
        onWatchlistChanged: onWatchlistChanged,
    }
}

// GetContainer returns the container
func (w *MinimalWatchlistContainer) GetContainer() *fyne.Container {
    return w.container
}

// LoadWatchlists is a simplified implementation that won't panic
func (w *MinimalWatchlistContainer) LoadWatchlists() {
    // Get active profile
    profile := data.GetActiveProfile()
    if profile == nil || len(profile.Watchlists) == 0 {
        return
    }
    
    // Just store the first watchlist
    w.currentWatchlist = &profile.Watchlists[0]
    
    // If we have a callback, call it
    if w.onWatchlistChanged != nil {
        w.onWatchlistChanged(w.currentWatchlist.ID)
    }
}

// LoadWatchlistItems loads stocks - simplified version
func (w *MinimalWatchlistContainer) LoadWatchlistItems() {
    if w.currentWatchlist == nil {
        return
    }
    
    // Get stocks for symbols in the watchlist
    var items []*models.Stock
    for _, symbol := range w.currentWatchlist.Symbols {
        stock, err := data.GetStockBySymbol(symbol)
        if err == nil {
            items = append(items, stock)
        }
    }
    
    // Create a container for the items
    listContainer := container.NewVBox()
    
    // Add items to the list
    if len(items) == 0 {
        emptyText := widget.NewLabel("No stocks in this watchlist")
        emptyText.Alignment = fyne.TextAlignCenter
        listContainer.Add(emptyText)
    } else {
        for _, stock := range items {
            // Create a simple stock item - fixed version
            stockSymbol := stock.Symbol // Store symbol for closure
            
            // Create a button for each stock
            btn := widget.NewButton(fmt.Sprintf("%s: $%.2f (%+.2f)", 
                stock.Symbol, 
                stock.Price, 
                stock.Change), 
                func() {
                    if w.onSelectStock != nil {
                        w.onSelectStock(stockSymbol)
                    }
                })
            
            // Style the button based on stock change
            if stock.Change >= 0 {
                btn.Importance = widget.HighImportance
            } else {
                btn.Importance = widget.DangerImportance
            }
            
            listContainer.Add(btn)
        }
    }
    
    // Update the container - safely
    if w.container != nil && len(w.container.Objects) > 1 {
        w.container.Objects[1] = listContainer
        w.container.Refresh()
    }
}

// AddSymbol adds a symbol to the watchlist
func (w *MinimalWatchlistContainer) AddSymbol(symbol string) {
    if w.currentWatchlist == nil {
        return
    }
    
    // Check if symbol already exists
    for _, s := range w.currentWatchlist.Symbols {
        if s == symbol {
            return // Already in the watchlist
        }
    }
    
    // Add symbol to the watchlist
    w.currentWatchlist.Symbols = append(w.currentWatchlist.Symbols, symbol)
    
    // Update the profile
    profile := data.GetActiveProfile()
    if profile == nil {
        return
    }
    
    // Find and update this watchlist
    for i, wl := range profile.Watchlists {
        if wl.ID == w.currentWatchlist.ID {
            profile.Watchlists[i] = *w.currentWatchlist
            break
        }
    }
    
    // Save profile
    data.SaveProfile(profile)
    
    // Refresh display
    w.LoadWatchlistItems()
}

// RemoveSymbol removes a symbol from the watchlist
func (w *MinimalWatchlistContainer) RemoveSymbol(symbol string) {
    if w.currentWatchlist == nil {
        return
    }
    
    // Filter out the symbol
    var newSymbols []string
    for _, s := range w.currentWatchlist.Symbols {
        if s != symbol {
            newSymbols = append(newSymbols, s)
        }
    }
    
    w.currentWatchlist.Symbols = newSymbols
    
    // Update the profile
    profile := data.GetActiveProfile()
    if profile == nil {
        return
    }
    
    // Find and update this watchlist
    for i, wl := range profile.Watchlists {
        if wl.ID == w.currentWatchlist.ID {
            profile.Watchlists[i] = *w.currentWatchlist
            break
        }
    }
    
    // Save profile
    data.SaveProfile(profile)
    
    // Refresh display
    w.LoadWatchlistItems()
}