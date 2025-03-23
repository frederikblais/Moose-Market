// File: internal/ui/components/watchlist.go
package components

import (
    "fmt"
    "image/color"
	"time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
    
    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// WatchlistContainer represents the watchlist component
type WatchlistContainer struct {
    container      *fyne.Container
    tabContainer   *container.TabItem
    watchlists     []models.Watchlist
    currentWatchlist *models.Watchlist
    items          []*models.Stock
    onSelectStock  func(string)
    onWatchlistChanged func(string)
}

// CreateWatchlistContainer creates the watchlist component
func CreateWatchlistContainer(onSelectStock func(string), onWatchlistChanged func(string)) *WatchlistContainer {
    // Create an empty container to hold the watchlist
    listContainer := container.NewVBox()
    
    // Create the tab container
    tabContainer := container.NewAppTabs()
    
    // Create watchlist container - this is the "w" we're referring to
    watchlistContainer := &WatchlistContainer{
        container:     container.NewBorder(nil, nil, nil, nil, listContainer),
        tabContainer:  nil,
        onSelectStock: onSelectStock,
        onWatchlistChanged: onWatchlistChanged,
    }
    
    // Create an "Add Watchlist" button
    // This is where we're making the change
    addWatchlistBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
        // Call the method on our watchlistContainer instance
        watchlistContainer.CreateWatchlistDialog()
    })
    
    // Create the header with tabs and add button
    header := container.NewBorder(nil, nil, nil, addWatchlistBtn, tabContainer)
    
    // Update the main container
    watchlistContainer.container = container.NewBorder(header, nil, nil, nil, listContainer)
    
    return watchlistContainer
}

// GetContainer returns the container for the watchlist
func (w *WatchlistContainer) GetContainer() *fyne.Container {
    return w.container
}

// AddWatchlist adds a new watchlist
func (w *WatchlistContainer) AddWatchlist(name string) {
    // Get the active profile
    profile := data.GetActiveProfile()
    if profile == nil {
        return
    }
    
    // Create a new watchlist
    newWatchlist := models.Watchlist{
        ID:       fmt.Sprintf("watchlist_%d", time.Now().Unix()),
        Name:     name,
        Symbols:  []string{},
        CreatedAt: time.Now(),
    }
    
    // Add watchlist to profile
    profile.Watchlists = append(profile.Watchlists, newWatchlist)
    
    // Save the profile
    data.SaveProfile(profile)
    
    // Refresh UI
    w.LoadWatchlists()
}

// LoadWatchlists loads all watchlists from the active profile
func (w *WatchlistContainer) LoadWatchlists() {
    // Get the active profile
    profile := data.GetActiveProfile()
    if profile == nil {
        // Just create an empty state
        w.container = container.NewVBox(widget.NewLabel("No profile loaded"))
        return
    }
    
    // Store the watchlists
    w.watchlists = profile.Watchlists
    
    // If no watchlists, create a default one
    if len(w.watchlists) == 0 {
        // Create a default watchlist
        defaultWatchlist := models.Watchlist{
            ID:        fmt.Sprintf("watchlist_%d", time.Now().Unix()),
            Name:      "Default",
            Symbols:   []string{},
            CreatedAt: time.Now(),
        }
        profile.Watchlists = append(profile.Watchlists, defaultWatchlist)
        w.watchlists = profile.Watchlists
        data.SaveProfile(profile)
    }
    
    // Create new tab container and tabs
    tabContainer := container.NewAppTabs()
    
    // Create a tab for each watchlist
    for i := range w.watchlists {
        // Local variable to avoid closure issues
        watchlistIndex := i
        watchlist := &w.watchlists[watchlistIndex]
        
        // Create container for this watchlist's items
        listContainer := container.NewVBox()
        
        // Add a tab for this watchlist
        tab := container.NewTabItem(watchlist.Name, listContainer)
        tabContainer.Append(tab)
    }
    
    // Set tab selection callback
    tabContainer.OnSelected = func(tab *container.TabItem) {
        // Find the corresponding watchlist
        for i, wl := range w.watchlists {
            if tab.Text == wl.Name {
                w.currentWatchlist = &w.watchlists[i]
                w.LoadWatchlistItems()
                if w.onWatchlistChanged != nil {
                    w.onWatchlistChanged(wl.ID)
                }
                break
            }
        }
    }
    
    // Select the first tab if we have any
    if len(w.watchlists) > 0 {
        w.currentWatchlist = &w.watchlists[0]
        tabContainer.SelectIndex(0)
    }
    
    // Create the add button
    addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), w.CreateWatchlistDialog)
    
    // Create the header with tabs and add button
    header := container.NewBorder(nil, nil, nil, addButton, tabContainer)
    
    // Create a new VBox for watchlist items
    itemsContainer := container.NewVBox()
    
    // Create the complete container from scratch
    w.container = container.NewBorder(header, nil, nil, nil, itemsContainer)
    
    // Load the items for the current watchlist
    if w.currentWatchlist != nil {
        w.LoadWatchlistItems()
    }
}

// LoadWatchlistItems loads the stocks for the current watchlist
func (w *WatchlistContainer) LoadWatchlistItems() {
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
    
    w.items = items
    
    // Get the content container
    var listContainer *fyne.Container
    if tabContainer, ok := w.container.Objects[0].(*fyne.Container).Objects[0].(*container.AppTabs); ok {
        if tab := tabContainer.CurrentTab(); tab != nil {
            listContainer = tab.Content.(*fyne.Container)
        }
    }
    
    if listContainer == nil {
        return
    }
    
    // Clear existing items
    listContainer.Objects = nil
    
    // Add items to the list
    if len(items) == 0 {
        emptyText := widget.NewLabel("No stocks in this watchlist")
        emptyText.Alignment = fyne.TextAlignCenter
        listContainer.Add(emptyText)
    } else {
        for _, stock := range items {
            // Create stock item row
            stockItem := w.createWatchlistItem(stock)
            listContainer.Add(stockItem)
        }
    }
    
    // Refresh the container
    listContainer.Refresh()
}

// AddSymbol adds a symbol to the current watchlist
func (w *WatchlistContainer) AddSymbol(symbol string) {
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
    
    // Refresh UI
    w.LoadWatchlistItems()
}

// RemoveSymbol removes a symbol from the current watchlist
func (w *WatchlistContainer) RemoveSymbol(symbol string) {
    if w.currentWatchlist == nil {
        return
    }
    
    // Remove symbol from watchlist
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
    
    // Refresh UI
    w.LoadWatchlistItems()
}

// createWatchlistItem creates a UI row for a stock in the watchlist
func (w *WatchlistContainer) createWatchlistItem(stock *models.Stock) fyne.CanvasObject {
    // Create elements for the stock item
    symbolText := canvas.NewText(stock.Symbol, color.White)
    symbolText.TextStyle = fyne.TextStyle{Bold: true}
    symbolText.TextSize = 16
    
    priceText := canvas.NewText(fmt.Sprintf("$%.2f", stock.Price), color.White)
    priceText.TextSize = 16
    
    var changeColor color.Color
    if stock.Change >= 0 {
        changeColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green
    } else {
        changeColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // Red
    }
    
    changeText := canvas.NewText(fmt.Sprintf("%.2f (%.2f%%)", stock.Change, stock.ChangePercent), changeColor)
    changeText.TextSize = 14
    
    // Container for text elements
    textContainer := container.NewVBox(
        container.NewHBox(symbolText, priceText),
        changeText,
    )
    
    // Delete button
    deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
        w.RemoveSymbol(stock.Symbol)
    })
    deleteBtn.Importance = widget.LowImportance
    
    // Create the item container
    itemContainer := container.NewBorder(
        nil, nil, textContainer, deleteBtn,
        nil,
    )
    
    // Make the whole item clickable
    itemButton := widget.NewButton("", func() {
        if w.onSelectStock != nil {
            w.onSelectStock(stock.Symbol)
        }
    })
    itemButton.Importance = widget.LowImportance
    
    // Overlay the button on the container
    combined := container.NewStack(
        itemContainer,
        itemButton,
    )
    
    return combined
}

func (w *WatchlistContainer) CreateWatchlistDialog() {
    // Get a valid window reference
    var parentWindow fyne.Window
    if len(fyne.CurrentApp().Driver().AllWindows()) > 0 {
        parentWindow = fyne.CurrentApp().Driver().AllWindows()[0]
    }
    
    // Create an entry for the watchlist name
    nameEntry := widget.NewEntry()
    nameEntry.SetPlaceHolder("Enter watchlist name")
    
    // Create a form dialog
    formItems := []*widget.FormItem{
        widget.NewFormItem("Name", nameEntry),
    }
    
    dialog.ShowForm("New Watchlist", "Create", "Cancel", formItems, 
        func(confirmed bool) {
            if confirmed && nameEntry.Text != "" {
                w.AddWatchlist(nameEntry.Text)
            }
        }, 
        parentWindow,
    )
}