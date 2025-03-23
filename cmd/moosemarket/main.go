// File: cmd/moosemarket/main.go
package main

import (
    "fyne.io/fyne/v2/app"
    
    "github.com/frederikblais/Moose-Market/internal/ui"
)

func main() {
    // Create the Fyne application
    a := app.New()
    w := a.NewWindow("Moose Market")
    
    // Create and set up the dashboard
    dashboard := ui.NewDashboard(a, w)
    dashboard.Setup()
    
    // Run the application
    w.ShowAndRun()
}