package tui

import (
	"os/exec"
	"strings"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DiskModel handles disk selection and partitioning.
type DiskModel struct {
	config    *model.Config
	cursor    int
	disks     []string
	diskSizes map[string]string
	subStep   int
	Next      bool
	GoBack    bool
}

func (m DiskModel) totalItems() int {
	switch m.subStep {
	case 0:
		return len(m.disks) + 2
	case 1:
		return 2 + 2 // GPT/MBR + [Next] [Back]
	case 2:
		return 2 + 2 // Auto/Manual + [Next] [Back]
	default:
		return 2
	}
}

// NewDiskModel creates the disk selection screen with detected disks.
func NewDiskModel(config *model.Config) DiskModel {
	disks, sizes := detectDisks()
	return DiskModel{
		config:    config,
		disks:     disks,
		diskSizes: sizes,
		subStep:   0,
	}
}

func detectDisks() ([]string, map[string]string) {
	sizes := make(map[string]string)
	out, err := exec.Command("lsblk", "-d", "-o", "NAME,SIZE", "-n").CombinedOutput()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var disks []string
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				d := "/dev/" + fields[0]
				disks = append(disks, d)
				sizes[d] = fields[1]
			}
		}
		if len(disks) > 0 {
			return disks, sizes
		}
	}
	return []string{"/dev/sda", "/dev/sdb", "/dev/nvme0n1"}, sizes
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
			if m.cursor < m.totalItems()-1 {
				m.cursor++
			}
		case "enter":
			total := m.totalItems()
			if m.cursor == total-2 {
				// [Next]
				switch m.subStep {
				case 0:
					if m.cursor < len(m.disks) {
						m.config.DiskDevice = m.disks[m.cursor]
						m.config.DiskSize = m.diskSizes[m.disks[m.cursor]]
						m.subStep = 1
						m.cursor = 0
					}
				case 1:
					m.config.PartitionScheme = "gpt"
					if m.cursor == 1 {
						m.config.PartitionScheme = "mbr"
					}
					m.subStep = 2
					m.cursor = 0
				case 2:
					if m.cursor == 0 {
						m.config.PartitionMode = "auto"
					} else {
						m.config.PartitionMode = "manual"
					}
					m.Next = true
				}
			} else if m.cursor == total-1 {
				if m.subStep > 0 {
					m.subStep--
					m.cursor = 0
				} else {
					m.GoBack = true
				}
			} else if m.subStep == 0 && m.cursor < len(m.disks) {
				m.config.DiskDevice = m.disks[m.cursor]
				m.config.DiskSize = m.diskSizes[m.disks[m.cursor]]
			} else if m.subStep == 1 {
				m.config.PartitionScheme = "gpt"
				if m.cursor == 1 {
					m.config.PartitionScheme = "mbr"
				}
			}
		case "esc":
			if m.subStep > 0 {
				m.subStep--
				m.cursor = 0
			} else {
				m.GoBack = true
			}
		}
	}
	return m, nil
}

func (m DiskModel) View() string {
	title := TitleStyle.Render("Disk Configuration")
	switch m.subStep {
	case 0:
		return m.viewDiskSelect(title)
	case 1:
		return m.viewPartitionScheme(title)
	case 2:
		return m.viewPartitionMode(title)
	default:
		return ""
	}
}

func (m DiskModel) viewDiskSelect(title string) string {
	w := lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("⚠  All data on the selected disk will be erased!")
	if len(m.disks) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, title, "", ErrorBox("No disks detected."))
	}
	var items string
	for i, d := range m.disks {
		sz := m.diskSizes[d]
		label := d
		if sz != "" {
			label += " [" + sz + "]"
		}
		items += ListItem(i == m.cursor, m.config.DiskDevice == d, label) + "\n"
	}
	items += "\n" + renderNavButtons(m.cursor, len(m.disks), len(m.disks)+1)
	return lipgloss.JoinVertical(lipgloss.Left, title, "", w, "", BoxStyle.Render(items))
}

func (m DiskModel) viewPartitionScheme(title string) string {
	s := SubtitleStyle.Render("Partition table for " + m.config.DiskDevice)
	info := InfoBox("GPT: required for UEFI\nMBR: legacy, max 4 partitions")
	var items string
	items += ListItem(m.cursor == 0, m.config.PartitionScheme == "gpt", RadioButton(m.config.PartitionScheme == "gpt", "GPT (recommended)")) + "\n"
	items += ListItem(m.cursor == 1, m.config.PartitionScheme == "mbr", RadioButton(m.config.PartitionScheme == "mbr", "MBR / DOS")) + "\n"
	items += "\n" + renderNavButtons(m.cursor, 2, 3)
	return lipgloss.JoinVertical(lipgloss.Left, title, s, "", info, "", BoxStyle.Render(items))
}

func (m DiskModel) viewPartitionMode(title string) string {
	desc := "Auto: 512MB EFI + rest as root"
	if m.config.PartitionScheme == "mbr" {
		desc = "Auto: 1MB boot + rest as root"
	}
	var items string
	items += ListItem(m.cursor == 0, m.config.PartitionMode == "auto", RadioButton(m.config.PartitionMode == "auto", "Auto partition")) + "\n"
	items += "  " + HelpStyle.Render(desc) + "\n\n"
	items += ListItem(m.cursor == 1, m.config.PartitionMode == "manual", RadioButton(m.config.PartitionMode == "manual", "Manual (advanced)")) + "\n"
	items += "\n" + renderNavButtons(m.cursor, 2, 3)
	return lipgloss.JoinVertical(lipgloss.Left, title, "", BoxStyle.Render(items))
}
