// File: internal/ui/components/chart.go
package components

import (
    "fmt"
    "image/color"
    "strconv"
    "time"
    "sort"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    
    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// ChartContainer represents the chart area with controls
type ChartContainer struct {
    container         *fyne.Container
    symbol            string
    timeframe         string
    candleData        *models.CandleData
    chartCanvas       *canvas.Rectangle
    onAddToWatchlist  func(string)
    alphaVantageClient *data.AlphaVantageClient
    loadingIndicator  *widget.ProgressBarInfinite
    symbolInfoLabel   *widget.Label
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
    
    // Loading indicator
    loadingIndicator := widget.NewProgressBarInfinite()
    loadingIndicator.Hide()
    
    // Symbol info label
    symbolInfoLabel := widget.NewLabel("")
    symbolInfoLabel.Hide()
    
    // Timeframe selectors
    timeframeOptions := []string{"5min", "15min", "30min", "60min", "1d", "1w"}
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
    
    // Create chart area with loading indicator and placeholder
    chartArea := container.NewStack(
        chartPlaceholder,
        container.NewCenter(emptyText),
        loadingIndicator,
    )
    
    // Combine chart and controls
    mainContainer := container.NewBorder(
        symbolInfoLabel,
        controlBar,
        nil,
        nil,
        chartArea,
    )
    
    chartContainer := &ChartContainer{
        container:         mainContainer,
        chartCanvas:       chartPlaceholder,
        timeframe:         "1d",
        onAddToWatchlist:  onAddToWatchlist,
        alphaVantageClient: data.NewAlphaVantageClient(),
        loadingIndicator:   loadingIndicator,
        symbolInfoLabel:    symbolInfoLabel,
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
    
    // Show loading indicator
    c.loadingIndicator.Show()
    
    // First try to get quote data
    go func() {
        // Try to get quote data first
        if c.alphaVantageClient.APIKey != "" {
            quote, err := c.alphaVantageClient.GetQuote(symbol)
            if err == nil {
                // Update symbol info label
                priceStr := quote.Price
                changeStr := quote.Change
                changePctStr := quote.ChangePercent
                
                info := fmt.Sprintf("%s - $%s | %s (%s)", 
                    symbol, priceStr, changeStr, changePctStr)
                
                // Update UI on the main thread
                updateLabelTextSafely(c.symbolInfoLabel, info)
            }
        }
        
        // Now load the chart data
        c.loadChartData(symbol)
    }()
    
    // Enable the add to watchlist button
    if button, ok := c.container.Objects[1].(*fyne.Container).Objects[1].(*widget.Button); ok {
        button.Enable()
    }
}

// loadChartData loads candle data for the chart
func (c *ChartContainer) loadChartData(symbol string) {
    var candleData *models.CandleData
    var err error
    
    // If Alpha Vantage API key is set, try to use it
    if c.alphaVantageClient.APIKey != "" {
        candleData, err = c.loadAlphaVantageData(symbol)
    } else {
        // Fall back to mock data
        candleData, err = data.GetCandleData(symbol, c.timeframe, 100)
    }
    
    // Update the UI on the main thread
    // Hide loading indicator
    if c.loadingIndicator != nil {
        // We need to use the canvas to refresh UI elements on the main thread
        canvas := fyne.CurrentApp().Driver().CanvasForObject(c.loadingIndicator)
        if canvas != nil {
            canvas.Refresh(c.loadingIndicator)
            c.loadingIndicator.Hide()
        }
    }
    
    if err != nil {
        // Show error message
        errorText := canvas.NewText("Error loading chart data", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
        errorText.Alignment = fyne.TextAlignCenter
        errorText.TextSize = 16
        
        // Find the stack container and update it
        if stack, ok := c.container.Objects[2].(*fyne.Container); ok {
            stack.Objects[1] = container.NewCenter(errorText)
            
            // Refresh the container
            canvas := fyne.CurrentApp().Driver().CanvasForObject(stack)
            if canvas != nil {
                canvas.Refresh(stack)
            }
        }
        
        return
    }
    
    c.candleData = candleData
    
    // For now, we'll just display a basic chart representation
    chartContent := createSimpleCandleChart(candleData, c.chartCanvas.Size())
    
    // Find the stack container and update it
    if stack, ok := c.container.Objects[2].(*fyne.Container); ok {
        stack.Objects[1] = chartContent
        
        // Refresh the container
        canvas := fyne.CurrentApp().Driver().CanvasForObject(stack)
        if canvas != nil {
            canvas.Refresh(stack)
        }
    }
}

// loadAlphaVantageData loads candle data from Alpha Vantage
func (c *ChartContainer) loadAlphaVantageData(symbol string) (*models.CandleData, error) {
    // Choose the appropriate API call based on timeframe
    var seriesData map[string]data.TimeSeriesData
    var err error
    
    if c.timeframe == "1d" {
        // Daily data
        seriesData, err = c.alphaVantageClient.GetDailyTimeSeries(symbol, true)
    } else {
        // Intraday data
        seriesData, err = c.alphaVantageClient.GetIntradayTimeSeries(symbol, c.timeframe, true)
    }
    
    if err != nil {
        return nil, err
    }
    
    // Convert to CandleData format
    candleData := &models.CandleData{
        Symbol:    symbol,
        Timeframe: c.timeframe,
        Candles:   make([]models.CandleStick, 0, len(seriesData)),
    }
    
    // Process the data
    for dateStr, timeSeries := range seriesData {
        // Parse date
        date, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            // Try intraday format
            date, err = time.Parse("2006-01-02 15:04:05", dateStr)
            if err != nil {
                continue // Skip dates we can't parse
            }
        }
        
        // Parse numeric values
        open, _ := strconv.ParseFloat(timeSeries.Open, 64)
        high, _ := strconv.ParseFloat(timeSeries.High, 64)
        low, _ := strconv.ParseFloat(timeSeries.Low, 64)
        close, _ := strconv.ParseFloat(timeSeries.Close, 64)
        volume, _ := strconv.ParseInt(timeSeries.Volume, 10, 64)
        
        // Create candle stick
        candle := models.CandleStick{
            Time:   date,
            Open:   open,
            High:   high,
            Low:    low,
            Close:  close,
            Volume: volume,
        }
        
        candleData.Candles = append(candleData.Candles, candle)
    }
    
    // Sort candles by time
    sort.Slice(candleData.Candles, func(i, j int) bool {
        return candleData.Candles[i].Time.Before(candleData.Candles[j].Time)
    })
    
    return candleData, nil
}

// Helper function to update a label's text safely from a goroutine
func updateLabelTextSafely(label *widget.Label, text string) {
    if label == nil {
        return
    }
    
    canvas := fyne.CurrentApp().Driver().CanvasForObject(label)
    if canvas == nil {
        return
    }
    
    label.SetText(text)
    label.Show()
    canvas.Refresh(label)
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
    case "1min", "5min", "15min", "30min", "60min":
        return t.Format("15:04")
    case "1d":
        return t.Format("Jan 02")
    case "1w":
        return t.Format("Jan 02, 2006")
    default:
        return t.Format("Jan 02")
    }
}