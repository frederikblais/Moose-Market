package components

import (
    "image/color"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/theme"
)

// CustomTheme implements a dark theme for the application
type CustomTheme struct{}

// Color returns the color for the specified ThemeColorName
func (m CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNameBackground:
        return color.NRGBA{R: 30, G: 30, B: 30, A: 255}
    case theme.ColorNameButton:
        return color.NRGBA{R: 51, G: 51, B: 51, A: 255}
    case theme.ColorNameDisabled:
        return color.NRGBA{R: 50, G: 50, B: 50, A: 255}
    case theme.ColorNameDisabledButton:
        return color.NRGBA{R: 40, G: 40, B: 40, A: 255}
    case theme.ColorNameForeground:
        return color.NRGBA{R: 240, G: 240, B: 240, A: 255}
    case theme.ColorNameHover:
        return color.NRGBA{R: 60, G: 60, B: 60, A: 255}
    case theme.ColorNamePlaceHolder:
        return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
    case theme.ColorNamePressed:
        return color.NRGBA{R: 70, G: 70, B: 70, A: 255}
    case theme.ColorNamePrimary:
        return color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Green
    case theme.ColorNameScrollBar:
        return color.NRGBA{R: 80, G: 80, B: 80, A: 255}
    case theme.ColorNameShadow:
        return color.NRGBA{R: 0, G: 0, B: 0, A: 80}
    default:
        return theme.DefaultTheme().Color(name, variant)
    }
}

// Font returns the font resource for the specified TextStyle
func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
    return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the specified IconName
func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
    return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the specified ThemeSizeName
func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
    switch name {
    case theme.SizeNamePadding:
        return 4
    case theme.SizeNameInlineIcon:
        return 20
    case theme.SizeNameScrollBar:
        return 10
    case theme.SizeNameScrollBarSmall:
        return 5
    case theme.SizeNameText:
        return 14
    case theme.SizeNameHeadingText:
        return 24
    case theme.SizeNameSubHeadingText:
        return 18
    case theme.SizeNameCaptionText:
        return 11
    case theme.SizeNameInputBorder:
        return 2
    default:
        return theme.DefaultTheme().Size(name)
    }
}