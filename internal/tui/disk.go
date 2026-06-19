package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DiskModel handles disk selection and partitioning.
type DiskModel struct {
	config *model.Config
	cursor int
	disks  []string
	Next   bool
}

// NewDiskModel creates the disk selection screen.
func NewDiskModel(config *model.Config) DiskModel {
	return DiskModel{
		config: config,
		cursor: 0,
		disks:  []string{"/dev/sda", "/dev/sdb", "/dev/nvme0n1", "/dev/nvme1n1", "/dev/vda", "/dev/mmcblk0"},
	}
}

func (m DiskModel) Init() tea.Cmd { return nil }

func (m DiskModel) Update(msg tea.Msg) (DiskModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.disks)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.config.DiskDevice = m.disks[m.cursor]
			m.Next = true
		case "tab":
			m.config.DiskDevice = m.disks[m.cursor]
			m.Next = true
		}
	}
	return m, nil
}

func (m DiskModel) View() string {
	title := TitleStyle.Render("Disk Selection")
	subtitle := SubtitleStyle.Render("Select the target disk for installation.")

	warning := lipgloss.NewStyle().
		Foreground(ColorWarning).
		Bold(true).
		Render("⚠  WARNING: All data on the selected disk will be erased!")

	var items string
	for i, disk := range m.disks {
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if m.config.DiskDevice == disk {
			selected = SuccessStyle.Render(" ✓")
		}
		items += style.Render(prefix+disk+selected) + "\n"
	}

	// Auto-partitioning toggle info
	partitionInfo := BoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			CheckBox(m.config.UseAutoPartitioning, "Auto partitioning (recommended)"),
			"",
			InfoBox("Auto partitioning creates: 512MB EFI + rest as root"),
		),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		warning,
		"",
		BoxStyle.Render(items),
		"",
		partitionInfo,
	)
}