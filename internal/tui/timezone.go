package tui

import (
	"fmt"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const tzViewportHeight = 8

// TimezoneModel handles timezone and locale configuration.
type TimezoneModel struct {
	config  *model.Config
	cursor  int
	scroll  int
	regions []string
	locales []string
	Next    bool
	GoBack  bool
}

// totalContentItems returns the count of cursorable content items (regions + locales + nav).
func (m TimezoneModel) totalContentItems() int {
	return len(m.regions) + len(m.locales) + 2
}

func (m TimezoneModel) viewHeight() int { return tzViewportHeight }

// NewTimezoneModel creates the timezone configuration screen.
func NewTimezoneModel(config *model.Config) TimezoneModel {
	return TimezoneModel{
		config:  config,
		cursor:  0,
		scroll:  0,
		regions: []string{"UTC", "Asia/Shanghai", "Asia/Tokyo", "Asia/Seoul", "Asia/Singapore", "Asia/Hong_Kong", "Asia/Taipei", "Asia/Kolkata", "Asia/Dubai", "America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles", "America/Sao_Paulo", "Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Moscow", "Europe/Amsterdam", "Australia/Sydney", "Australia/Melbourne", "Pacific/Auckland", "Africa/Cairo", "Africa/Johannesburg"},
		locales: []string{"en_US.UTF-8", "en_GB.UTF-8", "zh_CN.UTF-8", "zh_TW.UTF-8", "ja_JP.UTF-8", "ko_KR.UTF-8", "de_DE.UTF-8", "fr_FR.UTF-8", "es_ES.UTF-8", "pt_BR.UTF-8", "ru_RU.UTF-8", "it_IT.UTF-8", "pl_PL.UTF-8", "sv_SE.UTF-8", "nb_NO.UTF-8", "da_DK.UTF-8", "fi_FI.UTF-8", "nl_NL.UTF-8", "cs_CZ.UTF-8", "hu_HU.UTF-8", "ro_RO.UTF-8", "bg_BG.UTF-8", "el_GR.UTF-8", "tr_TR.UTF-8"},
	}
}

func localeSelected(locales []string, locale string) bool {
	for _, l := range locales {
		if l == locale {
			return true
		}
	}
	return false
}

func toggleLocale(locales []string, locale string) []string {
	for i, l := range locales {
		if l == locale {
			return append(locales[:i], locales[i+1:]...)
		}
	}
	return append(locales, locale)
}

func (m TimezoneModel) Init() tea.Cmd { return nil }

func (m TimezoneModel) Update(msg tea.Msg) (TimezoneModel, tea.Cmd) {
	total := m.totalContentItems()
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
			if m.cursor < total-1 {
				m.cursor++
				if m.cursor >= m.scroll+m.viewHeight() {
					m.scroll++
				}
			}
		case " ":
			contentLen := len(m.regions) + len(m.locales)
			if m.cursor >= len(m.regions) && m.cursor < contentLen {
				localeIdx := m.cursor - len(m.regions)
				if localeIdx < len(m.locales) {
					m.config.Locales = toggleLocale(m.config.Locales, m.locales[localeIdx])
				}
			}
		case "enter":
			if m.cursor == total-2 {
				if m.config.TimezoneRegion == "" && len(m.regions) > 0 {
					m.config.TimezoneRegion = m.regions[0]
				}
				m.Next = true
			} else if m.cursor == total-1 {
				m.GoBack = true
			} else if m.cursor < len(m.regions) {
				m.config.TimezoneRegion = m.regions[m.cursor]
			}
		}
	}
	return m, nil
}

func (m TimezoneModel) View() string {
	title := TitleStyle.Render("Timezone & Locale")
	subtitle := SubtitleStyle.Render("SPACE to toggle locales, ENTER on [Next] to confirm.")

	contentLen := len(m.regions) + len(m.locales)

	// Build a flat list of content items only (no dividers in the scrollable area)
	var items string
	visibleEnd := m.scroll + m.viewHeight()
	if visibleEnd > contentLen {
		visibleEnd = contentLen
	}

	// Track which content index each visible line corresponds to
	for i := m.scroll; i < visibleEnd; i++ {
		var line string
		var isRegion = i < len(m.regions)

		if isRegion {
			region := m.regions[i]
			sel := m.config.TimezoneRegion == region
			line = ListItem(i == m.cursor, sel, RadioButton(sel, region))
		} else {
			localeIdx := i - len(m.regions)
			locale := m.locales[localeIdx]
			checked := localeSelected(m.config.Locales, locale)
			prefix := "☐ "
			if checked {
				prefix = "☑ "
			}
			line = ListItem(i == m.cursor, false, prefix+locale)
		}
		items += line + "\n"
	}

	// Nav buttons (always visible at bottom, not scrolled)
	items += "\n" + renderNavButtons(m.cursor, contentLen, contentLen+1)

	scrollInfo := ""
	if contentLen > m.viewHeight() {
		scrollInfo = lipgloss.NewStyle().Foreground(ColorGray).Render(
			fmt.Sprintf("  [%d-%d of %d]  ▼", m.scroll+1, visibleEnd, contentLen),
		)
	}

	// Show divider labels inline above the content
	regionDivider := DividerStyle.Render(" Timezone ")
	localeDivider := DividerStyle.Render(" Locales (SPACE to toggle) ")

	summary := ""
	if len(m.config.Locales) > 0 {
		locStr := "Selected: "
		for i, l := range m.config.Locales {
			if i > 0 {
				locStr += ", "
			}
			locStr += l
		}
		summary = lipgloss.NewStyle().Foreground(ColorSuccess).Render(locStr)
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", scrollInfo, "", BoxStyle.MaxWidth(60).Render(regionDivider+"\n"+items+localeDivider), "", summary)
}
