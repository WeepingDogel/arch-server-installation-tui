package tui

import (
	"fmt"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RootModel is the top-level Bubble Tea model that manages navigation
// between all installation wizard steps.
type RootModel struct {
	config *model.Config
	step   int // 1-based step number

	// Sub-models for each step
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

	// Window size
	width  int
	height int

	err error
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

// Init initializes the Bubble Tea program.
func (m *RootModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles messages and key events.
func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step == TotalSteps && m.config.InstallStarted {
				return m, nil
			}
			return m, tea.Quit

		case "esc":
			if m.step > 1 {
				m.step--
			}
			return m, nil
		}
	}

	// Delegate to current step's model first (so it can save inputs)
	var cmd tea.Cmd
	switch m.step {
	case 1:
		newModel, c := m.welcome.Update(msg)
		m.welcome = newModel
		cmd = c
	case 2:
		newModel, c := m.keyboard.Update(msg)
		m.keyboard = newModel
		cmd = c
	case 3:
		newModel, c := m.network.Update(msg)
		m.network = newModel
		cmd = c
	case 4:
		newModel, c := m.mirror.Update(msg)
		m.mirror = newModel
		cmd = c
	case 5:
		newModel, c := m.disk.Update(msg)
		m.disk = newModel
		cmd = c
	case 6:
		newModel, c := m.filesystem.Update(msg)
		m.filesystem = newModel
		cmd = c
	case 7:
		newModel, c := m.bootloader.Update(msg)
		m.bootloader = newModel
		cmd = c
	case 8:
		newModel, c := m.timezone.Update(msg)
		m.timezone = newModel
		cmd = c
	case 9:
		newModel, c := m.users.Update(msg)
		m.users = newModel
		cmd = c
	case 10:
		newModel, c := m.ssh.Update(msg)
		m.ssh = newModel
		cmd = c
	case 11:
		newModel, c := m.packages.Update(msg)
		m.packages = newModel
		cmd = c
	case 12:
		newModel, c := m.summary.Update(msg)
		m.summary = newModel
		cmd = c
	case 13:
		newModel, c := m.install.Update(msg)
		m.install = newModel
		cmd = c
	}

	// After sub-model processed the message, check if it wants to advance
	// and validate the step before advancing
	if m.err != nil {
		m.err = nil // clear previous error on new attempt
	}

	switch m.step {
	case 1:
		if m.welcome.Next {
			m.welcome.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 2:
		if m.keyboard.Next {
			m.keyboard.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 3:
		if m.network.Next {
			m.network.Next = false
			if err := m.validateStep(); err != nil {
				m.err = err
			} else if m.step < TotalSteps {
				m.step++
			}
		}
	case 4:
		if m.mirror.Next {
			m.mirror.Next = false
			if err := m.validateStep(); err != nil {
				m.err = err
			} else if m.step < TotalSteps {
				m.step++
			}
		}
	case 5:
		if m.disk.Next {
			m.disk.Next = false
			if err := m.validateStep(); err != nil {
				m.err = err
			} else if m.step < TotalSteps {
				m.step++
			}
		}
	case 6:
		if m.filesystem.Next {
			m.filesystem.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 7:
		if m.bootloader.Next {
			m.bootloader.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 8:
		if m.timezone.Next {
			m.timezone.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 9:
		if m.users.Next {
			m.users.Next = false
			if err := m.validateStep(); err != nil {
				m.err = err
			} else if m.step < TotalSteps {
				m.step++
			}
		}
	case 10:
		if m.ssh.Next {
			m.ssh.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 11:
		if m.packages.Next {
			m.packages.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
	case 12:
		if m.summary.Next {
			m.summary.Next = false
			if m.step < TotalSteps {
				m.step++
			}
		}
		if m.summary.Confirmed {
			m.summary.Confirmed = false
			m.step = TotalSteps
			return m, m.install.StartInstall()
		}
	}

	return m, cmd
}

// View renders the current step's UI.
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

// validateStep validates the current step's input before allowing navigation.
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
	Percent float64
	Message string
	Done    bool
	Err     error
}
