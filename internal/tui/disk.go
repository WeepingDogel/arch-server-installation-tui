package tui

import (
	"fmt"
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
	subStep   int // 0=disk select, 1=partition scheme, 2=partition mode
	Next      bool
	GoBack    bool
}

// NewDiskModel creates the disk selection screen with detected disks.
func NewDiskModel(config *model.Config) DiskModel {
	disks, sizes := detectDisks()
	return DiskModel{
		config:    config,
		cursor:    0,
		disks:     disks,
		diskSizes: sizes,
		subStep:   0,
	}
}

func detectDisks() ([]string, map[string]string) {
	sizes := make(map[string]string)
	out, err := exec.Command("lsblk", "-d", "-o", "NAME,SIZE", "-n").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var disks []string
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				name, size := fields[0], fields[1]
				device := "/dev/" + name
				disks = append(disks, device)
				sizes[device] = size
			}
		}
		if len(disks) > 0 {
			return disks, sizes
		}
	}
	out, err = exec.Command("cat", "/proc/partitions").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var disks []string
		for _, line := range lines[2:] {
			fields := strings.Fields(line)
			if len(fields) == 4 {
				name := fields[3]
				isDisk := (strings.HasPrefix(name, "sd") && len(name) == 3) ||
					(strings.HasPrefix(name, "nvme") && strings.Contains(name, "n") && !strings.Contains(name[5:], "p")) ||
					(strings.HasPrefix(name, "vd") && len(name) == 3) ||
					(strings.HasPrefix(name, "mmcblk") && !strings.Contains(name, "p"))
				if isDisk {
					device := "/dev/" + name
					disks = append(disks, device)
					sizes[device] = fields[2] + "K"
				}
			}
		}
		if len(disks) > 0 {
			return disks, sizes
		}
	}
	return []string{"/dev/sda", "/dev/sdb", "/dev/nvme0n1", "/dev/mmcblk0"}, sizes
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
			mx := m.itemCount() - 1
			if m.cursor < mx {
				m.cursor++
			}
		case "enter":
			m.handleSelect()
		case "esc":
			if m.subStep > 0 {
				m.subStep--
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m *DiskModel) itemCount() int {
	switch m.subStep {
	case 0:
		return len(m.disks)
	case 1:
		return 2 // GPT or MBR
	case 2:
		return 3 // Auto, Manual, or confirm
	default:
		return 0
	}
}

func (m *DiskModel) handleSelect() {
	switch m.subStep {
	case 0:
		if m.cursor < len(m.disks) {
			m.config.DiskDevice = m.disks[m.cursor]
			m.config.DiskSize = m.diskSizes[m.disks[m.cursor]]
			m.subStep = 1
			m.cursor = 0
		}
	case 1:
		if m.cursor == 0 {
			m.config.PartitionScheme = "gpt"
		} else {
			m.config.PartitionScheme = "mbr"
		}
		m.subStep = 2
		m.cursor = 0
	case 2:
		if m.cursor == 0 {
			m.config.PartitionMode = "auto"
			m.Next = true
		} else if m.cursor == 1 {
			m.config.PartitionMode = "manual"
			m.Next = true
		}
		// cursor 2 = back (handled by esc)
	}
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
	subtitle := SubtitleStyle.Render("Select the target disk for installation.")
	warning := lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("⚠  All data on the selected disk will be erased!")

	if len(m.disks) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", ErrorBox("No disks detected."))
	}

	var items string
	for i, disk := range m.disks {
		sel := m.config.DiskDevice == disk
		size := m.diskSizes[disk]
		label := disk
		if size != "" {
			label += "  [" + size + "]"
		}
		items += ListItem(i == m.cursor, sel, label) + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", warning, "", BoxStyle.Render(items))
}

func (m DiskModel) viewPartitionScheme(title string) string {
	subtitle := SubtitleStyle.Render("Choose partition table type for " + m.config.DiskDevice)
	info := InfoBox("GPT: required for UEFI, supports 128+ partitions\nMBR: legacy, max 4 primary partitions")

	var items string
	items += ListItem(m.cursor == 0, m.config.PartitionScheme == "gpt", RadioButton(m.config.PartitionScheme == "gpt", "GPT (recommended)")) + "\n"
	items += ListItem(m.cursor == 1, m.config.PartitionScheme == "mbr", RadioButton(m.config.PartitionScheme == "mbr", "MBR / DOS")) + "\n"

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", info, "", BoxStyle.Render(items))
}

func (m DiskModel) viewPartitionMode(title string) string {
	subtitle := SubtitleStyle.Render(fmt.Sprintf("%s on %s (%s)", strings.ToUpper(m.config.PartitionScheme), m.config.DiskDevice, m.config.DiskSize))

	var autoDesc string
	if m.config.PartitionScheme == "gpt" {
		autoDesc = "Auto: 512MB EFI + rest as root"
	} else {
		autoDesc = "Auto: 1MB boot + rest as root"
	}

	var items string
	items += ListItem(m.cursor == 0, m.config.PartitionMode == "auto", RadioButton(m.config.PartitionMode == "auto", "Auto partition")) + "\n"
	items += HelpStyle.Render("  "+autoDesc) + "\n\n"
	items += ListItem(m.cursor == 1, m.config.PartitionMode == "manual", RadioButton(m.config.PartitionMode == "manual", "Manual (advanced)")) + "\n"
	items += HelpStyle.Render("  Custom: set partition sizes manually") + "\n"

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", BoxStyle.Render(items))
}
