package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BootloaderModel handles bootloader selection.
type BootloaderModel struct {
	config *model.Config
	cursor int
	Next   bool
}

// NewBootloaderModel creates the bootloader selection screen.
func NewBootloaderModel(config *model.Config) BootloaderModel {
	cursor := 0
	if config.BootloaderType == "systemd-boot" {
		cursor = 1
	}
	return BootloaderModel{config: config, cursor: cursor}
}

func (m BootloaderModel) Init() tea.Cmd { return nil }

func (m BootloaderModel) Update(msg tea.Msg) (BootloaderModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 3 {
				m.cursor++
			}
		case "enter", " ":
			switch m.cursor {
			case 0:
				m.config.BootloaderType = "grub"
			case 1:
				m.config.BootloaderType = "systemd-boot"
			case 2:
				m.config.UEFIMode = true
			case 3:
				m.config.UEFIMode = false
			}
			if m.cursor < 2 {
				m.Next = true
			}
		case "tab":
			if m.cursor < 2 {
				m.Next = true
			}
		}
	}
	return m, nil
}

func (m BootloaderModel) View() string {
	title := TitleStyle.Render("Bootloader Configuration")
	subtitle := SubtitleStyle.Render("Select bootloader and firmware mode.")

	bootloaderOptions := []struct {
		name string
		desc string
		sel  string
	}{
		{name: "GRUB", desc: "Most compatible, supports BIOS and UEFI", sel: m.config.BootloaderType},
		{name: "systemd-boot", desc: "Simple, UEFI only (recommended for UEFI)", sel: m.config.BootloaderType},
	}

	var items string
	for i, opt := range bootloaderOptions {
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if m.config.BootloaderType == opt.name {
			selected = SuccessStyle.Render(" ✓")
		}
		desc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — " + opt.desc)
		items += style.Render(prefix+opt.name+desc+selected) + "\n"
	}

	firmwareModes := []struct {
		name string
		sel  bool
		desc string
	}{
		{name: "UEFI Mode", sel: m.config.UEFIMode, desc: "Modern firmware, faster boot, secure boot support"},
		{name: "BIOS/Legacy Mode", sel: !m.config.UEFIMode, desc: "Traditional BIOS compatibility"},
	}

	items += "\n" + DividerStyle.Render(" Firmware Mode ") + "\n\n"
	for i, fm := range firmwareModes {
		idx := i + 2
		style := ListItemStyle
		prefix := "  "
		if idx == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if fm.sel {
			selected = SuccessStyle.Render(" ✓")
		}
		desc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — " + fm.desc)
		items += style.Render(prefix+fm.name+desc+selected) + "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(items),
	)
}
