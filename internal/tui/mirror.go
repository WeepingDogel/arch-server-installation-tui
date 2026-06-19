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
	config       *model.Config
	cursor       int
	scroll       int // scroll offset for viewport
	mirrors      []mirror.Mirror
	filtered     []mirror.Mirror
	showSpecial  bool
	showAll      bool
	enableArchCN bool
	focusBottom  bool // false=mirror list, true=Arch CN toggle
	Next         bool
}

const viewportHeight = 10

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
		focusBottom:  false,
	}
}

func (m MirrorModel) Init() tea.Cmd { return nil }

func (m MirrorModel) Update(msg tea.Msg) (MirrorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.focusBottom {
				m.focusBottom = false
				m.cursor = m.scroll + viewportHeight - 1
				if m.cursor >= len(m.filtered) {
					m.cursor = len(m.filtered) - 1
				}
			} else if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll--
				}
			}
		case "down", "j":
			if !m.focusBottom && m.cursor < len(m.filtered)-1 {
				m.cursor++
				if m.cursor >= m.scroll+viewportHeight {
					m.scroll++
				}
			} else {
				m.focusBottom = true
			}
		case " ":
			if m.focusBottom {
				m.enableArchCN = !m.enableArchCN
				m.config.EnableArchCN = m.enableArchCN
				return m, nil
			}
			// Space on a mirror item selects it
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				selected := m.filtered[m.cursor]
				m.config.MirrorURL = selected.URL
				m.config.MirrorCountry = selected.Country
			}
		case "enter":
			if m.focusBottom {
				// Confirm and proceed
				if m.config.MirrorURL == "" && len(m.filtered) > 0 {
					selected := m.filtered[0]
					m.config.MirrorURL = selected.URL
					m.config.MirrorCountry = selected.Country
				}
				m.config.EnableArchCN = m.enableArchCN
				m.Next = true
			} else if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				selected := m.filtered[m.cursor]
				m.config.MirrorURL = selected.URL
				m.config.MirrorCountry = selected.Country
			}
		case "tab":
			if m.config.MirrorURL == "" && len(m.filtered) > 0 {
				m.config.MirrorURL = m.filtered[0].URL
			}
			m.config.EnableArchCN = m.enableArchCN
			m.Next = true
		case "s":
			m.showSpecial = !m.showSpecial
			m.showAll = !m.showAll
			m.cursor = 0
			m.scroll = 0
			m.focusBottom = false
			if m.showAll {
				m.filtered = m.mirrors
			} else {
				m.filtered = mirror.FilterSpecial(m.mirrors)
			}
		case "a":
			m.showAll = !m.showAll
			m.showSpecial = !m.showAll
			m.cursor = 0
			m.scroll = 0
			m.focusBottom = false
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
		"Select a mirror. 's'=special, 'a'=all, %d mirrors. Use ↑/↓ to scroll.",
		len(m.filtered),
	))

	modeStr := "[Special Mirrors ★]"
	if m.showAll {
		modeStr = "[All Mirrors]"
	}
	mode := lipgloss.NewStyle().Foreground(ColorAccent).Italic(true).Render(modeStr)

	// Scroll indicator
	totalItems := len(m.filtered)
	visibleEnd := m.scroll + viewportHeight
	if visibleEnd > totalItems {
		visibleEnd = totalItems
	}
	scrollInfo := ""
	if totalItems > viewportHeight {
		scrollInfo = lipgloss.NewStyle().Foreground(ColorGray).Render(
			fmt.Sprintf("  [%d-%d of %d]  ▼", m.scroll+1, visibleEnd, totalItems),
		)
	}

	var items string
	for i := m.scroll; i < visibleEnd; i++ {
		mrr := m.filtered[i]
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor && !m.focusBottom {
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
			prefix, mrr.Name, country, star, selected,
		)
		items += style.Render(itemStr) + "\n"
	}

	// Arch Linux CN toggle
	archCNLabel := lipgloss.NewStyle().Foreground(ColorWhite).Render("Arch Linux CN repository: ")
	var archCNBtn string
	if m.enableArchCN {
		archCNBtn = lipgloss.NewStyle().
			Background(ColorSuccess).
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 2).
			Render("  Enabled  ")
	} else {
		archCNBtn = lipgloss.NewStyle().
			Background(ColorGray).
			Foreground(ColorWhite).
			Padding(0, 2).
			Render("  Disabled ")
	}

	archCNToggle := lipgloss.JoinHorizontal(lipgloss.Center, archCNLabel, archCNBtn)
	archCNHint := HelpStyle.Render("SPACE to toggle Arch Linux CN (adds AUR/Chinese packages)")

	// Bottom confirm
	bottomStyle := ListItemStyle
	bottomPrefix := "  "
	if m.focusBottom {
		bottomStyle = ListItemSelectedStyle
		bottomPrefix = "▶ "
	}

	bottomContent := lipgloss.JoinVertical(
		lipgloss.Left,
		bottomStyle.Render(bottomPrefix+archCNToggle),
		"",
		archCNHint,
	)

	// Confirm hint
	confirmHint := InfoBox("ENTER to confirm selection and continue.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		mode,
		scrollInfo,
		"",
		BoxStyle.Render(items),
		"",
		BoxStyle.Render(bottomContent),
		"",
		confirmHint,
	)
}
