// File: internal/ui/components/heatmap.go
package components

import (
    "fmt"
    "image/color"
    "math"
    
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    
    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// HeatmapContainer represents the market heatmap component
type HeatmapContainer struct {
    container    *fyne.Container
    watchlistID  string
    stocks       []*models.Stock
    onSelectStock func(string)
}

// CreateHeatmapContainer creates the market heatmap container
func CreateHeatmapContainer(onSelectStock func(string)) *HeatmapContainer {
    // Title for the heatmap
    title := widget.NewLabel("Market Heatmap")
    title.TextStyle = fyne.TextStyle{Bold: true}
    
    // Create container for heatmap tiles
    heatmapGrid := container.NewGridWrap(fyne.NewSize(80, 80))
    
    // Scrollable container for the grid
    scrollContainer := container.NewScroll(heatmapGrid)
    
    // Put everything together
    heatmapContainer := container.NewBorder(title, nil, nil, nil, scrollContainer)
    
    return &HeatmapContainer{
        container:    heatmapContainer,
        onSelectStock: onSelectStock,
    }
}

// GetContainer returns the container for the heatmap
func (h *HeatmapContainer) GetContainer() *fyne.Container {
    return h.container
}

// SetWatchlist sets the watchlist to display in the heatmap
func (h *HeatmapContainer) SetWatchlist(watchlistID string) {
    h.watchlistID = watchlistID
    h.RefreshHeatmap()
}

// RefreshHeatmap refreshes the heatmap with the current watchlist data
func (h *HeatmapContainer) RefreshHeatmap() {
    // Get the active profile
    profile := data.GetActiveProfile()
    if profile == nil {
        return
    }
    
    // Find the watchlist
    var watchlist *models.Watchlist
    for i := range profile.Watchlists {
        if profile.Watchlists[i].ID == h.watchlistID {
            watchlist = &profile.Watchlists[i]
            break
        }
    }
    
    if watchlist == nil {
        return
    }
    
    // Get stocks for this watchlist
    var stocks []*models.Stock
    for _, symbol := range watchlist.Symbols {
        stock, err := data.GetStockBySymbol(symbol)
        if err == nil {
            stocks = append(stocks, stock)
        }
    }
    
    h.stocks = stocks
    
    // Get the grid container
    var heatmapGrid *fyne.Container
    if scrollContainer, ok := h.container.Objects[1].(*container.Scroll); ok {
        heatmapGrid = scrollContainer.Content.(*fyne.Container)
    }
    
    if heatmapGrid == nil {
        return
    }
    
    // Clear the grid
    heatmapGrid.Objects = nil
    
    // Create tiles for each stock
    for _, stock := range h.stocks {
        tile := h.createHeatmapTile(stock)
        heatmapGrid.Add(tile)
    }
    
    // Add placeholder message if empty
    if len(h.stocks) == 0 {
        emptyText := widget.NewLabel("No stocks in this watchlist")
        emptyText.Alignment = fyne.TextAlignCenter
        heatmapGrid.Add(emptyText)
    }
    
    // Refresh the container
    heatmapGrid.Refresh()
}

// createHeatmapTile creates a tile for the heatmap
func (h *HeatmapContainer) createHeatmapTile(stock *models.Stock) fyne.CanvasObject {
    // Calculate color based on percent change
    percent := stock.ChangePercent
    var tileColor color.Color
    
    if percent > 0 {
        // Green intensity based on positive percent
        intensity := uint8(math.Min(255, 100 + math.Abs(percent)*10))
        tileColor = color.NRGBA{R: 0, G: intensity, B: 0, A: 255}
    } else {
        // Red intensity based on negative percent
        intensity := uint8(math.Min(255, 100 + math.Abs(percent)*10))
        tileColor = color.NRGBA{R: intensity, G: 0, B: 0, A: 255}
    }
    
    // Create the tile
    tile := canvas.NewRectangle(tileColor)
    
    // Symbol text
    symbolText := canvas.NewText(stock.Symbol, color.White)
    symbolText.TextStyle = fyne.TextStyle{Bold: true}
    symbolText.Alignment = fyne.TextAlignCenter
    
    // Percent change text
    percentText := canvas.NewText(fmt.Sprintf("%.2f%%", stock.ChangePercent), color.White)
    percentText.TextSize = 12
    percentText.Alignment = fyne.TextAlignCenter
    
    // Layout the content
    content := container.NewVBox(
        symbolText,
        percentText,
    )
    
    // Make the tile clickable
    tileButton := widget.NewButton("", func() {
        if h.onSelectStock != nil {
            h.onSelectStock(stock.Symbol)
        }
    })
    tileButton.Importance = widget.LowImportance
    
    // Overlay everything
    return container.NewStack(
        tile,
        content,
        tileButton,
    )
}