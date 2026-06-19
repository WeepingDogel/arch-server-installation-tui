package tui

import (
	"strconv"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SSHModel handles SSH server configuration.
type SSHModel struct {
	config     *model.Config
	Next       bool
	focusIndex int
	inputs     []textinput.Model
}

// NewSSHModel creates the SSH configuration screen.
func NewSSHModel(config *model.Config) SSHModel {
	m := SSHModel{
		config: config,
		inputs: make([]textinput.Model, 1),
	}

	m.inputs[0] = newTextInput("SSH Port (default: 22)", strconv.Itoa(config.SSHPort))
	m.inputs[0].SetValue("22")

	return m
}

func (m SSHModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SSHModel) Update(msg tea.Msg) (SSHModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.focusIndex > 0 {
				m.focusIndex--
				m.updateFocus()
			}
			return m, nil
		case "down", "j":
			if m.focusIndex < len(m.inputs)-1 {
				m.focusIndex++
				m.updateFocus()
			}
			return m, nil
		case " ", "enter":
			// Toggle SSH enable/disable
			if m.focusIndex == 0 && m.config.EnableSSH {
				m.config.EnableSSH = !m.config.EnableSSH
				return m, nil
			}
			m.saveInputs()
			m.Next = true
			return m, nil
		case "tab":
			m.saveInputs()
			m.Next = true
			return m, nil
		case "r":
			m.config.AllowRootLogin = !m.config.AllowRootLogin
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.focusIndex >= 0 && m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	}
	return m, cmd
}

func (m *SSHModel) updateFocus() {
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].TextStyle = InputFocusStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].TextStyle = InputStyle
		}
	}
}

func (m *SSHModel) saveInputs() {
	portStr := m.inputs[0].Value()
	if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
		m.config.SSHPort = port
	}
}

func (m SSHModel) View() string {
	title := TitleStyle.Render("SSH Configuration")
	subtitle := SubtitleStyle.Render("Configure OpenSSH server settings.")

	// SSH enable/disable toggle
	enableLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Enable SSH: ")
	var enableBtn string
	if m.config.EnableSSH {
		enableBtn = lipgloss.NewStyle().
			Background(ColorSuccess).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Render("  Enabled  ")
	} else {
		enableBtn = lipgloss.NewStyle().
			Background(ColorError).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Render(" Disabled ")
	}

	// Root login toggle
	rootLoginLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Allow Root Login: ")
	var rootBtn string
	if m.config.AllowRootLogin {
		rootBtn = lipgloss.NewStyle().
			Background(ColorWarning).
			Foreground(ColorDark).
			Bold(true).
			Padding(0, 2).
			Render("   Yes   ")
	} else {
		rootBtn = lipgloss.NewStyle().
			Background(ColorGray).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Render("   No    ")
	}

	portLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("SSH Port:")
	portInput := m.inputs[0].View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, enableLabel, enableBtn),
		HelpStyle.Render("Press SPACE to toggle SSH"),
		"",
		portLabel,
		portInput,
		"",
		lipgloss.JoinHorizontal(lipgloss.Center, rootLoginLabel, rootBtn),
		HelpStyle.Render("Press 'r' to toggle root login"),
		"",
		InfoBox("Disabling root login is recommended for security."),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(content),
	)
}