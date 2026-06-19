package tui

import (
	"fmt"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/mirror"
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const mirrorViewHeight = 8

// MirrorModel handles mirror selection.
type MirrorModel struct {
	config       *model.Config
	cursor       int
	scroll       int
	mirrors      []mirror.Mirror
	filtered     []mirror.Mirror
	showSpecial  bool
	showAll      bool
	enableArchCN bool
	Next         bool
	GoBack       bool
}

func (m MirrorModel) totalItems() int {
	return len(m.filtered) + 2 + 1 // items + [Next] [Back] + Arch CN toggle
}

// NewMirrorModel creates the mirror selection screen.
func NewMirrorModel(config *model.Config) MirrorModel {
	all := mirror.DefaultMirrors()
	return MirrorModel{
		config:       config,
		cursor:       0,
		scroll:       0,
		mirrors:      all,
		filtered:     mirror.FilterSpecial(all),
		showSpecial:  true,
		showAll:      false,
		enableArchCN: config.EnableArchCN,
	}
}

func (m MirrorModel) Init() tea.Cmd { return nil }

func (m MirrorModel) Update(msg tea.Msg) (MirrorModel, tea.Cmd) {
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
			if m.cursor < m.totalItems()-1 {
				m.cursor++
				if m.cursor < len(m.filtered) && m.cursor >= m.scroll+mirrorViewHeight {
					m.scroll++
				}
			}
		case " ":
			itemCount := len(m.filtered)
			if m.cursor == itemCount {
				m.enableArchCN = !m.enableArchCN
				m.config.EnableArchCN = m.enableArchCN
			} else if m.cursor < itemCount {
				selected := m.filtered[m.cursor]
				m.config.MirrorURL = selected.URL
				m.config.MirrorCountry = selected.Country
			}
		case "enter":
			total := m.totalItems()
			if m.cursor == total-2 {
				m.config.EnableArchCN = m.enableArchCN
				if m.config.MirrorURL == "" && len(m.filtered) > 0 {
					m.config.MirrorURL = m.filtered[0].URL
				}
				m.Next = true
			} else if m.cursor == total-1 {
				m.GoBack = true
			} else if m.cursor < len(m.filtered) {
				selected := m.filtered[m.cursor]
				m.config.MirrorURL = selected.URL
				m.config.MirrorCountry = selected.Country
			}
		case "s":
			m.showSpecial = true
			m.showAll = false
			m.filtered = mirror.FilterSpecial(m.mirrors)
			m.cursor = 0
			m.scroll = 0
		case "a":
			m.showAll = true
			m.showSpecial = false
			m.filtered = m.mirrors
			m.cursor = 0
			m.scroll = 0
		}
	}
	return m, nil
}

func (m MirrorModel) View() string {
	title := TitleStyle.Render("Mirror Selection")
	subtitle := SubtitleStyle.Render(fmt.Sprintf("'s'=special 'a'=all • %d mirrors • ↓ to [Next]", len(m.filtered)))

	modeStr := " [Special Mirrors ★] "
	if m.showAll {
		modeStr = " [All Mirrors] "
	}
	mode := lipgloss.NewStyle().Foreground(ColorAccent).Italic(true).Render(modeStr)

	total := len(m.filtered)
	visibleEnd := m.scroll + mirrorViewHeight
	if visibleEnd > total {
		visibleEnd = total
	}
	scrollInfo := ""
	if total > mirrorViewHeight {
		scrollInfo = lipgloss.NewStyle().Foreground(ColorGray).Render(fmt.Sprintf("  [%d-%d of %d]", m.scroll+1, visibleEnd, total))
	}

	var items string
	for i := m.scroll; i < visibleEnd; i++ {
		mrr := m.filtered[i]
		sel := m.config.MirrorURL == mrr.URL
		label := mrr.Name + " [" + mrr.Country + "]"
		if mrr.Special {
			label += " ★"
		}
		items += ListItem(i == m.cursor, sel, RadioButton(sel, label)) + "\n"
	}

	// Arch CN toggle
	cnLine := ""
	if m.cursor == total {
		cnLine = ListItem(true, m.enableArchCN, Checkbox(m.enableArchCN, "Arch Linux CN repository"))
	} else {
		cnLine = ListItem(false, m.enableArchCN, Checkbox(m.enableArchCN, "Arch Linux CN repository"))
	}
	items += "\n" + cnLine + "\n"
	items += renderNavButtons(m.cursor, total+1, total+2)

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", mode, scrollInfo, "", BoxStyle.Render(items))
}
