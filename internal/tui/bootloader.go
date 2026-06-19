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
		case "enter":
			switch m.cursor {
			case 0:
				m.config.BootloaderType = "grub"
				m.Next = true
			case 1:
				m.config.BootloaderType = "systemd-boot"
				m.Next = true
			case 2:
				m.config.UEFIMode = true
			case 3:
				m.config.UEFIMode = false
			}
		}
	}
	return m, nil
}

func (m BootloaderModel) View() string {
	title := TitleStyle.Render("Bootloader Configuration")
	subtitle := SubtitleStyle.Render("Select bootloader and firmware mode.")

	var items string
	// Bootloader section
	items += DividerStyle.Render(" Bootloader ") + "\n\n"
	isGRUB := m.config.BootloaderType == "grub"
	isSDBoot := m.config.BootloaderType == "systemd-boot"
	grubDesc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — Supports BIOS and UEFI")
	sdDesc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — UEFI only, simpler")

	items += ListItem(m.cursor == 0, isGRUB, RadioButton(isGRUB, "GRUB"+grubDesc)) + "\n"
	items += ListItem(m.cursor == 1, isSDBoot, RadioButton(isSDBoot, "systemd-boot"+sdDesc)) + "\n"

	// Firmware section
	items += "\n" + DividerStyle.Render(" Firmware Mode ") + "\n\n"
	isUEFI := m.config.UEFIMode
	isBIOS := !m.config.UEFIMode
	uefiDesc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — Modern, secure boot, GPT")
	biosDesc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — Legacy, MBR")

	items += ListItem(m.cursor == 2, isUEFI, RadioButton(isUEFI, "UEFI"+uefiDesc)) + "\n"
	items += ListItem(m.cursor == 3, isBIOS, RadioButton(isBIOS, "BIOS/Legacy"+biosDesc)) + "\n"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(items),
	)
}
