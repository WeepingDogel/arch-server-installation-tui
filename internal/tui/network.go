package tui

import (
	"os"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NetworkModel handles network configuration with NIC detection.
type NetworkModel struct {
	config      *model.Config
	Next        bool
	focusIndex  int
	dhcpToggle  bool
	inputs      []textinput.Model
	showStatic  bool
	interfaces  []string // detected network interfaces
	ifaceCursor int
}

// detectNICs reads /sys/class/net for real network interfaces.
func detectNICs() []string {
	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		// Fallback
		return []string{"eth0", "enp0s3", "wlan0"}
	}
	var ifaces []string
	for _, e := range entries {
		name := e.Name()
		if name != "lo" {
			ifaces = append(ifaces, name)
		}
	}
	if len(ifaces) == 0 {
		return []string{"eth0", "enp0s3", "wlan0"}
	}
	return ifaces
}

// NewNetworkModel creates the network configuration screen with detected NICs.
func NewNetworkModel(config *model.Config) NetworkModel {
	ifaces := detectNICs()
	m := NetworkModel{
		config:      config,
		dhcpToggle:  true,
		interfaces:  ifaces,
		ifaceCursor: 0,
		inputs:      make([]textinput.Model, 5),
	}

	m.inputs[0] = newTextInput("Hostname", config.Hostname)
	m.inputs[1] = newTextInput("IP Address", config.IPAddress)
	m.inputs[2] = newTextInput("Netmask", config.Netmask)
	m.inputs[3] = newTextInput("Gateway", config.Gateway)
	m.inputs[4] = newTextInput("DNS Servers", config.DNSServers)

	// Set detected interface
	if len(ifaces) > 0 {
		config.NetworkIface = ifaces[0]
	}

	return m
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
			}
			return m, nil
		case "left":
			if m.focusIndex == 0 && !m.showStatic && m.ifaceCursor > 0 {
				m.ifaceCursor--
				m.config.NetworkIface = m.interfaces[m.ifaceCursor]
			}
			return m, nil
		case "right":
			if m.focusIndex == 0 && !m.showStatic && m.ifaceCursor < len(m.interfaces)-1 {
				m.ifaceCursor++
				m.config.NetworkIface = m.interfaces[m.ifaceCursor]
			}
			return m, nil
		case " ":
			if m.focusIndex == 0 && !m.showStatic {
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
			if m.showStatic && m.focusIndex == len(m.inputs)-1 {
				m.saveInputs()
				m.Next = true
			}
			return m, nil
		case "enter":
			m.saveInputs()
			m.Next = true
			return m, nil
		case "tab":
			m.saveInputs()
			m.Next = true
			return m, nil
		}
	}

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

	// NIC selector
	ifaceLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Interface: ")
	var ifaceBtn string
	if len(m.interfaces) > 0 {
		ifaceBtn = lipgloss.NewStyle().
			Background(ColorAccent).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Render(" " + m.interfaces[m.ifaceCursor] + " ")
	}
	ifaceSelector := lipgloss.JoinHorizontal(lipgloss.Center, ifaceLabel, ifaceBtn)
	ifaceHint := HelpStyle.Render("←/→ to change interface")

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
		ifaceSelector,
		ifaceHint,
		"",
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
