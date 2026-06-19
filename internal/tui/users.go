package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UsersModel handles user account configuration.
type UsersModel struct {
	config     *model.Config
	Next       bool
	focusIndex int
	inputs     []textinput.Model
}

// NewUsersModel creates the user configuration screen.
func NewUsersModel(config *model.Config) UsersModel {
	m := UsersModel{
		config: config,
		inputs: make([]textinput.Model, 3),
	}

	m.inputs[0] = newPasswordInput("Root Password", config.RootPassword)
	m.inputs[1] = newTextInput("Username (for sudo user)", config.UserName)
	m.inputs[2] = newPasswordInput("User Password", config.UserPassword)

	return m
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

func (m UsersModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UsersModel) Update(msg tea.Msg) (UsersModel, tea.Cmd) {
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
		case "tab":
			m.saveInputs()
			if m.config.RootPassword == "" {
				return m, nil
			}
			m.Next = true
			return m, nil
		case "enter":
			m.saveInputs()
			if m.focusIndex == len(m.inputs)-1 && m.config.RootPassword != "" {
				m.Next = true
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m *UsersModel) updateFocus() {
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

func (m *UsersModel) saveInputs() {
	m.config.RootPassword = m.inputs[0].Value()
	m.config.UserName = m.inputs[1].Value()
	m.config.UserPassword = m.inputs[2].Value()
	m.config.CreateUser = m.config.UserName != ""
}

func (m UsersModel) View() string {
	title := TitleStyle.Render("Users & Passwords")
	subtitle := SubtitleStyle.Render("Set root password and create a sudo user.")

	rootLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Root Password:")
	userLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Username:")
	passLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("User Password:")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("⚠ Root password must be at least 8 characters"),
		"",
		rootLabel,
		m.inputs[0].View(),
		"",
		DividerStyle.Render(" Optional Sudo User "),
		"",
		userLabel,
		m.inputs[1].View(),
		"",
		passLabel,
		m.inputs[2].View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(content),
	)
}