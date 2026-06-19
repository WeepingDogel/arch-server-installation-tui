package tui

import (
	"fmt"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/mirror"
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MirrorModel handles mirror selection.
type MirrorModel struct {
	config      *model.Config
	cursor      int
	mirrors     []mirror.Mirror
	filtered    []mirror.Mirror
	showSpecial bool // Toggle: show all or only special mirrors
	showAll     bool // Show all countries
	Next        bool
}

// NewMirrorModel creates the mirror selection screen.
func NewMirrorModel(config *model.Config) MirrorModel {
	all := mirror.DefaultMirrors()
	return MirrorModel{
		config:      config,
		cursor:      0,
		mirrors:     all,
		filtered:    mirror.FilterSpecial(all),
		showSpecial: true,
		showAll:     false,
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
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				selected := m.filtered[m.cursor]
				m.config.MirrorURL = selected.URL
				m.config.MirrorCountry = selected.Country
				m.Next = true
			}
		case "tab":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				selected := m.filtered[m.cursor]
				m.config.MirrorURL = selected.URL
				m.config.MirrorCountry = selected.Country
				m.Next = true
			}
		case "s":
			// Toggle special/all mirrors
			m.showSpecial = !m.showSpecial
			m.showAll = !m.showAll
			m.cursor = 0
			if m.showAll {
				m.filtered = m.mirrors
			} else {
				m.filtered = mirror.FilterSpecial(m.mirrors)
			}
		case "a":
			// Show all mirrors
			m.showAll = !m.showAll
			m.showSpecial = !m.showAll
			m.cursor = 0
			if m.showAll {
				m.filtered = m.mirrors
			} else {
				m.filtered = mirror.FilterSpecial(m.mirrors)
			}
		}
	}
	return m, nil
}

func (m MirrorModel) View() string {
	title := TitleStyle.Render("Mirror Selection")
	subtitle := SubtitleStyle.Render(fmt.Sprintf(
		"Choose a mirror server. Press 's' for special mirrors, 'a' for all. Found: %d mirrors",
		len(m.filtered),
	))

	// Mode indicator
	modeStr := "[Special Mirrors ★]"
	if m.showAll {
		modeStr = "[All Mirrors]"
	}
	mode := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Italic(true).
		Render(modeStr)

	var items string
	for i, mrr := range m.filtered {
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if m.config.MirrorURL == mrr.URL {
			selected = SuccessStyle.Render(" ✓")
		}
		country := lipgloss.NewStyle().Foreground(ColorGray).Render(" [" + mrr.Country + "]")
		star := ""
		if mrr.Special {
			star = lipgloss.NewStyle().Foreground(ColorWarning).Render(" ★")
		}

		itemStr := fmt.Sprintf("%s%s%s%s%s",
			prefix,
			mrr.Name,
			country,
			star,
			selected,
		)
		items += style.Render(itemStr) + "\n"
	}

	// Custom mirror input hint
	customHint := InfoBox("Use ESC to go back, ENTER to select. Press 's' or 'a' to toggle mirror groups.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		mode,
		"",
		BoxStyle.Render(items),
		"",
		customHint,
	)
}