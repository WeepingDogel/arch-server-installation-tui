package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
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

// NavBar renders styled Back / Next navigation buttons for inline use in screens.
func NavBar(step, total, focus int) string {
	var backBtn, nextBtn string
	if step > 1 {
		if focus == 0 {
			backBtn = lipgloss.NewStyle().Background(ColorAccent).Foreground(ColorWhite).Bold(true).Padding(0, 4).Render(" ◀ Back ")
		} else {
			backBtn = lipgloss.NewStyle().Background(ColorDark).Foreground(ColorAccent).Padding(0, 4).Render(" ◀ Back ")
		}
	}
	if step < total {
		if focus == 1 {
			nextBtn = lipgloss.NewStyle().Background(ColorPrimary).Foreground(ColorWhite).Bold(true).Padding(0, 4).Render(" Next ▶ ")
		} else {
			nextBtn = lipgloss.NewStyle().Background(ColorDark).Foreground(ColorPrimary).Padding(0, 4).Render(" Next ▶ ")
		}
	}
	if step == total-1 {
		if focus == 1 {
			nextBtn = lipgloss.NewStyle().Background(ColorSuccess).Foreground(ColorWhite).Bold(true).Padding(0, 4).Render(" Install ▶ ")
		} else {
			nextBtn = lipgloss.NewStyle().Background(ColorDark).Foreground(ColorSuccess).Padding(0, 4).Render(" Install ▶ ")
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, backBtn, "  ", nextBtn)
}

// SimpleFooter returns footer hints.
func SimpleFooter() string {
	return HelpStyle.Render("↑/↓  •  Enter  •  Esc Back  •  Ctrl+C Quit")
}

func renderHeader(step int) string {
	logo := MiniArchLogo()
	stepInd := StepIndicator(step, TotalSteps)
	stepName := lipgloss.NewStyle().Foreground(ColorWhite).Bold(true).Render(" " + StepName(step) + " ")
	return HeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, logo, stepInd, stepName))
}

func InfoBox(msg string) string {
	return lipgloss.NewStyle().Foreground(ColorAccent).Italic(true).Render("ℹ " + msg)
}

func ErrorBox(err string) string {
	if err == "" {
		return ""
	}
	return ErrorStyle.Render("✗ " + err)
}

func SuccessBox(msg string) string {
	return SuccessStyle.Render("✓ " + msg)
}

func newTextInput(placeholder, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.Width = 50
	ti.TextStyle = InputStyle
	return ti
}

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

func RadioButton(selected bool, label string) string {
	if selected {
		return lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("◉") + " " + label
	}
	return lipgloss.NewStyle().Foreground(ColorGray).Render("○") + " " + label
}

func Checkbox(checked bool, label string) string {
	if checked {
		return lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("☑") + " " + label
	}
	return lipgloss.NewStyle().Foreground(ColorGray).Render("☐") + " " + label
}

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
