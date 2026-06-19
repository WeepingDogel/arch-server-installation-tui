package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FilesystemModel handles filesystem type selection.
type FilesystemModel struct {
	config  *model.Config
	cursor  int
	fsTypes []fsOption
	Next    bool
}

type fsOption struct {
	name string
	desc string
}

// NewFilesystemModel creates the filesystem selection screen.
func NewFilesystemModel(config *model.Config) FilesystemModel {
	return FilesystemModel{
		config: config,
		cursor: 0,
		fsTypes: []fsOption{
			{name: "ext4", desc: "Most compatible, reliable, default Linux filesystem"},
			{name: "btrfs", desc: "Advanced features: snapshots, compression, subvolumes"},
			{name: "xfs", desc: "High performance, good for large files and servers"},
			{name: "f2fs", desc: "Optimized for flash storage and SSDs"},
		},
	}
}

func (m FilesystemModel) Init() tea.Cmd { return nil }

func (m FilesystemModel) Update(msg tea.Msg) (FilesystemModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.fsTypes)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.config.FilesystemType = m.fsTypes[m.cursor].name
			m.Next = true
		case "tab":
			m.config.FilesystemType = m.fsTypes[m.cursor].name
			m.Next = true
		}
	}
	return m, nil
}

func (m FilesystemModel) View() string {
	title := TitleStyle.Render("Filesystem Selection")
	subtitle := SubtitleStyle.Render("Choose the filesystem type for the root partition.")

	var items string
	for i, fs := range m.fsTypes {
		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		selected := ""
		if m.config.FilesystemType == fs.name {
			selected = SuccessStyle.Render(" ✓")
		}
		desc := lipgloss.NewStyle().Foreground(ColorGray).Render(" — " + fs.desc)
		items += style.Render(prefix+fs.name+desc+selected) + "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(items),
	)
}
