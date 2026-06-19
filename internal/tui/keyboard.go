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
	GoBack  bool
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

// itemCount returns layouts + Next + Back
func (m KeyboardModel) totalItems() int {
	return len(m.layouts) + 2 // 2 for Next, Back
}

func (m KeyboardModel) Update(msg tea.Msg) (KeyboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < len(m.layouts) && m.cursor < m.scroll {
					m.scroll--
				}
			}
		case "down", "j":
			if m.cursor < m.totalItems()-1 {
				m.cursor++
				if m.cursor < len(m.layouts) && m.cursor >= m.scroll+kbViewportHeight {
					m.scroll++
				}
			}
		case "enter":
			total := m.totalItems()
			if m.cursor == total-2 { // Next button - use current selection
				if m.config.KeyboardLayout == "" {
					m.config.KeyboardLayout = m.layouts[0]
				}
				m.Next = true
			} else if m.cursor == total-1 { // Back button
				m.GoBack = true
			} else if m.cursor < len(m.layouts) {
				m.config.KeyboardLayout = m.layouts[m.cursor]
			}
		}
	}
	return m, nil
}

func (m KeyboardModel) View() string {
	title := TitleStyle.Render("Keyboard Layout")
	subtitle := SubtitleStyle.Render("↑/↓ to select, ENTER on [Next] to confirm.")

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

	// Render nav buttons
	items += "\n" + renderNavButtons(m.cursor, total, total+1)

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", scrollInfo, "", BoxStyle.Render(items))
}

// renderNavButtons shows Back/Next as cursor-selectable items.
func renderNavButtons(cursor, nextIdx, backIdx int) string {
	var items string
	if cursor == nextIdx {
		items += lipgloss.NewStyle().Background(ColorPrimary).Foreground(ColorWhite).Bold(true).Padding(0, 4).Render(" [ Next ▶ ] ") + "\n"
	} else {
		items += lipgloss.NewStyle().Background(ColorDark).Foreground(ColorPrimary).Padding(0, 4).Render(" [ Next ▶ ] ") + "\n"
	}
	if cursor == backIdx {
		items += lipgloss.NewStyle().Background(ColorAccent).Foreground(ColorWhite).Bold(true).Padding(0, 4).Render(" [ ◀ Back ] ")
	} else {
		items += lipgloss.NewStyle().Background(ColorDark).Foreground(ColorAccent).Padding(0, 4).Render(" [ ◀ Back ] ")
	}
	return items
}
