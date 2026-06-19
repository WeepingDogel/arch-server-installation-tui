package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Screen returns a rendered screen with header, content, navigation bar, and footer.
func Screen(step int, content, footer string, navFocus int) string {
	header := renderHeader(step)
	contentBox := lipgloss.NewStyle().
		Padding(0, 1).
		Render(content)

	navBar := NavBar(step, TotalSteps, navFocus, StepName(step))

	footerBox := FooterStyle.Render(footer)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		"",
		contentBox,
		"",
		navBar,
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
		HelpStyle.Render("←/→ Nav buttons  "),
		HelpStyle.Render("Ctrl+C Quit"),
	)
}

// NavBar renders styled Back / Next navigation buttons.
func NavBar(step, total, focus int, stepName string) string {
	var backBtn string
	if step > 1 {
		if focus == 0 {
			backBtn = lipgloss.NewStyle().
				Background(ColorAccent).
				Foreground(ColorWhite).
				Bold(true).
				Padding(0, 4).
				Render("◀  Back")
		} else {
			backBtn = lipgloss.NewStyle().
				Background(ColorDark).
				Foreground(ColorAccent).
				Padding(0, 4).
				Render("◀  Back")
		}
	}

	var nextBtn string
	if step < total {
		if focus == 1 {
			nextBtn = lipgloss.NewStyle().
				Background(ColorPrimary).
				Foreground(ColorWhite).
				Bold(true).
				Padding(0, 4).
				Render("Next  ▶")
		} else {
			nextBtn = lipgloss.NewStyle().
				Background(ColorDark).
				Foreground(ColorPrimary).
				Padding(0, 4).
				Render("Next  ▶")
		}
	}

	// On step 12 (last before install), show Install button
	if step == total-1 {
		if focus == 1 {
			nextBtn = lipgloss.NewStyle().
				Background(ColorSuccess).
				Foreground(ColorWhite).
				Bold(true).
				Padding(0, 4).
				Render("Install  ▶")
		} else {
			nextBtn = lipgloss.NewStyle().
				Background(ColorDark).
				Foreground(ColorSuccess).
				Padding(0, 4).
				Render("Install  ▶")
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		backBtn,
		"  ",
		nextBtn,
	)
}

// SimpleFooter returns a minimal footer with keyboard hints.
func SimpleFooter() string {
	return HelpStyle.Render("↑/↓  •  Enter  •  ←/→ Nav  •  Esc Back  •  Ctrl+C Quit")
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

// newTextInput creates a text input with consistent styling.
func newTextInput(placeholder, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.Width = 50
	ti.TextStyle = InputStyle
	return ti
}

// newPasswordInput creates a password text input with consistent styling.
func newPasswordInput(placeholder, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.Width = 50
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.TextStyle = InputStyle
	return ti
}

// RadioButton renders a radio button indicator.
func RadioButton(selected bool, label string) string {
	var indicator string
	if selected {
		indicator = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("◉ ")
	} else {
		indicator = lipgloss.NewStyle().Foreground(ColorGray).Render("○ ")
	}
	return indicator + label
}

// Checkbox renders a checkbox indicator.
func Checkbox(checked bool, label string) string {
	var indicator string
	if checked {
		indicator = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("☑ ")
	} else {
		indicator = lipgloss.NewStyle().Foreground(ColorGray).Render("☐ ")
	}
	return indicator + label
}

// ListItem renders a single list item with cursor indicator.
func ListItem(isCursor, isSelected bool, label string) string {
	style := ListItemStyle
	prefix := "  "
	if isCursor {
		style = ListItemSelectedStyle
		prefix = "▶ "
	}
	suffix := ""
	if isSelected {
		suffix = " " + lipgloss.NewStyle().Foreground(ColorSuccess).Render("✓")
	}
	return style.Render(prefix + label + suffix)
}
