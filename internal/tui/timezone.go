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
	regions []string
	locales []string
	Next    bool
	GoBack  bool
	scroll  int // scroll offset for low-resolution
}

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
	total := len(m.regions) + len(m.locales)
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
				if m.cursor >= m.scroll+tzViewportHeight {
					m.scroll++
				}
			}
		case " ":
			if m.cursor >= len(m.regions) {
				localeIdx := m.cursor - len(m.regions)
				if localeIdx < len(m.locales) {
					m.config.Locales = toggleLocale(m.config.Locales, m.locales[localeIdx])
				}
			}
		case "enter":
			if m.cursor < len(m.regions) {
				m.config.TimezoneRegion = m.regions[m.cursor]
			}
			if m.config.TimezoneRegion != "" {
				m.Next = true
			}
		case "tab":
			if m.config.TimezoneRegion != "" {
				m.Next = true
			}
		}
	}
	return m, nil
}

func (m TimezoneModel) View() string {
	title := TitleStyle.Render("Timezone & Locale")
	subtitle := SubtitleStyle.Render("Select timezone. SPACE to toggle multiple locales, ENTER to confirm.")

	total := len(m.regions) + len(m.locales)

	// Scroll indicator
	visibleEnd := m.scroll + tzViewportHeight
	if visibleEnd > total {
		visibleEnd = total
	}
	scrollInfo := ""
	if total > tzViewportHeight {
		scrollInfo = lipgloss.NewStyle().Foreground(ColorGray).Render(
			fmt.Sprintf("  [%d-%d of %d]  ▼", m.scroll+1, visibleEnd, total),
		)
	}

	var items string
	// Build all items first, then pick the visible window
	var allLines []string

	// Timezone section header
	allLines = append(allLines, DividerStyle.Render(" Timezone "))

	for i, region := range m.regions {
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if m.config.TimezoneRegion == region {
			selected = SuccessStyle.Render(" ✓")
		}
		allLines = append(allLines, style.Render(prefix+region+selected))
	}

	// Locale section header
	allLines = append(allLines, "", DividerStyle.Render(" Locales (SPACE to toggle) "))

	for i, locale := range m.locales {
		idx := i + len(m.regions)
		style := ListItemStyle
		prefix := "  "
		if idx == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		checked := ""
		if localeSelected(m.config.Locales, locale) {
			checked = SuccessStyle.Render("☑")
		}
		allLines = append(allLines, style.Render(prefix+locale+checked))
	}

	// Render only the visible window
	end := m.scroll + tzViewportHeight
	if end > len(allLines) {
		end = len(allLines)
	}
	for i := m.scroll; i < end; i++ {
		items += allLines[i] + "\n"
	}

	// Selection summary
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

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		scrollInfo,
		"",
		BoxStyle.Render(items),
		"",
		summary,
	)
}
