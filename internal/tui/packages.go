package tui

import (
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PackagesModel handles package selection.
type PackagesModel struct {
	config   *model.Config
	cursor   int
	packages []pkgOption
	Next     bool
}

type pkgOption struct {
	name    string
	enabled bool
	field   string // Which config field to toggle
}

// NewPackagesModel creates the package selection screen.
func NewPackagesModel(config *model.Config) PackagesModel {
	return PackagesModel{
		config: config,
		cursor: 0,
		packages: []pkgOption{
			{name: "Kernel: linux", enabled: config.KernelType == "linux", field: "kernel_linux"},
			{name: "Kernel: linux-lts", enabled: config.KernelType == "linux-lts", field: "kernel_lts"},
			{name: "Kernel: linux-zen", enabled: config.KernelType == "linux-zen", field: "kernel_zen"},
			{name: "Kernel: linux-hardened", enabled: config.KernelType == "linux-hardened", field: "kernel_hardened"},
			{name: "base-devel (build tools)", enabled: config.InstallBaseDev, field: "base_devel"},
			{name: "Docker (container runtime)", enabled: config.InstallDocker, field: "docker"},
			{name: "Nginx (web server)", enabled: config.InstallNginx, field: "nginx"},
			{name: "PostgreSQL (database)", enabled: config.InstallPostgres, field: "postgres"},
			{name: "MariaDB (database)", enabled: config.InstallMariaDB, field: "mariadb"},
			{name: "Redis (caching)", enabled: config.InstallRedis, field: "redis"},
			{name: "Fail2ban (security)", enabled: config.InstallFail2ban, field: "fail2ban"},
			{name: "UFW (firewall)", enabled: config.InstallUfw, field: "ufw"},
			{name: "Git (version control)", enabled: config.InstallGit, field: "git"},
			{name: "Vim (editor)", enabled: config.InstallVim, field: "vim"},
		},
	}
}

func (m PackagesModel) Init() tea.Cmd { return nil }

func (m PackagesModel) Update(msg tea.Msg) (PackagesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.packages)-1 {
				m.cursor++
			}
		case " ":
			// Toggle the selected package
			pkg := &m.packages[m.cursor]
			pkg.enabled = !pkg.enabled
			m.applyToggle(pkg)
		case "enter":
			// Apply and continue
			m.applyAll()
			m.Next = true
		case "tab":
			m.applyAll()
			m.Next = true
		}
	}
	return m, nil
}

func (m *PackagesModel) applyToggle(pkg *pkgOption) {
	switch pkg.field {
	case "kernel_linux":
		m.config.KernelType = "linux"
	case "kernel_lts":
		m.config.KernelType = "linux-lts"
	case "kernel_zen":
		m.config.KernelType = "linux-zen"
	case "kernel_hardened":
		m.config.KernelType = "linux-hardened"
	case "base_devel":
		m.config.InstallBaseDev = pkg.enabled
	case "docker":
		m.config.InstallDocker = pkg.enabled
	case "nginx":
		m.config.InstallNginx = pkg.enabled
	case "postgres":
		m.config.InstallPostgres = pkg.enabled
	case "mariadb":
		m.config.InstallMariaDB = pkg.enabled
	case "redis":
		m.config.InstallRedis = pkg.enabled
	case "fail2ban":
		m.config.InstallFail2ban = pkg.enabled
	case "ufw":
		m.config.InstallUfw = pkg.enabled
	case "git":
		m.config.InstallGit = pkg.enabled
	case "vim":
		m.config.InstallVim = pkg.enabled
	}
}

func (m *PackagesModel) applyAll() {
	for i := range m.packages {
		m.applyToggle(&m.packages[i])
	}
}

func (m PackagesModel) View() string {
	title := TitleStyle.Render("Package Selection")
	subtitle := SubtitleStyle.Render("Select packages to install. Use SPACE to toggle, ENTER to continue.")

	var items string
	// Group by category
	currentCategory := ""
	for i, pkg := range m.packages {
		cat := ""
		switch {
		case i < 4:
			cat = " Kernel "
		case i == 4:
			cat = " Development "
		case i >= 5 && i <= 10:
			cat = " Server Packages "
		default:
			cat = " Utilities "
		}

		if cat != currentCategory {
			currentCategory = cat
			items += "\n" + DividerStyle.Render(cat) + "\n\n"
		}

		style := ListItemStyle
		prefix := "  "
		if i == m.cursor {
			style = ListItemSelectedStyle
			prefix = "▶ "
		}
		checkbox := CheckBox(pkg.enabled, "")
		items += style.Render(prefix+checkbox+" "+pkg.name) + "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		BoxStyle.Render(items),
	)
}
