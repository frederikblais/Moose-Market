// File: internal/ui/components/chart.go
package components

import (
    "fmt"
    "image/color"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    
    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// ChartContainer represents the chart area with controls
type ChartContainer struct {
    container   *fyne.Container
    symbol      string
    timeframe   string
    candleData  *models.CandleData
    chartCanvas *canvas.Rectangle
    onAddToWatchlist func(string)
}

// CreateChartContainer creates the stock chart container
func CreateChartContainer(onAddToWatchlist func(string)) *ChartContainer {
    // Create placeholder for chart
    chartPlaceholder := canvas.NewRectangle(color.NRGBA{R: 40, G: 40, B: 40, A: 255})
    chartPlaceholder.SetMinSize(fyne.NewSize(600, 400))
    
    // Create text overlay for empty state
    emptyText := canvas.NewText("Select a stock to display chart", color.White)
    emptyText.Alignment = fyne.TextAlignCenter
    emptyText.TextSize = 18
    
    // Timeframe selectors
    timeframeOptions := []string{"1m", "5m", "15m", "1h", "1d", "1w"}
    timeframeSelect := widget.NewSelect(timeframeOptions, func(selected string) {
        // Will be implemented in the returned struct
    })
    timeframeSelect.Selected = "1d"
    
    // Add to watchlist button
    addButton := widget.NewButton("Add to Watchlist", func() {
        // Will be implemented in the returned struct
    })
    addButton.Disable() // Disabled until a stock is selected
    
    // Control bar
    controlBar := container.NewHBox(
        timeframeSelect,
        addButton,
    )
    
    // Combine chart and controls
    chartArea := container.NewBorder(nil, controlBar, nil, nil, container.NewCenter(emptyText, chartPlaceholder))
    
    chartContainer := &ChartContainer{
        container:   chartArea,
        chartCanvas: chartPlaceholder,
        timeframe:   "1d",
        onAddToWatchlist: onAddToWatchlist,
    }
    
    // Set up the timeframe callback
    timeframeSelect.OnChanged = func(selected string) {
        chartContainer.timeframe = selected
        if chartContainer.symbol != "" {
            chartContainer.LoadChart(chartContainer.symbol)
        }
    }
    
    // Set up add to watchlist callback
    addButton.OnTapped = func() {
        if chartContainer.symbol != "" && chartContainer.onAddToWatchlist != nil {
            chartContainer.onAddToWatchlist(chartContainer.symbol)
        }
    }
    
    return chartContainer
}

// GetContainer returns the container for the chart
func (c *ChartContainer) GetContainer() *fyne.Container {
    return c.container
}

// LoadChart loads the chart for a specific symbol
func (c *ChartContainer) LoadChart(symbol string) {
    c.symbol = symbol
    
    // Fetch candle data for the selected symbol and timeframe
    candleData, err := data.GetCandleData(symbol, c.timeframe, 100)
    if err != nil {
        // Show error message
        errorText := canvas.NewText("Error loading chart data", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
        errorText.Alignment = fyne.TextAlignCenter
        errorText.TextSize = 16
        c.container.Objects[0] = container.NewCenter(errorText)
        c.container.Refresh()
        return
    }
    
    c.candleData = candleData
    
    // For now, we'll just display a basic chart representation
    // In a real implementation, you would use a proper charting library
    chartTitle := canvas.NewText(symbol+" - "+c.timeframe, color.White)
    chartTitle.Alignment = fyne.TextAlignCenter
    chartTitle.TextSize = 18
    
    // Update UI
    chartContent := container.NewVBox(
        chartTitle,
        createSimpleCandleChart(candleData, c.chartCanvas.Size()),
    )
    
    c.container.Objects[0] = chartContent
    c.container.Refresh()
    
    // Enable the add to watchlist button
    if button, ok := c.container.Objects[1].(*fyne.Container).Objects[1].(*widget.Button); ok {
        button.Enable()
    }
}

// createSimpleCandleChart creates a simple representation of candles
// This is a placeholder - a real implementation would use a proper charting library
func createSimpleCandleChart(data *models.CandleData, size fyne.Size) fyne.CanvasObject {
    container := container.NewWithoutLayout()
    
    // This is just a placeholder implementation
    // In a real app, you would use a proper chart rendering library
    
    // Find min/max for scaling
    var min, max float64
    if len(data.Candles) > 0 {
        min = data.Candles[0].Low
        max = data.Candles[0].High
        
        for _, candle := range data.Candles {
            if candle.Low < min {
                min = candle.Low
            }
            if candle.High > max {
                max = candle.High
            }
        }
    }
    
    // Add some buffer
    priceRange := max - min
    min -= priceRange * 0.05
    max += priceRange * 0.05
    
    // Create candlesticks
    canvasWidth := size.Width
    canvasHeight := size.Height - 40 // Leave space for labels
    
    if len(data.Candles) > 0 {
        candleWidth := canvasWidth / float32(len(data.Candles)) * 0.8
        spacing := canvasWidth / float32(len(data.Candles)) * 0.2
        
        for i, candle := range data.Candles {
            x := float32(i) * (candleWidth + spacing)
            
            // Scale prices to canvas height
            highY := canvasHeight - float32((candle.High-min)/(max-min))*canvasHeight
            lowY := canvasHeight - float32((candle.Low-min)/(max-min))*canvasHeight
            openY := canvasHeight - float32((candle.Open-min)/(max-min))*canvasHeight
            closeY := canvasHeight - float32((candle.Close-min)/(max-min))*canvasHeight
            
            // Draw candle body
            var bodyColor color.Color
            if candle.Close >= candle.Open {
                bodyColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green for up
            } else {
                bodyColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // Red for down
            }
            
            body := canvas.NewRectangle(bodyColor)
            body.Move(fyne.NewPos(x, min32(openY, closeY)))
            body.Resize(fyne.NewSize(candleWidth, abs32(closeY-openY)+1))
            
            // Draw wick
            wick := canvas.NewLine(color.White)
            wick.StrokeWidth = 1
            wick.Position1 = fyne.NewPos(x+candleWidth/2, highY)
            wick.Position2 = fyne.NewPos(x+candleWidth/2, lowY)
            
            container.Add(wick)
            container.Add(body)
        }
        
        // Add price scale on the right side
        priceStep := (max - min) / 5
        for i := 0; i <= 5; i++ {
            price := min + priceStep*float64(i)
            y := canvasHeight - float32(float64(i)/5)*canvasHeight
            
            priceText := canvas.NewText(formatPrice(price), color.NRGBA{R: 180, G: 180, B: 180, A: 255})
            priceText.Alignment = fyne.TextAlignTrailing
            priceText.TextSize = 12
            priceText.Move(fyne.NewPos(canvasWidth-70, y-10))
            
            container.Add(priceText)
        }
        
        // Add date labels at bottom
        if len(data.Candles) > 0 {
            startTime := data.Candles[0].Time
            endTime := data.Candles[len(data.Candles)-1].Time
            
            startLabel := canvas.NewText(formatDate(startTime, data.Timeframe), color.NRGBA{R: 180, G: 180, B: 180, A: 255})
            startLabel.TextSize = 12
            startLabel.Move(fyne.NewPos(10, canvasHeight+10))
            
            endLabel := canvas.NewText(formatDate(endTime, data.Timeframe), color.NRGBA{R: 180, G: 180, B: 180, A: 255})
            endLabel.TextSize = 12
            endLabel.Alignment = fyne.TextAlignTrailing
            endLabel.Move(fyne.NewPos(canvasWidth-100, canvasHeight+10))
            
            container.Add(startLabel)
            container.Add(endLabel)
        }
    }
    
    return container
}

// Helper functions
func abs32(x float32) float32 {
    if x < 0 {
        return -x
    }
    return x
}

func min32(a, b float32) float32 {
    if a < b {
        return a
    }
    return b
}

func formatPrice(price float64) string {
    return "$" + widget.NewLabel(fmt.Sprintf("%.2f", price)).Text
}

func formatDate(t time.Time, timeframe string) string {
    switch timeframe {
    case "1m", "5m", "15m", "1h":
        return t.Format("15:04")
    case "1d":
        return t.Format("Jan 02")
    case "1w":
        return t.Format("Jan 02, 2006")
    default:
        return t.Format("Jan 02")
    }
}