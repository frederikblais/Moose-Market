package ui

import (
    "fmt"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

func NewMooseMarketApp() {
    // Create the Fyne application
    a := app.New()
    w := a.NewWindow("Moose Market")
    w.Resize(fyne.NewSize(800, 600))

    // Container that will hold stock list
    listContainer := container.NewVBox()

    // “Refresh” button to load new data
    refreshButton := widget.NewButton("Refresh Data", func() {
        stocks, err := data.FetchMarketData()
        if err != nil {
            fmt.Println("Error fetching data:", err)
            return
        }
        if err := data.SaveStocks(stocks); err != nil {
            fmt.Println("Error saving data:", err)
        }
        updateStockList(listContainer, stocks)
    })

    // On startup, try to load existing data
    savedStocks, err := data.LoadStocks()
    if err == nil && len(savedStocks) > 0 {
        updateStockList(listContainer, savedStocks)
    } else {
        // If no saved data, fetch new data
        newStocks, _ := data.FetchMarketData()
        _ = data.SaveStocks(newStocks)
        updateStockList(listContainer, newStocks)
    }

    // Put everything into a layout
    content := container.NewBorder(
        refreshButton, // top
        nil,           // bottom
        nil,           // left
        nil,           // right
        listContainer, // center
    )

    w.SetContent(content)
    w.ShowAndRun()
}

// Helper function to update the container with stock data
func updateStockList(listContainer *fyne.Container, stocks []models.Stock) {
    // Clear old items
    listContainer.Objects = listContainer.Objects[:0]

    for _, s := range stocks {
        priceStr := fmt.Sprintf("%.2f", s.Price)
        changeStr := fmt.Sprintf("%.2f", s.Change)

        label := widget.NewLabel(
            s.Symbol + " | " + s.Name + " | $" + priceStr + " | Δ " + changeStr,
        )
        listContainer.Add(label)
    }

    // Refresh the container so the UI updates
    listContainer.Refresh()
}
