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
	diskSizes map[string]string // device -> size
	Next      bool
}

// NewDiskModel creates the disk selection screen with detected disks.
func NewDiskModel(config *model.Config) DiskModel {
	disks, sizes := detectDisks()
	return DiskModel{
		config:    config,
		cursor:    0,
		disks:     disks,
		diskSizes: sizes,
	}
}

// detectDisks probes the system for available block devices and their sizes.
func detectDisks() ([]string, map[string]string) {
	diskSizes := make(map[string]string)

	// Try lsblk first (Linux) with size column
	out, err := exec.Command("lsblk", "-d", "-o", "NAME,SIZE", "-n").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var disks []string
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				name := fields[0]
				size := fields[1]
				device := "/dev/" + name
				disks = append(disks, device)
				diskSizes[device] = size
			}
		}
		if len(disks) > 0 {
			return disks, diskSizes
		}
	}

	// Try /proc/partitions (Linux)
	out, err = exec.Command("cat", "/proc/partitions").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var disks []string
		for _, line := range lines[2:] {
			fields := strings.Fields(line)
			if len(fields) == 4 {
				name := fields[3]
				sizeKB := fields[2] // size in 1K blocks
				isDisk := false
				if strings.HasPrefix(name, "sd") && len(name) == 3 {
					isDisk = true
				} else if strings.HasPrefix(name, "nvme") && strings.Contains(name, "n") && !strings.Contains(name[5:], "p") {
					isDisk = true
				} else if strings.HasPrefix(name, "vd") && len(name) == 3 {
					isDisk = true
				} else if strings.HasPrefix(name, "mmcblk") && !strings.Contains(name, "p") {
					isDisk = true
				}
				if isDisk {
					device := "/dev/" + name
					disks = append(disks, device)
					// Convert KB to human-readable
					diskSizes[device] = formatSize(sizeKB)
				}
			}
		}
		if len(disks) > 0 {
			return disks, diskSizes
		}
	}

	// Fallback
	return []string{"/dev/sda", "/dev/sdb", "/dev/nvme0n1", "/dev/nvme1n1", "/dev/vda", "/dev/mmcblk0"}, diskSizes
}

// formatSize converts a size string in 1K blocks to human-readable format.
func formatSize(sizeKB string) string {
	// Parse the size KB string
	var kb int64
	for _, c := range sizeKB {
		if c >= '0' && c <= '9' {
			kb = kb*10 + int64(c-'0')
		}
	}
	if kb < 1024 {
		return fmt.Sprintf("%dK", kb)
	} else if kb < 1024*1024 {
		return fmt.Sprintf("%.1fM", float64(kb)/1024)
	}
	return fmt.Sprintf("%.1fG", float64(kb)/(1024*1024))
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
			m.config.DiskSize = m.diskSizes[m.disks[m.cursor]]
			m.Next = true
		case "tab":
			m.config.DiskDevice = m.disks[m.cursor]
			m.config.DiskSize = m.diskSizes[m.disks[m.cursor]]
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

	if len(m.disks) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			ErrorBox("No disks detected. Make sure you are running as root on a Linux system."),
		)
	}

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
		size := m.diskSizes[disk]
		sizeStr := ""
		if size != "" {
			sizeStr = lipgloss.NewStyle().Foreground(ColorGray).Render(" [" + size + "]")
		}
		items += style.Render(prefix+disk+sizeStr+selected) + "\n"
	}

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
