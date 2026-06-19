package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Screen returns a rendered screen with header, content, and footer.
func Screen(step int, content, footer string) string {
	header := renderHeader(step)
	contentBox := lipgloss.NewStyle().
		Padding(0, 1).
		Render(content)

	footerBox := FooterStyle.Render(footer)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		"",
		contentBox,
		"",
		footerBox,
	)
}

// renderHeader creates the styled header with step indicator.
func renderHeader(step int) string {
	logo := MiniArchLogo()
	stepInd := StepIndicator(step, TotalSteps)
	stepName := lipgloss.NewStyle().
		Foreground(ColorWhite).
		Bold(true).
		Render(" " + StepName(step) + " ")

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Center,
		logo,
		stepInd,
		stepName,
	)

	return HeaderStyle.Render(headerContent)
}

// ContentBox wraps content in a styled box.
func ContentBox(content string, width int) string {
	return BoxStyle.Width(width).Render(content)
}

// FooterHelp returns a help text for navigation.
func FooterHelp() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		HelpStyle.Render("↑/↓ Navigate  "),
		HelpStyle.Render("Enter Select  "),
		HelpStyle.Render("Esc Back  "),
		HelpStyle.Render("Tab Next  "),
		HelpStyle.Render("Ctrl+C Quit"),
	)
}

// SimpleFooter returns a minimal footer.
func SimpleFooter() string {
	return HelpStyle.Render("↑/↓  •  Enter  •  Esc  •  Tab  •  Ctrl+C")
}

// ErrorBox renders an error message in a styled box.
func ErrorBox(err string) string {
	if err == "" {
		return ""
	}
	return ErrorStyle.Render("✗ " + err)
}

// SuccessBox renders a success message in a styled box.
func SuccessBox(msg string) string {
	return SuccessStyle.Render("✓ " + msg)
}

// InfoBox renders an info message in a styled box.
func InfoBox(msg string) string {
	return lipgloss.NewStyle().
		Foreground(ColorAccent).
		Italic(true).
		Render("ℹ " + msg)
}

// CheckBox returns a styled checkbox for the given state.
func CheckBox(checked bool, label string) string {
	var prefix string
	if checked {
		prefix = lipgloss.NewStyle().Foreground(ColorSuccess).Render("✓ ")
	} else {
		prefix = lipgloss.NewStyle().Foreground(ColorGray).Render("○ ")
	}
	return prefix + label
}
