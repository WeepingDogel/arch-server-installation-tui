package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TimezoneModel handles timezone and locale configuration.
type TimezoneModel struct {
	config  *model.Config
	cursor  int
	regions []string
	locales []string
	Next    bool
}

// NewTimezoneModel creates the timezone configuration screen.
func NewTimezoneModel(config *model.Config) TimezoneModel {
	return TimezoneModel{
		config:  config,
		cursor:  0,
		regions: []string{"UTC", "Asia/Shanghai", "Asia/Tokyo", "Asia/Seoul", "Asia/Singapore", "Asia/Hong_Kong", "Asia/Taipei", "Asia/Kolkata", "Asia/Dubai", "America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles", "America/Sao_Paulo", "Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Moscow", "Europe/Amsterdam", "Australia/Sydney", "Australia/Melbourne", "Pacific/Auckland", "Africa/Cairo", "Africa/Johannesburg"},
		locales: []string{"en_US.UTF-8", "en_GB.UTF-8", "zh_CN.UTF-8", "zh_TW.UTF-8", "ja_JP.UTF-8", "ko_KR.UTF-8", "de_DE.UTF-8", "fr_FR.UTF-8", "es_ES.UTF-8", "pt_BR.UTF-8", "ru_RU.UTF-8", "it_IT.UTF-8", "pl_PL.UTF-8", "sv_SE.UTF-8", "nb_NO.UTF-8", "da_DK.UTF-8", "fi_FI.UTF-8", "nl_NL.UTF-8", "cs_CZ.UTF-8", "hu_HU.UTF-8", "ro_RO.UTF-8", "bg_BG.UTF-8", "el_GR.UTF-8", "tr_TR.UTF-8"},
	}
}

func (m TimezoneModel) Init() tea.Cmd { return nil }

func (m TimezoneModel) Update(msg tea.Msg) (TimezoneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.regions)+len(m.locales)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor < len(m.regions) {
				m.config.TimezoneRegion = m.regions[m.cursor]
			} else {
				localeIdx := m.cursor - len(m.regions)
				if localeIdx < len(m.locales) {
					m.config.Locale = m.locales[localeIdx]
				}
			}
			// Only advance to next if both timezone and locale are selected
			if m.config.TimezoneRegion != "" && m.config.Locale != "" {
				m.Next = true
			}
		case "tab":
			if m.config.TimezoneRegion != "" && m.config.Locale != "" {
				m.Next = true
			}
		}
	}
	return m, nil
}

func (m TimezoneModel) View() string {
	title := TitleStyle.Render("Timezone & Locale")
	subtitle := SubtitleStyle.Render("Select your timezone and system locale.")

	var items string
	// Show timezone section
	items += DividerStyle.Render(" Timezone Regions ") + "\n\n"
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
		items += style.Render(prefix+region+selected) + "\n"
	}

	// Show locale section
	items += "\n" + DividerStyle.Render(" Locales ") + "\n\n"
	for i, locale := range m.locales {
		idx := i + len(m.regions)
		style := ListItemStyle
		prefix := "  "
		if idx == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if m.config.Locale == locale {
			selected = SuccessStyle.Render(" ✓")
		}
		items += style.Render(prefix+locale+selected) + "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(items),
	)
}