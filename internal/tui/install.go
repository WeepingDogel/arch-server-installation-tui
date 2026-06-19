package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/installer"
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InstallModel shows the installation progress with log viewer.
type InstallModel struct {
	config       *model.Config
	spinnerFrame int
	completed    bool
	currentStep  int
	installer    *installer.Installer
	showLogs     bool
	logs         []string
}

var installSteps = []string{
	"Partitioning disk...",
	"Formatting filesystems...",
	"Mounting partitions...",
	"Installing base system...",
	"Generating fstab...",
	"Configuring timezone & locale...",
	"Setting up hostname...",
	"Configuring network...",
	"Setting root password...",
	"Creating user account...",
	"Installing bootloader...",
	"Configuring SSH...",
	"Installing additional packages...",
	"Finalizing installation...",
	"Installation complete!",
}

// NewInstallModel creates the installation progress screen.
func NewInstallModel(config *model.Config) InstallModel {
	return InstallModel{
		config:       config,
		spinnerFrame: 0,
		installer:    installer.New(config),
		logs:         make([]string, 0, 100),
	}
}

func (m InstallModel) Init() tea.Cmd { return nil }

// StartInstall returns the command to begin installation.
func (m InstallModel) StartInstall() tea.Cmd {
	m.config.InstallStarted = true
	m.logs = append(m.logs, "Installation started...")
	return func() tea.Msg {
		progressCh := make(chan installer.ProgressUpdate)
		go m.installer.Install(progressCh)

		// Read first progress message
		update, ok := <-progressCh
		if !ok {
			return installProgressMsg{
				Percent: 100,
				Message: "Installation complete!",
				Done:    true,
			}
		}
		go func() {
			for p := range progressCh {
				_ = p
			}
		}()
		return installProgressMsg{
			Percent:   update.Percent,
			Message:   update.Message,
			LogOutput: update.Message,
			Done:      update.Done,
			Err:       update.Err,
		}
	}
}

func (m InstallModel) Update(msg tea.Msg) (InstallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case installProgressMsg:
		m.currentStep++
		m.spinnerFrame = (m.spinnerFrame + 1) % 10
		m.config.ProgressPercent = msg.Percent
		m.config.ProgressMessage = msg.Message

		if msg.LogOutput != "" {
			m.logs = append(m.logs, msg.LogOutput)
		}

		if msg.Done {
			m.completed = true
			m.config.InstallComplete = true
			m.logs = append(m.logs, "Installation complete!")
			return m, nil
		}

		// Schedule next progress poll
		return m, func() tea.Msg {
			time.Sleep(600 * time.Millisecond)
			return installProgressMsg{
				Percent:   m.config.ProgressPercent + 100.0/float64(len(installSteps)),
				Message:   fmt.Sprintf("Installing... (%.0f%%)", m.config.ProgressPercent),
				LogOutput: "",
				Done:      m.config.ProgressPercent >= 99,
			}
		}

	case tea.KeyMsg:
		if m.completed {
			switch msg.String() {
			case "enter", "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		// Toggle log view with 'l'
		if msg.String() == "l" {
			m.showLogs = !m.showLogs
			return m, nil
		}
	}

	return m, nil
}

func (m InstallModel) View() string {
	if m.completed {
		return m.completedView()
	}
	if m.showLogs {
		return m.logView()
	}
	return m.inProgressView()
}

func (m InstallModel) inProgressView() string {
	logo := MiniArchLogo()
	title := TitleStyle.Render("Installing Arch Linux Server...")

	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinnerFrames[m.spinnerFrame]

	spinnerView := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(spinner + " " + m.config.ProgressMessage)

	barWidth := 50
	filled := int(m.config.ProgressPercent / 100.0 * float64(barWidth))
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar = lipgloss.NewStyle().Foreground(ColorPrimary).Render(bar)
	percentText := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(fmt.Sprintf("%.0f%%", m.config.ProgressPercent))

	logHint := HelpStyle.Render("Press 'l' to view full logs")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		title,
		"",
		BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, spinnerView, "", bar, "", percentText)),
		"",
		logHint,
	)

	return Screen(TotalSteps, content, "Installation in progress...", 1)
}

// logView shows the full installation log output.
func (m InstallModel) logView() string {
	title := TitleStyle.Render("Installation Logs")
	subtitle := SubtitleStyle.Render("Press 'l' to return to progress view. ↑/↓ to scroll.")

	var logContent string
	// Show last 20 log lines
	start := 0
	if len(m.logs) > 20 {
		start = len(m.logs) - 20
	}
	for i, line := range m.logs[start:] {
		prefix := lipgloss.NewStyle().Foreground(ColorGray).Render(fmt.Sprintf("%3d ", start+i+1))
		logContent += prefix + line + "\n"
	}

	if len(m.logs) == 0 {
		logContent = HelpStyle.Render("Waiting for installation output...")
	}

	logBox := BoxStyle.Render(strings.TrimRight(logContent, "\n"))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		logBox,
	)
}

func (m InstallModel) completedView() string {
	logo := ArchLogo()
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSuccess).
		Render("✓ Installation Complete!")
	subtitle := SubtitleStyle.Render("Arch Linux Server has been successfully installed.")

	summaryInfo := BoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf("Disk: %s", m.config.DiskDevice),
			fmt.Sprintf("Partition: %s / %s", m.config.PartitionScheme, m.config.PartitionMode),
			fmt.Sprintf("Filesystem: %s", m.config.FilesystemType),
			fmt.Sprintf("Bootloader: %s (%s)", m.config.BootloaderType, bootModeStr(m.config.UEFIMode)),
			fmt.Sprintf("Hostname: %s", m.config.Hostname),
			fmt.Sprintf("SSH: Port %d, Root Login: %s", m.config.SSHPort, boolStr(m.config.AllowRootLogin, "Allowed", "Denied")),
			"",
			SubtitleStyle.Render("Reboot to start your new Arch Linux server!"),
		),
	)

	quitHint := lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true).
		Render("Press ENTER, Q or Ctrl+C to quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		logo,
		"",
		title,
		"",
		subtitle,
		"",
		summaryInfo,
		"",
		quitHint,
	)
}
