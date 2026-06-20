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
	GoBack bool
}

func (m BootloaderModel) totalItems() int {
	return 4 + 2 // 4 options + [Next] [Back]
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
			if m.cursor < m.totalItems()-1 {
				m.cursor++
			}
		case "enter":
			total := m.totalItems()
			if m.cursor == total-2 {
				if m.config.BootloaderType == "" {
					m.config.BootloaderType = "grub"
				}
				m.Next = true
			} else if m.cursor == total-1 {
				m.GoBack = true
			} else if m.cursor < 2 {
				// Bootloader choice
				if m.cursor == 0 {
					m.config.BootloaderType = "grub"
				} else {
					m.config.BootloaderType = "systemd-boot"
				}
			} else {
				// Firmware mode
				if m.cursor == 2 {
					m.config.UEFIMode = true
				} else {
					m.config.UEFIMode = false
				}
			}
		}
	}
	return m, nil
}

func (m BootloaderModel) View() string {
	title := TitleStyle.Render("Bootloader Configuration")
	subtitle := SubtitleStyle.Render("↑/↓ select, ENTER on [Next] to confirm.")

	var items string
	items += DividerStyle.Render(" Bootloader ") + "\n\n"
	isGRUB := m.config.BootloaderType == "grub"
	isSDBoot := m.config.BootloaderType == "systemd-boot"
	items += ListItem(m.cursor == 0, isGRUB, RadioButton(isGRUB, "GRUB — Supports BIOS and UEFI")) + "\n"
	items += ListItem(m.cursor == 1, isSDBoot, RadioButton(isSDBoot, "systemd-boot — UEFI only, simpler")) + "\n"

	items += "\n" + DividerStyle.Render(" Firmware Mode ") + "\n\n"
	items += ListItem(m.cursor == 2, m.config.UEFIMode, RadioButton(m.config.UEFIMode, "UEFI — Modern, secure boot")) + "\n"
	items += ListItem(m.cursor == 3, !m.config.UEFIMode, RadioButton(!m.config.UEFIMode, "BIOS/Legacy — MBR mode")) + "\n"

	items += "\n" + renderNavButtons(m.cursor, 4, 5)

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", BoxStyle.Render(items))
}
