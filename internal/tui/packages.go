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
	GoBack   bool
}

type pkgOption struct {
	name    string
	enabled bool
	field   string
}

func (m PackagesModel) totalItems() int {
	return len(m.packages) + 2
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
			{name: "Docker", enabled: config.InstallDocker, field: "docker"},
			{name: "Nginx", enabled: config.InstallNginx, field: "nginx"},
			{name: "PostgreSQL", enabled: config.InstallPostgres, field: "postgres"},
			{name: "MariaDB", enabled: config.InstallMariaDB, field: "mariadb"},
			{name: "Redis", enabled: config.InstallRedis, field: "redis"},
			{name: "Fail2ban", enabled: config.InstallFail2ban, field: "fail2ban"},
			{name: "UFW", enabled: config.InstallUfw, field: "ufw"},
			{name: "Git", enabled: config.InstallGit, field: "git"},
			{name: "Vim", enabled: config.InstallVim, field: "vim"},
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
			if m.cursor < m.totalItems()-1 {
				m.cursor++
			}
		case " ":
			if m.cursor < len(m.packages) {
				pkg := &m.packages[m.cursor]
				pkg.enabled = !pkg.enabled
				m.applyToggle(pkg)
			}
		case "enter":
			total := m.totalItems()
			if m.cursor == total-2 {
				m.applyAll()
				m.Next = true
			} else if m.cursor == total-1 {
				m.GoBack = true
			}
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
	subtitle := SubtitleStyle.Render("SPACE to toggle, ENTER on [Next] to confirm.")

	var items string
	for i, pkg := range m.packages {
		items += ListItem(i == m.cursor, false, Checkbox(pkg.enabled, pkg.name)) + "\n"
	}
	items += "\n" + renderNavButtons(m.cursor, len(m.packages), len(m.packages)+1)

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", BoxStyle.Render(items))
}
