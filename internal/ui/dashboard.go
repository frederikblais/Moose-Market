package ui

import (
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
    "github.com/frederikblais/Moose-Market/internal/ui/components"
)

type WatchlistInterface interface {
    GetContainer() *fyne.Container
    LoadWatchlists()
    LoadWatchlistItems()
    AddSymbol(symbol string)
}

// Dashboard represents the main UI of the application
type Dashboard struct {
    app               fyne.App
    window            fyne.Window
    header            *fyne.Container
    chartContainer    *components.ChartContainer
    watchlistContainer WatchlistInterface
    heatmapContainer  *components.HeatmapContainer
    activeProfile     *models.Profile
}

// NewDashboard creates a new dashboard UI
func NewDashboard(app fyne.App, window fyne.Window) *Dashboard {
    // Initialize the dashboard
    dashboard := &Dashboard{
        app:    app,
        window: window,
    }

    // Set window properties
    window.SetTitle("Moose Market")
    window.Resize(fyne.NewSize(1280, 800))

    // Set custom theme
    app.Settings().SetTheme(&components.CustomTheme{})

    return dashboard
}

// Setup initializes the dashboard UI
func (d *Dashboard) Setup() {
    // Initialize data storage
    if err := data.Initialize(); err != nil {
        dialog.ShowError(err, d.window)
        return
    }

    // Load existing profile or create default
    d.initializeProfile()

    // Create components
    d.createUI()

    // Set window content
    d.window.SetContent(d.createLayout())
}

// initializeProfile loads or creates a profile
func (d *Dashboard) initializeProfile() {
    // Try to load profiles
    profiles, err := data.GetProfiles()
    if err != nil || len(profiles) == 0 {
        // Create a default profile if none exists
        defaultProfile, err := data.CreateProfile("Default Profile")
        if err != nil {
            dialog.ShowError(err, d.window)
            return
        }
        
        d.activeProfile = defaultProfile
    } else {
        // Load the first profile
        d.activeProfile = &profiles[0]
    }

    // Set as active
    data.SetActiveProfile(d.activeProfile)
}

// createUI creates all UI components
func (d *Dashboard) createUI() {
    // Create the compact header with improved search functionality
    d.header = components.CreateCompactHeader(
        d.window, 
        func(symbol string) {
            // Handle stock selection from search
            d.chartContainer.LoadChart(symbol)
        },
        d.activeProfile,
        func(profileName string) {
            // For the header's profile button, we'll show the profile manager
            d.showProfileManager()
        },
    )

    // Create chart container
    d.chartContainer = components.CreateChartContainer(func(symbol string) {
        // Add to watchlist callback
        d.watchlistContainer.AddSymbol(symbol)
    })

    // Create watchlist container
    d.watchlistContainer = components.CreateMinimalWatchlistContainer(
        func(symbol string) {
            // Select stock callback
            d.chartContainer.LoadChart(symbol)
        },
        func(watchlistID string) {
            // Watchlist changed callback
            d.heatmapContainer.SetWatchlist(watchlistID)
        },
    )

    // Create heatmap container
    d.heatmapContainer = components.CreateHeatmapContainer(func(symbol string) {
        // Select stock callback
        d.chartContainer.LoadChart(symbol)
    })

    // Load watchlists from the active profile
    d.watchlistContainer.LoadWatchlists()

    // Set up periodic refresh
    go d.setupPeriodicRefresh()
}

// createLayout creates the main layout for the dashboard
func (d *Dashboard) createLayout() fyne.CanvasObject {
    // Create the right panel with watchlist and heatmap
    rightPanel := container.NewVSplit(
        d.watchlistContainer.GetContainer(),
        d.heatmapContainer.GetContainer(),
    )
    rightPanel.SetOffset(0.6) // 60% watchlist, 40% heatmap

    // Create the main content area with chart and right panel
    mainContent := container.NewHSplit(
        d.chartContainer.GetContainer(),
        rightPanel,
    )
    mainContent.SetOffset(0.7) // 70% chart, 30% right panel

    // Put everything together with a clean layout
    return container.NewBorder(
        d.header,   // Top
        nil,        // Bottom
        nil,        // Left
        nil,        // Right
        mainContent, // Center
    )
}

// showProfileManager displays the profile management dialog
func (d *Dashboard) showProfileManager() {
    // Create a list of profiles
    profiles, _ := data.GetProfiles()
    var items []string
    for _, profile := range profiles {
        items = append(items, profile.Name)
    }

    // Create list widget
    list := widget.NewList(
        func() int {
            return len(items)
        },
        func() fyne.CanvasObject {
            return container.NewBorder(
                nil, nil, nil,
                widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
                widget.NewLabel("Template"),
            )
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*fyne.Container).Objects[1].(*widget.Label).SetText(items[i])
            o.(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
                // Delete profile confirmation
                dialog.ShowConfirm("Delete Profile", 
                    "Are you sure you want to delete "+items[i]+"?", 
                    func(confirm bool) {
                        if confirm {
                            data.DeleteProfile(profiles[i].ID)
                            d.refreshProfileList()
                        }
                    }, 
                    d.window)
            }
        },
    )

    // Create add profile button
    addButton := widget.NewButton("Add Profile", func() {
        // Show dialog to add a new profile
        input := widget.NewEntry()
        input.SetPlaceHolder("Profile Name")
        
        dialog.ShowCustomConfirm("New Profile", "Create", "Cancel", 
            container.NewVBox(
                widget.NewLabel("Enter profile name:"),
                input,
            ),
            func(confirm bool) {
                if confirm && input.Text != "" {
                    _, err := data.CreateProfile(input.Text)
                    if err != nil {
                        dialog.ShowError(err, d.window)
                    } else {
                        d.refreshProfileList()
                    }
                }
            },
            d.window)
    })

    // Show profile manager dialog
    dialog.ShowCustom("Manage Profiles", "Close", 
        container.NewBorder(nil, addButton, nil, nil, list),
        d.window)
}

// refreshProfileList refreshes the profile dropdown
func (d *Dashboard) refreshProfileList() {
    // Reload profiles
    profiles, _ := data.GetProfiles()
    
    // Update active profile if it was deleted
    found := false
    for _, profile := range profiles {
        if profile.ID == d.activeProfile.ID {
            found = true
            break
        }
    }
    
    if !found && len(profiles) > 0 {
        d.activeProfile = &profiles[0]
        data.SetActiveProfile(d.activeProfile)
        d.watchlistContainer.LoadWatchlists()
    }
    
    // Refresh UI
    d.window.Content().Refresh()
}

// setupPeriodicRefresh refreshes stock data periodically
func (d *Dashboard) setupPeriodicRefresh() {
    for {
        // Sleep for the refresh interval (default 60 seconds)
        time.Sleep(time.Duration(d.activeProfile.Settings.RefreshInterval) * time.Second)
        
        // Refresh the watchlist
        d.watchlistContainer.LoadWatchlistItems()
        
        // Refresh the heatmap
        d.heatmapContainer.RefreshHeatmap()
    }
}