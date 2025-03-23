# Moose Market

A comprehensive portfolio manager to facilitate trade and market readability across North American stock exchanges.

![Moose Market Screenshot](https://github.com/frederikblais/Moose-Market/raw/main/docs/screenshots/dashboard.png)

## Features

- **Multi-profile support**: Create and manage different investment profiles
- **Portfolio tracking**: Monitor multiple account types (TFSA, RRSP, FHSA, etc.)
- **Interactive charts**: View candlestick charts with multiple timeframes
- **Customizable watchlists**: Create and organize stock watchlists with real-time updates
- **Market heatmap**: Visualize market performance with color-coded tiles
- **Cross-platform**: Works on Windows, macOS, and Linux

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/frederikblais/Moose-Market.git
   cd Moose-Market
   ```

2. Build the application:
   ```bash
   go build -o moosemarket ./cmd/moosemarket
   ```

3. Run the application:
   ```bash
   ./moosemarket
   ```

## Usage Guide

### Creating a Profile

1. Launch Moose Market
2. Use the profile dropdown to select "Manage Profiles"
3. Click "Add Profile" and enter a name
4. Your new profile will be created with default settings

### Adding Stocks to Watchlist

1. Press `/` to focus the search bar
2. Type a stock symbol or company name
3. Select the stock to view its chart
4. Click "Add to Watchlist" to add it to your current watchlist

### Creating Watchlists

1. Click the "+" button in the watchlist panel
2. Enter a name for your new watchlist
3. Click "Create"
4. Your new empty watchlist will be available in the tabs

### Using the Chart

- Select different timeframes using the dropdown below the chart
- View price information including open, high, low, and close
- Charts automatically update at your configured refresh interval

## Development

### Project Structure

```
Moose-Market/
├── cmd/
│   └── moosemarket/      # Application entry point
├── internal/
│   ├── data/             # Data management and storage
│   ├── models/           # Data structures
│   └── ui/               # User interface components
├── docs/                 # Documentation
├── README.md             # This file
└── LICENSE               # MIT License
```

### Building from Source

```bash
go build ./cmd/moosemarket
```

### Running Tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Fyne](https://fyne.io/) - The UI toolkit used
- [Go](https://golang.org/) - The programming language

## Contact

Frederik Blais - [@frederikblais](https://github.com/frederikblais)

Project Link: [https://github.com/frederikblais/Moose-Market](https://github.com/frederikblais/Moose-Market)