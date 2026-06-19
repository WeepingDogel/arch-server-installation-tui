package tui

import "github.com/charmbracelet/lipgloss"

// Theme contains all style definitions for the TUI.
// Using Arch Linux brand colors with a dark theme.
var (
	// Colors
	ColorPrimary   = lipgloss.Color("#1793D1") // Arch cyan
	ColorSecondary = lipgloss.Color("#2E3440") // Dark background
	ColorAccent    = lipgloss.Color("#81A1C1") // Light blue-gray
	ColorSuccess   = lipgloss.Color("#A3BE8C") // Green
	ColorError     = lipgloss.Color("#BF616A") // Red
	ColorWarning   = lipgloss.Color("#EBCB8B") // Yellow
	ColorWhite     = lipgloss.Color("#ECEFF4") // Snow white
	ColorGray      = lipgloss.Color("#4C566A") // Gray
	ColorDark      = lipgloss.Color("#1E2229") // Darker background

	// Layout styles
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Background(ColorSecondary).
			Foreground(ColorWhite)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Padding(0, 1)

	ContentStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorWhite)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Padding(0, 1).
			MarginTop(1)

	// Navigation buttons
	ButtonStyle = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 3).
			MarginRight(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary)

	ButtonDisabledStyle = lipgloss.NewStyle().
				Background(ColorGray).
				Foreground(ColorDark).
				Padding(0, 3).
				MarginRight(1).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorGray)

	BackButtonStyle = lipgloss.NewStyle().
			Background(ColorDark).
			Foreground(ColorAccent).
			Padding(0, 3).
			MarginRight(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(0, 1).
			Width(40)

	InputFocusStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1).
			Width(40)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorWhite)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Background(ColorPrimary).
				Foreground(ColorWhite).
				Padding(0, 2).
				Bold(true)

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Italic(true).
			Padding(0, 2)

	// Error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true).
			Padding(0, 2)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			Padding(0, 2)

	// Progress bar
	ProgressBarEmpty = lipgloss.Color("#3B4252")
	ProgressBarFill  = ColorPrimary

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Padding(0, 1).
			Width(50).
			Align(lipgloss.Center)

	// Step indicator
	StepStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Padding(0, 1).
			MarginRight(1)

	StepActiveStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	StepDoneStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Padding(0, 1)

	// Box style for grouping
	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorGray).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
			Width(60)
)

// DividerText returns a centered divider with the given text.
func DividerText(text string) string {
	return DividerStyle.Render("─── " + text + " ───")
}
