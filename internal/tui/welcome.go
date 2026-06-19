package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WelcomeModel is the first screen of the installer.
type WelcomeModel struct {
	config *model.Config
	Next   bool
}

// NewWelcomeModel creates the welcome screen.
func NewWelcomeModel(config *model.Config) WelcomeModel {
	return WelcomeModel{config: config}
}

func (m WelcomeModel) Init() tea.Cmd { return nil }

func (m WelcomeModel) Update(msg tea.Msg) (WelcomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ", "tab":
			m.Next = true
		}
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	banner := WelcomeBanner()

	info := InfoBox("This tool will guide you through installing Arch Linux as a production-ready server.\n" +
		"Configure keyboard, network, disks, packages, and more step by step.")

	requirements := BoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			TitleStyle.Render("Requirements:"),
			"",
			CheckBox(true, "Active internet connection"),
			CheckBox(true, "Arch Linux ISO booted"),
			CheckBox(true, "64-bit (x86_64) or ARM64 processor"),
			CheckBox(true, "At least 8GB of disk space"),
			"",
			SubtitleStyle.Render("Press ENTER or TAB to begin ▶"),
		),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		banner,
		"",
		info,
		"",
		requirements,
	)
}
