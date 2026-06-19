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
	scroll  int
	layouts []string
	Next    bool
}

const kbViewportHeight = 10

// NewKeyboardModel creates the keyboard layout selection screen.
func NewKeyboardModel(config *model.Config) KeyboardModel {
	return KeyboardModel{
		config:  config,
		cursor:  0,
		scroll:  0,
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
				if m.cursor < m.scroll {
					m.scroll--
				}
			}
		case "down", "j":
			if m.cursor < len(m.layouts)-1 {
				m.cursor++
				if m.cursor >= m.scroll+kbViewportHeight {
					m.scroll++
				}
			}
		case "enter":
			m.config.KeyboardLayout = m.layouts[m.cursor]
			m.Next = true
		}
	}
	return m, nil
}

func (m KeyboardModel) View() string {
	title := TitleStyle.Render("Keyboard Layout")
	subtitle := SubtitleStyle.Render("Select your keyboard layout using ↑/↓, press ENTER to confirm.")

	// Scroll indicator
	total := len(m.layouts)
	end := m.scroll + kbViewportHeight
	if end > total {
		end = total
	}
	scrollInfo := ""
	if total > kbViewportHeight {
		scrollInfo = lipgloss.NewStyle().Foreground(ColorGray).Render(
			"  [" + intToStr(m.scroll+1) + "-" + intToStr(end) + " of " + intToStr(total) + "]  ▼",
		)
	}

	var items string
	visible := m.layouts[m.scroll:end]
	for i, layout := range visible {
		idx := m.scroll + i
		sel := m.config.KeyboardLayout == layout
		items += ListItem(idx == m.cursor, sel, RadioButton(sel, layout)) + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", scrollInfo, "", BoxStyle.Render(items))
}
