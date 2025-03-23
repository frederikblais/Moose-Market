// File: internal/ui/components/chart.go
package components

import (
    "fmt"
    "image/color"
    "strconv"
    "time"
    "sort"
    "math"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    
    "github.com/frederikblais/Moose-Market/internal/data"
    "github.com/frederikblais/Moose-Market/internal/models"
)

// ChartContainer represents the chart area with controls
type ChartContainer struct {
    container          *fyne.Container
    symbol             string
    timeframe          string
    candleData         *models.CandleData
    chartCanvas        *canvas.Rectangle
    chartContent       *fyne.Container
    onAddToWatchlist   func(string)
    alphaVantageClient *data.AlphaVantageClient
    loadingIndicator   *widget.ProgressBarInfinite
    symbolInfoLabel    *widget.Label
    chartArea          *fyne.Container
    emptyText          *canvas.Text
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
    
    // Initialize chart content container
    chartContent := container.NewWithoutLayout()
    
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
        chartContent,
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
        container:          mainContainer,
        chartCanvas:        chartPlaceholder,
        chartContent:       chartContent,
        timeframe:          "1d",
        onAddToWatchlist:   onAddToWatchlist,
        alphaVantageClient: data.NewAlphaVantageClient(),
        loadingIndicator:   loadingIndicator,
        symbolInfoLabel:    symbolInfoLabel,
        chartArea:          chartArea,
        emptyText:          emptyText,
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
    c.emptyText.Hide()
    
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
    
     // Find the control bar at the bottom of the container
     if controlBar, ok := c.container.Objects[1].(*fyne.Container); ok {
        // Find the button (should be the second element in the control bar)
        if len(controlBar.Objects) > 1 {
            if button, ok := controlBar.Objects[1].(*widget.Button); ok {
                button.Enable()
            }
        }
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
            c.loadingIndicator.Hide()
            canvas.Refresh(c.loadingIndicator)
        }
    }
    
    if err != nil {
        // Show error message
        errorText := canvas.NewText("Error loading chart data", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
        errorText.Alignment = fyne.TextAlignCenter
        errorText.TextSize = 16
        
        c.chartContent.Objects = nil
        c.chartContent.Add(container.NewCenter(errorText))
        c.chartContent.Refresh()
        return
    }
    
    c.candleData = candleData
    
    // Get the current size of the chart area
    chartSize := c.chartCanvas.Size()
    if chartSize.Width < 10 || chartSize.Height < 10 {
        // Set a default size if the chart area is too small
        chartSize = fyne.NewSize(600, 400)
    }
    
    // Create a new chart visualization
    newChartContent := createImprovedCandleChart(candleData, chartSize)
    
    // Update the chart content
    c.chartContent.Objects = nil
    c.chartContent.Add(newChartContent)
    
    // Refresh the container
    c.chartContent.Refresh()
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

// createImprovedCandleChart creates a better implementation of the candle chart
func createImprovedCandleChart(data *models.CandleData, size fyne.Size) fyne.CanvasObject {
    chartContainer := container.NewWithoutLayout()
    
    if data == nil || len(data.Candles) == 0 {
        noDataText := canvas.NewText("No chart data available", color.White)
        noDataText.Alignment = fyne.TextAlignCenter
        noDataText.TextSize = 16
        noDataText.Move(fyne.NewPos(size.Width/2-100, size.Height/2-10))
        chartContainer.Add(noDataText)
        return chartContainer
    }
    
    // Find min/max values for scaling
    var min, max float64
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
    
    // Add some buffer
    priceRange := max - min
    min -= priceRange * 0.05
    max += priceRange * 0.05
    
    // Chart dimensions
    margin := float32(40)
    chartWidth := size.Width - margin*2
    chartHeight := size.Height - margin*2
    
    // Background
    bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
    bg.Resize(size)
    bg.Move(fyne.NewPos(0, 0))
    chartContainer.Add(bg)
    
    // Draw grid lines
    gridSteps := 5
    for i := 0; i <= gridSteps; i++ {
        y := margin + chartHeight - (chartHeight * float32(i) / float32(gridSteps))
        
        // Horizontal grid line
        line := canvas.NewLine(color.NRGBA{R: 60, G: 60, B: 60, A: 255})
        line.StrokeWidth = 1
        line.Move(fyne.NewPos(margin, y))
        line.Resize(fyne.NewSize(chartWidth, 1))
        chartContainer.Add(line)
        
        // Price label
        price := min + (max-min)*float64(i)/float64(gridSteps)
        priceLabel := canvas.NewText(fmt.Sprintf("$%.2f", price), color.NRGBA{R: 200, G: 200, B: 200, A: 255})
        priceLabel.TextSize = 12
        priceLabel.Move(fyne.NewPos(margin-35, y-8))
        chartContainer.Add(priceLabel)
    }
    
    // Chart title
    title := canvas.NewText(fmt.Sprintf("%s - %s Chart", data.Symbol, data.Timeframe), color.NRGBA{R: 220, G: 220, B: 220, A: 255})
    title.TextSize = 16
    title.TextStyle = fyne.TextStyle{Bold: true}
    title.Move(fyne.NewPos(margin, 10))
    chartContainer.Add(title)
    
    // Number of candles to display (limit to what can reasonably fit)
    maxCandles := minInt(len(data.Candles), 50)
    
    // Use only the most recent candles if we have too many
    startIdx := len(data.Candles) - maxCandles
    if startIdx < 0 {
        startIdx = 0
    }
    
    displayCandles := data.Candles[startIdx:]
    
    // Draw candles
    candleSpacing := float32(2)
    candleWidth := (chartWidth / float32(maxCandles)) - candleSpacing
    
    for i, candle := range displayCandles {
        // Calculate x position
        x := margin + float32(i)*(candleWidth+candleSpacing)
        
        // Scale prices to chart height
        highY := margin + chartHeight - float32((candle.High-min)/(max-min))*chartHeight
        lowY := margin + chartHeight - float32((candle.Low-min)/(max-min))*chartHeight
        openY := margin + chartHeight - float32((candle.Open-min)/(max-min))*chartHeight
        closeY := margin + chartHeight - float32((candle.Close-min)/(max-min))*chartHeight
        
        // Draw the wick
        wick := canvas.NewLine(color.White)
        wick.StrokeWidth = 1
        wick.Position1 = fyne.NewPos(x+candleWidth/2, highY)
        wick.Position2 = fyne.NewPos(x+candleWidth/2, lowY)
        chartContainer.Add(wick)
        
        // Draw the candle body
        var bodyColor color.Color
        if candle.Close >= candle.Open {
            bodyColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green for up
        } else {
            bodyColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255} // Red for down
        }
        
        body := canvas.NewRectangle(bodyColor)
        
        bodyTop := fyne.Min(openY, closeY)
        bodyHeight := fyne.Max(float32(1), float32(math.Abs(float64(closeY-openY))))
        
        body.Move(fyne.NewPos(x, bodyTop))
        body.Resize(fyne.NewSize(candleWidth, bodyHeight))
        
        chartContainer.Add(body)
        
        // Date labels for some candles
        if i == 0 || i == len(displayCandles)-1 || i%(maxCandles/5) == 0 {
            dateText := canvas.NewText(formatShortDate(candle.Time, data.Timeframe), color.NRGBA{R: 180, G: 180, B: 180, A: 255})
            dateText.TextSize = 10
            dateText.Move(fyne.NewPos(x, margin+chartHeight+5))
            chartContainer.Add(dateText)
        }
    }
    
    // Add some information text
    if len(data.Candles) > 0 {
        lastCandle := data.Candles[len(data.Candles)-1]
        infoText := fmt.Sprintf("O: %.2f  H: %.2f  L: %.2f  C: %.2f", 
            lastCandle.Open, lastCandle.High, lastCandle.Low, lastCandle.Close)
        
        info := canvas.NewText(infoText, color.NRGBA{R: 220, G: 220, B: 220, A: 255})
        info.TextSize = 14
        info.Move(fyne.NewPos(margin+chartWidth-250, 10))
        chartContainer.Add(info)
    }
    
    return chartContainer
}

// Helper function for date formatting
func formatShortDate(t time.Time, timeframe string) string {
    switch timeframe {
    case "1min", "5min", "15min", "30min", "60min":
        return t.Format("15:04")
    case "1d":
        return t.Format("01/02")
    case "1w":
        return t.Format("01/02")
    default:
        return t.Format("01/02")
    }
}

// Helper function to get minimum of two integers
func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}