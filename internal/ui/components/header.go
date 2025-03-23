// File: internal/ui/components/header.go
package components

import (
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/driver/desktop"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/theme"
    "image/color"

    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// SearchResult represents a search result to display in the dropdown
type SearchResult struct {
    Symbol   string
    Name     string
    Exchange string
}

// SearchBox is a custom widget that extends Entry to provide search functionality
type SearchBox struct {
    widget.Entry
    resultsContainer *fyne.Container
    onSelectStock    func(string)
    resultButtons    []*widget.Button
    selectedIndex    int
    isDropdownOpen   bool
    client           *data.AlphaVantageClient
    debounceTimer    *time.Timer
    overlay          *widget.PopUp
    window           fyne.Window
}

// NewSearchBox creates a new search box widget
func NewSearchBox(window fyne.Window, onSelectStock func(string)) *SearchBox {
    searchBox := &SearchBox{
        onSelectStock:    onSelectStock,
        resultsContainer: container.NewVBox(),
        resultButtons:    make([]*widget.Button, 0),
        selectedIndex:    -1,
        isDropdownOpen:   false,
        client:           data.NewAlphaVantageClient(),
        window:           window,
    }

    // Set up the entry
    searchBox.PlaceHolder = "Search stocks... (press / to focus)"
    searchBox.ActionItem = widget.NewIcon(theme.SearchIcon())

    // Set up event handlers
    searchBox.OnChanged = func(text string) {
        // Debounce the search to avoid too many API calls
        if searchBox.debounceTimer != nil {
            searchBox.debounceTimer.Stop()
        }
        
        searchBox.debounceTimer = time.AfterFunc(300*time.Millisecond, func() {
            // Only search if we have at least 2 characters
            if len(text) >= 2 {
                searchBox.Search(text)
            } else {
                searchBox.CloseDropdown()
            }
        })
    }

    searchBox.OnSubmitted = func(text string) {
        // If dropdown is open and something is selected, use that
        if searchBox.isDropdownOpen && searchBox.selectedIndex >= 0 && searchBox.selectedIndex < len(searchBox.resultButtons) {
            searchBox.resultButtons[searchBox.selectedIndex].OnTapped()
        } else if len(text) > 0 {
            // Otherwise just search with the text
            searchBox.Search(text)
        }
    }

    return searchBox
}

// Search performs a search and updates the dropdown
func (s *SearchBox) Search(query string) {
    // First try local search
    localResults := data.SearchSymbols(query)
    
    // If local search doesn't yield enough results, try Alpha Vantage
    if len(localResults) < 5 && s.client.APIKey != "" {
        go s.performAPISearch(query)
        return
    }
    
    // Convert local results to search results
    results := make([]SearchResult, 0, len(localResults))
    for _, stock := range localResults {
        results = append(results, SearchResult{
            Symbol:   stock.Symbol,
            Name:     stock.Name,
            Exchange: stock.Exchange,
        })
    }
    
    // Update the dropdown with results
    s.UpdateDropdown(results)
}

// performAPISearch searches using the Alpha Vantage API
func (s *SearchBox) performAPISearch(query string) {
    apiResults, err := s.client.SearchSymbols(query)
    if err != nil {
        return
    }
    
    // Convert API results to search results (limit to 5)
    results := make([]SearchResult, 0, 5)
    for i, result := range apiResults {
        if i >= 5 {
            break
        }
        results = append(results, SearchResult{
            Symbol:   result.Symbol,
            Name:     result.Name,
            Exchange: result.Region,
        })
    }
    
    // Update the dropdown on the main thread
    fyne.CurrentApp().Driver().CanvasForObject(s).Refresh(s)
    s.UpdateDropdown(results)
}

// UpdateDropdown updates the dropdown with the given search results
func (s *SearchBox) UpdateDropdown(results []SearchResult) {
    // Clear previous results
    s.resultsContainer.Objects = nil
    s.resultButtons = nil
    
    if len(results) == 0 {
        s.CloseDropdown()
        return
    }
    
    // Add new results
    for _, result := range results {
        btn := widget.NewButton(result.Symbol+" | "+result.Name, nil)
        btn.Alignment = widget.ButtonAlignLeading
        
        // Create a copy of the symbol for the closure
        symbol := result.Symbol
        btn.OnTapped = func() {
            s.onSelectStock(symbol)
            s.SetText("") // Clear the search box after selection
            s.CloseDropdown()
        }
        
        s.resultButtons = append(s.resultButtons, btn)
        s.resultsContainer.Add(btn)
    }
    
    s.selectedIndex = -1
    s.OpenDropdown()
}

// OpenDropdown opens the dropdown
func (s *SearchBox) OpenDropdown() {
    s.isDropdownOpen = true
    
    // Close existing popup if it exists
    if s.overlay != nil {
        s.overlay.Hide()
    }
    
    // Create a background for the results
    bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 240})
    
    // Limit to max 5 items
    displayResults := s.resultsContainer
    if len(s.resultButtons) > 5 {
        // Create a new container with only the first 5 items
        limitedContainer := container.NewVBox()
        for i := 0; i < 5; i++ {
            limitedContainer.Add(s.resultButtons[i])
        }
        displayResults = limitedContainer
    }
    
    resultBox := container.NewStack(
        bg,
        container.NewPadded(displayResults),
    )
    
    // Calculate position and show popup
    pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(s)
    s.overlay = widget.NewPopUp(resultBox, s.window.Canvas())
    s.overlay.ShowAtPosition(fyne.NewPos(pos.X, pos.Y+s.Size().Height))
    
    // Set minimum width to match search box width
    resultBox.Resize(fyne.NewSize(s.Size().Width, displayResults.MinSize().Height))
}

// CloseDropdown closes the dropdown
func (s *SearchBox) CloseDropdown() {
    s.isDropdownOpen = false
    if s.overlay != nil {
        s.overlay.Hide()
        s.overlay = nil
    }
}

// KeyDown handles key down events
func (s *SearchBox) KeyDown(key *fyne.KeyEvent) {
    switch key.Name {
    case fyne.KeyDown:
        if s.isDropdownOpen {
            s.selectedIndex = (s.selectedIndex + 1) % len(s.resultButtons)
            s.updateSelectedItem()
        }
    case fyne.KeyUp:
        if s.isDropdownOpen {
            s.selectedIndex = (s.selectedIndex - 1 + len(s.resultButtons)) % len(s.resultButtons)
            s.updateSelectedItem()
        }
    case fyne.KeyEscape:
        s.CloseDropdown()
    default:
        s.Entry.KeyDown(key)
    }
}

// updateSelectedItem updates the visual state of the selected item
func (s *SearchBox) updateSelectedItem() {
    for i, btn := range s.resultButtons {
        if i == s.selectedIndex {
            btn.Importance = widget.HighImportance
        } else {
            btn.Importance = widget.MediumImportance
        }
        btn.Refresh()
    }
}

// CreateCompactHeader creates a compact header with logo, search, and profile
func CreateCompactHeader(window fyne.Window, onSelectStock func(string), profile *models.Profile, onProfileSelect func(string)) *fyne.Container {
    // App title/logo on the left
    title := canvas.NewText("Moose Market", color.NRGBA{R: 76, G: 175, B: 80, A: 255})
    title.TextSize = 20
    title.TextStyle = fyne.TextStyle{Bold: true}
    
    // Create search box with popup behavior in the middle
    searchBox := NewSearchBox(window, onSelectStock)
    
    // Set up keyboard shortcut for / key
    if desktopCanvas, ok := window.Canvas().(desktop.Canvas); ok {
        desktopCanvas.SetOnKeyDown(func(event *fyne.KeyEvent) {
            if event.Name == fyne.KeySlash && window.Canvas().Focused() != searchBox {
                window.Canvas().Focus(searchBox)
                searchBox.SetText("")
            }
        })
    }
    
    // Create profile and settings buttons on the right
    profileButton := widget.NewButtonWithIcon("Default Profile", theme.AccountIcon(), func() {
        if onProfileSelect != nil {
            onProfileSelect(profile.Name)
        }
    })
    
    // Settings button
    settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), nil)
    
    // Create right side buttons container
    rightContainer := container.NewHBox(profileButton, settingsButton)
    
    // Create spacer to push title to left and buttons to right
    spacer1 := canvas.NewRectangle(color.Transparent)
    spacer1.SetMinSize(fyne.NewSize(400, 1))
    spacer2 := canvas.NewRectangle(color.Transparent)
    spacer2.SetMinSize(fyne.NewSize(400, 1))
    
    // Create a container with the title on the left, search in the middle, and buttons on the right
    topBar := container.NewHBox(
        title,  
        spacer1,
        searchBox,
        spacer2,
        rightContainer,
    )
    
    // Force search box size
    searchBox.Resize(fyne.NewSize(400, searchBox.Size().Height))
    
    // Final container
    return container.NewBorder(nil, nil, nil, nil, topBar)
}