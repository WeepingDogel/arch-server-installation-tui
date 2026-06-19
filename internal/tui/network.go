package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NetworkModel handles network configuration.
type NetworkModel struct {
	config     *model.Config
	Next       bool
	focusIndex int
	dhcpToggle bool // true = DHCP, false = Static
	inputs     []textinput.Model
	showStatic bool
}

// NewNetworkModel creates the network configuration screen.
func NewNetworkModel(config *model.Config) NetworkModel {
	m := NetworkModel{
		config:     config,
		dhcpToggle: true,
		inputs:     make([]textinput.Model, 5),
	}

	m.inputs[0] = newTextInput("Hostname", config.Hostname)
	m.inputs[1] = newTextInput("IP Address", config.IPAddress)
	m.inputs[2] = newTextInput("Netmask", config.Netmask)
	m.inputs[3] = newTextInput("Gateway", config.Gateway)
	m.inputs[4] = newTextInput("DNS Servers", config.DNSServers)

	return m
}

func newTextInput(placeholder, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.Width = 50
	ti.TextStyle = InputStyle
	return ti
}

func (m NetworkModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m NetworkModel) Update(msg tea.Msg) (NetworkModel, tea.Cmd) {
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
			if m.showStatic && m.focusIndex < len(m.inputs)-1 {
				m.focusIndex++
				m.updateFocus()
			} else if !m.showStatic && m.focusIndex < 0 {
				m.focusIndex = 0
				m.updateFocus()
			}
			return m, nil

		case " ", "enter":
			if m.focusIndex == 0 && !m.showStatic {
				// Toggle DHCP when focused on the toggle area
				m.dhcpToggle = !m.dhcpToggle
				m.config.NetworkDHCP = m.dhcpToggle
				m.showStatic = !m.dhcpToggle
				if m.dhcpToggle {
					m.focusIndex = 0
				} else {
					m.focusIndex = 1
				}
				m.updateFocus()
				return m, nil
			}
			if m.focusIndex == 0 && m.showStatic {
				m.dhcpToggle = !m.dhcpToggle
				m.config.NetworkDHCP = m.dhcpToggle
				m.showStatic = !m.dhcpToggle
				if m.dhcpToggle {
					m.focusIndex = 0
					m.updateBlur()
				}
				m.updateFocus()
				return m, nil
			}
			// Check if we're in the last input
			if m.showStatic && m.focusIndex == len(m.inputs)-1 {
				m.saveInputs()
				m.Next = true
			}
			return m, nil

		case "tab":
			m.saveInputs()
			m.Next = true
			return m, nil
		}
	}

	// Handle input updates
	var cmd tea.Cmd
	if m.showStatic {
		idx := m.focusIndex
		if idx >= 0 && idx < len(m.inputs) {
			m.inputs[idx], cmd = m.inputs[idx].Update(msg)
		}
	}
	return m, cmd
}

func (m *NetworkModel) updateFocus() {
	for i := range m.inputs {
		if m.showStatic && i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].TextStyle = InputFocusStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].TextStyle = InputStyle
		}
	}
}

func (m *NetworkModel) updateBlur() {
	for i := range m.inputs {
		m.inputs[i].Blur()
		m.inputs[i].TextStyle = InputStyle
	}
}

func (m *NetworkModel) saveInputs() {
	m.config.Hostname = m.inputs[0].Value()
	if m.showStatic {
		m.config.IPAddress = m.inputs[1].Value()
		m.config.Netmask = m.inputs[2].Value()
		m.config.Gateway = m.inputs[3].Value()
		m.config.DNSServers = m.inputs[4].Value()
	}
	m.config.NetworkDHCP = m.dhcpToggle
}

func (m NetworkModel) View() string {
	title := TitleStyle.Render("Network Configuration")
	subtitle := SubtitleStyle.Render("Configure network settings for your server.")

	// DHCP/Static toggle
	toggleLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Network Mode: ")
	var toggleBtn string
	if m.dhcpToggle {
		toggleBtn = lipgloss.NewStyle().
			Background(ColorSuccess).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Render("  DHCP  ")
	} else {
		toggleBtn = lipgloss.NewStyle().
			Background(ColorWarning).
			Foreground(ColorDark).
			Bold(true).
			Padding(0, 2).
			Render(" Static ")
	}

	toggle := lipgloss.JoinHorizontal(lipgloss.Center, toggleLabel, toggleBtn)
	toggleHint := HelpStyle.Render("Press SPACE to toggle network mode")

	// Hostname
	hostnameLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Hostname:")
	hostnameInput := m.inputs[0].View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		toggle,
		toggleHint,
		"",
		hostnameLabel,
		hostnameInput,
	)

	// Static fields
	if m.showStatic {
		labels := []string{"IP Address:", "Netmask:", "Gateway:", "DNS Servers:"}
		for i, label := range labels {
			idx := i + 1
			if idx < len(m.inputs) {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					content,
					"",
					lipgloss.NewStyle().Foreground(ColorWhite).Render(label),
					m.inputs[idx].View(),
				)
			}
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(content),
	)
}
