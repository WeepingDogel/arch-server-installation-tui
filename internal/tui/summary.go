package tui

import (
	"fmt"
	"strings"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SummaryModel displays the installation summary and confirmation.
type SummaryModel struct {
	config    *model.Config
	Next      bool
	GoBack    bool
	Confirmed bool
	cursor    int // 0 = Confirm, 1 = Back
}

// NewSummaryModel creates the summary screen.
func NewSummaryModel(config *model.Config) SummaryModel {
	return SummaryModel{config: config}
}

func (m SummaryModel) Init() tea.Cmd { return nil }

func (m SummaryModel) Update(msg tea.Msg) (SummaryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "down", "j":
			m.cursor = 1 - m.cursor
		case "enter":
			if m.cursor == 0 {
				m.Confirmed = true
			} else {
				m.Next = true
			}
		case "tab":
			m.Confirmed = true
		case "esc":
			m.Next = true
		}
	}
	return m, nil
}

func (m SummaryModel) View() string {
	title := TitleStyle.Render("Installation Summary")
	subtitle := SubtitleStyle.Render("Review your configuration before installation.")

	// Locale display
	localeStr := "en_US.UTF-8"
	if len(m.config.Locales) > 0 {
		localeStr = strings.Join(m.config.Locales, ", ")
	}

	sections := []struct {
		name  string
		items []string
	}{
		{
			name: "System",
			items: []string{
				fmt.Sprintf("Keyboard Layout: %s", m.config.KeyboardLayout),
				fmt.Sprintf("Hostname: %s", m.config.Hostname),
				fmt.Sprintf("Timezone: %s", m.config.TimezoneRegion),
				fmt.Sprintf("Locales: %s", localeStr),
			},
		},
		{
			name: "Network",
			items: []string{
				fmt.Sprintf("Interface: %s", ifElse(m.config.NetworkIface != "", m.config.NetworkIface, "auto")),
				fmt.Sprintf("Mode: %s", boolStr(m.config.NetworkDHCP, "DHCP", "Static")),
			},
		},
		{
			name: "Mirror",
			items: []string{
				fmt.Sprintf("Main Mirror: %s", m.config.MirrorURL),
				fmt.Sprintf("Arch Linux CN: %s", boolStr(m.config.EnableArchCN, "Enabled", "Disabled")),
			},
		},
		{
			name: "Storage",
			items: []string{
				fmt.Sprintf("Disk: %s (%s)", m.config.DiskDevice, m.config.DiskSize),
				fmt.Sprintf("Filesystem: %s", m.config.FilesystemType),
				fmt.Sprintf("Bootloader: %s (%s)", m.config.BootloaderType, bootModeStr(m.config.UEFIMode)),
			},
		},
		{
			name: "Users",
			items: []string{
				"Root Password: " + passwordMask(m.config.RootPassword),
				fmt.Sprintf("Sudo User: %s", ifElse(m.config.CreateUser, m.config.UserName, "(none)")),
			},
		},
		{
			name: "SSH",
			items: []string{
				fmt.Sprintf("SSH: %s", boolStr(m.config.EnableSSH, "Enabled", "Disabled")),
				fmt.Sprintf("Port: %d", m.config.SSHPort),
				fmt.Sprintf("Root Login: %s", boolStr(m.config.AllowRootLogin, "Allowed", "Denied")),
			},
		},
		{
			name: "Packages",
			items: []string{
				fmt.Sprintf("Kernel: %s", m.config.KernelType),
				packageLine(m.config.InstallDocker, "Docker"),
				packageLine(m.config.InstallNginx, "Nginx"),
				packageLine(m.config.InstallPostgres, "PostgreSQL"),
				packageLine(m.config.InstallMariaDB, "MariaDB"),
				packageLine(m.config.InstallRedis, "Redis"),
				packageLine(m.config.InstallFail2ban, "Fail2ban"),
				packageLine(m.config.InstallUfw, "UFW"),
				packageLine(m.config.InstallGit, "Git"),
				packageLine(m.config.InstallVim, "Vim"),
			},
		},
	}

	var summary string
	for _, section := range sections {
		summary += "\n" + lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(section.name) + "\n"
		for _, item := range section.items {
			summary += "  " + item + "\n"
		}
	}

	summaryBox := BoxStyle.Render(summary)

	warning := lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true).
		Render("⚠  WARNING: This will erase all data on " + m.config.DiskDevice + "!")

	var confirmBtn string
	if m.cursor == 0 {
		confirmBtn = ButtonStyle.Render("▶ Start Installation")
	} else {
		confirmBtn = ButtonDisabledStyle.Render("  Start Installation  ")
	}

	var backBtn string
	if m.cursor == 1 {
		backBtn = BackButtonStyle.Render("◀ Back to Edit")
	} else {
		backBtn = BackButtonStyle.Render("  Back to Edit  ")
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		backBtn,
		confirmBtn,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		warning,
		"",
		summaryBox,
		"",
		buttons,
	)
}

func boolStr(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func ifElse(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func passwordMask(pw string) string {
	if pw == "" {
		return "(not set)"
	}
	return "********"
}

func bootModeStr(uefi bool) string {
	if uefi {
		return "UEFI"
	}
	return "BIOS/Legacy"
}

func packageLine(enabled bool, name string) string {
	if enabled {
		return lipgloss.NewStyle().Foreground(ColorSuccess).Render("✓ " + name)
	}
	return lipgloss.NewStyle().Foreground(ColorGray).Render("○ " + name)
}
