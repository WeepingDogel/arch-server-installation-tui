package tui

import (
	"fmt"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RootModel is the top-level Bubble Tea model that manages navigation.
type RootModel struct {
	config *model.Config
	step   int

	welcome    WelcomeModel
	keyboard   KeyboardModel
	network    NetworkModel
	mirror     MirrorModel
	disk       DiskModel
	filesystem FilesystemModel
	bootloader BootloaderModel
	timezone   TimezoneModel
	users      UsersModel
	ssh        SSHModel
	packages   PackagesModel
	summary    SummaryModel
	install    InstallModel

	width  int
	height int
	err    error
}

// New creates the root model with default configuration.
func New() *RootModel {
	cfg := model.DefaultConfig()
	return &RootModel{
		config:     cfg,
		step:       1,
		welcome:    NewWelcomeModel(cfg),
		keyboard:   NewKeyboardModel(cfg),
		network:    NewNetworkModel(cfg),
		mirror:     NewMirrorModel(cfg),
		disk:       NewDiskModel(cfg),
		filesystem: NewFilesystemModel(cfg),
		bootloader: NewBootloaderModel(cfg),
		timezone:   NewTimezoneModel(cfg),
		users:      NewUsersModel(cfg),
		ssh:        NewSSHModel(cfg),
		packages:   NewPackagesModel(cfg),
		summary:    NewSummaryModel(cfg),
		install:    NewInstallModel(cfg),
	}
}

func (m *RootModel) Init() tea.Cmd { return tea.EnterAltScreen }

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		}
	}

	var cmd tea.Cmd
	switch m.step {
	case 1:
		newModel, c := m.welcome.Update(msg)
		m.welcome = newModel
		cmd = c
		if m.welcome.Next {
			m.welcome.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 2:
		newModel, c := m.keyboard.Update(msg)
		m.keyboard = newModel
		cmd = c
		if m.keyboard.Next || m.keyboard.GoBack {
			m.handleNav(&m.keyboard.Next, &m.keyboard.GoBack)
			return m, cmd
		}
	case 3:
		newModel, c := m.network.Update(msg)
		m.network = newModel
		cmd = c
		if m.network.Next || m.network.GoBack {
			m.handleNav(&m.network.Next, &m.network.GoBack)
			return m, cmd
		}
	case 4:
		newModel, c := m.mirror.Update(msg)
		m.mirror = newModel
		cmd = c
		if m.mirror.Next || m.mirror.GoBack {
			m.handleNav(&m.mirror.Next, &m.mirror.GoBack)
			return m, cmd
		}
	case 5:
		newModel, c := m.disk.Update(msg)
		m.disk = newModel
		cmd = c
		if m.disk.Next || m.disk.GoBack {
			m.handleNav(&m.disk.Next, &m.disk.GoBack)
			return m, cmd
		}
	case 6:
		newModel, c := m.filesystem.Update(msg)
		m.filesystem = newModel
		cmd = c
		if m.filesystem.Next || m.filesystem.GoBack {
			m.handleNav(&m.filesystem.Next, &m.filesystem.GoBack)
			return m, cmd
		}
	case 7:
		newModel, c := m.bootloader.Update(msg)
		m.bootloader = newModel
		cmd = c
		if m.bootloader.Next || m.bootloader.GoBack {
			m.handleNav(&m.bootloader.Next, &m.bootloader.GoBack)
			return m, cmd
		}
	case 8:
		newModel, c := m.timezone.Update(msg)
		m.timezone = newModel
		cmd = c
		if m.timezone.Next || m.timezone.GoBack {
			m.handleNav(&m.timezone.Next, &m.timezone.GoBack)
			return m, cmd
		}
	case 9:
		newModel, c := m.users.Update(msg)
		m.users = newModel
		cmd = c
		if m.users.Next || m.users.GoBack {
			m.handleNav(&m.users.Next, &m.users.GoBack)
			return m, cmd
		}
	case 10:
		newModel, c := m.ssh.Update(msg)
		m.ssh = newModel
		cmd = c
		if m.ssh.Next || m.ssh.GoBack {
			m.handleNav(&m.ssh.Next, &m.ssh.GoBack)
			return m, cmd
		}
	case 11:
		newModel, c := m.packages.Update(msg)
		m.packages = newModel
		cmd = c
		if m.packages.Next || m.packages.GoBack {
			m.handleNav(&m.packages.Next, &m.packages.GoBack)
			return m, cmd
		}
	case 12:
		newModel, c := m.summary.Update(msg)
		m.summary = newModel
		cmd = c
		if m.summary.Next || m.summary.GoBack {
			m.handleNav(&m.summary.Next, &m.summary.GoBack)
			return m, cmd
		}
		if m.summary.Confirmed {
			m.summary.Confirmed = false
			m.step = TotalSteps
			return m, m.install.StartInstall()
		}
	case 13:
		newModel, c := m.install.Update(msg)
		m.install = newModel
		cmd = c
	}

	return m, cmd
}

// handleNav processes GoBack/Next from sub-models with validation.
func (m *RootModel) handleNav(next, back *bool) {
	if *back {
		*back = false
		if m.step > 1 {
			m.step--
		}
		return
	}
	if *next {
		*next = false
		if err := m.validateStep(); err != nil {
			m.err = err
			return
		}
		m.err = nil
		if m.step < TotalSteps {
			m.step++
		}
	}
}

func (m *RootModel) View() string {
	var content string
	switch m.step {
	case 1:
		content = m.welcome.View()
	case 2:
		content = m.keyboard.View()
	case 3:
		content = m.network.View()
	case 4:
		content = m.mirror.View()
	case 5:
		content = m.disk.View()
	case 6:
		content = m.filesystem.View()
	case 7:
		content = m.bootloader.View()
	case 8:
		content = m.timezone.View()
	case 9:
		content = m.users.View()
	case 10:
		content = m.ssh.View()
	case 11:
		content = m.packages.View()
	case 12:
		content = m.summary.View()
	case 13:
		return m.install.View()
	}

	errMsg := ""
	if m.err != nil {
		errMsg = ErrorBox(m.err.Error())
	}

	screen := Screen(m.step, content, SimpleFooter())

	if errMsg != "" {
		screen = lipgloss.JoinVertical(lipgloss.Top, screen, "", errMsg)
	}

	return fmt.Sprintf("\n%s\n", screen)
}

func (m *RootModel) validateStep() error {
	switch m.step {
	case 3:
		if !m.config.NetworkDHCP {
			if !model.ValidateHostname(m.config.Hostname) {
				return fmt.Errorf("invalid hostname: must be 1-63 chars, alphanumeric, hyphens, dots")
			}
			if !model.ValidateIP(m.config.IPAddress) {
				return fmt.Errorf("invalid IP address")
			}
			if !model.ValidateIP(m.config.Gateway) {
				return fmt.Errorf("invalid gateway address")
			}
		} else {
			if !model.ValidateHostname(m.config.Hostname) {
				return fmt.Errorf("invalid hostname: must be 1-63 chars, alphanumeric, hyphens, dots")
			}
		}
	case 4:
		if m.config.MirrorURL == "" {
			return fmt.Errorf("please select a mirror")
		}
	case 5:
		if m.config.DiskDevice == "" {
			return fmt.Errorf("please select a disk device")
		}
	case 9:
		if m.config.RootPassword == "" {
			return fmt.Errorf("root password is required")
		}
		if len(m.config.RootPassword) < 8 {
			return fmt.Errorf("root password must be at least 8 characters")
		}
		if m.config.CreateUser && m.config.UserName != "" && m.config.UserPassword == "" {
			return fmt.Errorf("user password is required")
		}
	}
	return nil
}

// installProgressMsg is sent by the installer to update the UI.
type installProgressMsg struct {
	Percent   float64
	Message   string
	LogOutput string
	Stage     string
	Done      bool
	Err       error
}
