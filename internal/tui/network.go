package tui

import (
	"os"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NetworkModel handles network configuration.
type NetworkModel struct {
	config      *model.Config
	Next        bool
	GoBack      bool
	cursor      int // unified cursor: 0=inputs, last 2 = [Next] [Back]
	dhcpToggle  bool
	inputs      []textinput.Model
	showStatic  bool
	interfaces  []string
	ifaceCursor int
}

func (m NetworkModel) totalItems() int {
	if m.showStatic {
		return len(m.inputs) + 2 // inputs + [Next] [Back]
	}
	return 1 + 2 // hostname + [Next] [Back]
}

func detectNICs() []string {
	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return []string{"eth0", "enp0s3", "wlan0"}
	}
	var ifaces []string
	for _, e := range entries {
		if e.Name() != "lo" {
			ifaces = append(ifaces, e.Name())
		}
	}
	if len(ifaces) == 0 {
		return []string{"eth0", "enp0s3", "wlan0"}
	}
	return ifaces
}

// NewNetworkModel creates the network configuration screen.
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
	if len(ifaces) > 0 {
		config.NetworkIface = ifaces[0]
	}
	return m
}

func (m NetworkModel) Init() tea.Cmd { return textinput.Blink }

func (m NetworkModel) Update(msg tea.Msg) (NetworkModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.syncInputFocus()
			}
			return m, nil
		case "down", "j":
			if m.cursor < m.totalItems()-1 {
				m.cursor++
				m.syncInputFocus()
			}
			return m, nil
		case "left":
			if m.cursor == 0 && m.ifaceCursor > 0 {
				m.ifaceCursor--
				m.config.NetworkIface = m.interfaces[m.ifaceCursor]
			}
			return m, nil
		case "right":
			if m.cursor == 0 && m.ifaceCursor < len(m.interfaces)-1 {
				m.ifaceCursor++
				m.config.NetworkIface = m.interfaces[m.ifaceCursor]
			}
			return m, nil
		case " ":
			m.dhcpToggle = !m.dhcpToggle
			m.config.NetworkDHCP = m.dhcpToggle
			m.showStatic = !m.dhcpToggle
			m.syncInputFocus()
			return m, nil
		case "enter":
			total := m.totalItems()
			if m.cursor == total-2 {
				m.saveInputs()
				m.Next = true
			} else if m.cursor == total-1 {
				m.GoBack = true
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	// Update focused input
	focus := m.cursor
	if focus >= 0 && focus < len(m.inputs) && (m.showStatic || focus == 0) {
		m.inputs[focus], cmd = m.inputs[focus].Update(msg)
	}
	return m, cmd
}

func (m *NetworkModel) syncInputFocus() {
	for i := range m.inputs {
		shouldFocus := i == m.cursor && (m.showStatic || i == 0)
		if shouldFocus {
			m.inputs[i].Focus()
			m.inputs[i].TextStyle = InputFocusStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].TextStyle = InputStyle
		}
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
	subtitle := SubtitleStyle.Render("Configure network, ↓ to reach [Next].")

	nicStr := " [ " + m.interfaces[m.ifaceCursor] + " ] "
	ifaceLine := lipgloss.NewStyle().Foreground(ColorWhite).Render("Interface:") + " " +
		lipgloss.NewStyle().Background(ColorAccent).Foreground(ColorWhite).Bold(true).Padding(0, 2).Render(nicStr) +
		" " + HelpStyle.Render("←/→ change")

	modeStr := " DHCP "
	if !m.dhcpToggle {
		modeStr = " Static "
	}
	modeColor := ColorSuccess
	if !m.dhcpToggle {
		modeColor = ColorWarning
	}
	toggleLine := lipgloss.NewStyle().Foreground(ColorWhite).Render("Mode:") + " " +
		lipgloss.NewStyle().Background(modeColor).Foreground(ColorWhite).Bold(true).Padding(0, 2).Render(modeStr) +
		" " + HelpStyle.Render("SPACE toggle")

	items := ifaceLine + "\n" + toggleLine + "\n\n"
	items += lipgloss.NewStyle().Foreground(ColorWhite).Render("Hostname:") + "\n" + m.inputs[0].View() + "\n"

	if m.showStatic {
		labels := []string{"IP Address:", "Netmask:", "Gateway:", "DNS Servers:"}
		for i, label := range labels {
			idx := i + 1
			items += "\n" + lipgloss.NewStyle().Foreground(ColorWhite).Render(label) + "\n" + m.inputs[idx].View() + "\n"
		}
	}

	items += "\n" + renderNavButtons(m.cursor, m.totalItems()-2, m.totalItems()-1)

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", BoxStyle.Render(items))
}
