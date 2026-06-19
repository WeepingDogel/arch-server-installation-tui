package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// ArchLogo returns the Arch Linux ASCII art diamond logo.
func ArchLogo() string {
	logo := `            /\
           /  \
          /    \
         _\     \
        /        \
       /          \
      /     __   \_\
     /     /  \     \
    /__,--'    '--,__\`
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render(logo)
}

// MiniArchLogo returns a small Arch logo for the header.
func MiniArchLogo() string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("  /\\  ")
}

// WelcomeBanner returns the full welcome banner with logo and title.
func WelcomeBanner() string {
	logo := ArchLogo()
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render("Arch Linux Server Installer")
	subtitle := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Italic(true).
		Render("v1.0.0 \u2014 Interactive Server Setup")

	credit := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("by WeepingDogel")

	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(50)

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			logo,
			"",
			title,
			subtitle,
			"",
			credit,
		),
	)
}

// StepIndicator returns a styled step indicator string.
func StepIndicator(current, total int) string {
	var steps string
	for i := 1; i <= total; i++ {
		var s string
		if i == current {
			s = StepActiveStyle.Render("\u25CF")
		} else if i < current {
			s = StepDoneStyle.Render("\u2713")
		} else {
			s = StepStyle.Render("\u25CB")
		}
		if i > 1 {
			steps += " "
		}
		steps += s
	}
	label := lipgloss.NewStyle().Foreground(ColorAccent).Render(
		" Step " + intToStr(current) + "/" + intToStr(total) + " ",
	)
	return lipgloss.JoinHorizontal(lipgloss.Center, steps, label)
}

// intToStr is a simple integer to string converter.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}

// TotalSteps returns the total number of installation steps.
const TotalSteps = 13

// Step names for each step in the wizard.
var StepNames = []string{
	"Welcome",
	"Keyboard Layout",
	"Network",
	"Mirror Selection",
	"Disk Partitioning",
	"Filesystem",
	"Bootloader",
	"Timezone & Locale",
	"Users & Passwords",
	"SSH Configuration",
	"Package Selection",
	"Summary",
	"Installation",
}

// StepName returns the name of the given step (1-based).
func StepName(step int) string {
	if step < 1 || step > len(StepNames) {
		return "Unknown"
	}
	return StepNames[step-1]
}
