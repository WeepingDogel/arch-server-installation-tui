package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyboardModel handles keyboard layout selection.
type KeyboardModel struct {
	config  *model.Config
	cursor  int
	layouts []string
	Next    bool
}

// NewKeyboardModel creates the keyboard layout selection screen.
func NewKeyboardModel(config *model.Config) KeyboardModel {
	return KeyboardModel{
		config:  config,
		cursor:  0,
		layouts: []string{"us", "uk", "de", "fr", "jp", "cn", "kr", "br", "it", "es", "ru", "se", "no", "dk", "fi", "pl", "pt", "be", "ca", "sg", "tw", "hu", "cz", "sk", "ro", "bg", "gr", "tr", "is"},
	}
}

func (m KeyboardModel) Init() tea.Cmd { return nil }

func (m KeyboardModel) Update(msg tea.Msg) (KeyboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.layouts)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.config.KeyboardLayout = m.layouts[m.cursor]
			m.Next = true
		case "tab":
			m.config.KeyboardLayout = m.layouts[m.cursor]
			m.Next = true
		}
	}
	return m, nil
}

func (m KeyboardModel) View() string {
	var items string
	for i, layout := range m.layouts {
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		current := ""
		if m.config.KeyboardLayout == layout {
			current = SuccessStyle.Render(" ✓")
		}
		items += style.Render(prefix+layout+current) + "\n"
	}

	title := TitleStyle.Render("Keyboard Layout")
	subtitle := SubtitleStyle.Render("Select your keyboard layout using ↑/↓, press ENTER to confirm.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(items),
	)
}