package components

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/theme"
    "image/color"
)

// CreateHeader creates the top header with app title and search bar
func CreateHeader(onSearch func(string)) *fyne.Container {
    // App title
    title := canvas.NewText("Moose Market", color.NRGBA{R: 76, G: 175, B: 80, A: 255})
    title.TextSize = 24
    title.TextStyle = fyne.TextStyle{Bold: true}
    
    // Search entry
    searchEntry := widget.NewEntry()
    searchEntry.PlaceHolder = "Search stocks... (press / to focus)"
    searchEntry.OnSubmitted = onSearch
    
    // Set up keyboard shortcut for search
    searchEntry.ActionItem = widget.NewIcon(theme.SearchIcon())
    
    // Create the header container
    header := container.NewBorder(
        nil, nil, title, nil,
        container.NewPadded(searchEntry))
    
    return header
}