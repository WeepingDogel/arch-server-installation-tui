package tui

import (
	"fmt"
	"time"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InstallModel shows the installation progress.
type InstallModel struct {
	config         *model.Config
	spinnerFrame   int
	completed      bool
	err            error
	currentStep    int
}

// Installation steps shown during the process.
var installSteps = []string{
	"Setting up disk partitions...",
	"Formatting filesystems...",
	"Mounting partitions...",
	"Installing base system (pacstrap)...",
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
	}
}

func (m InstallModel) Init() tea.Cmd {
	return nil
}

// StartInstall returns the command to begin installation in a goroutine.
func (m InstallModel) StartInstall() tea.Cmd {
	m.config.InstallStarted = true
	return func() tea.Msg {
		// Run installation simulation in background
		// In production, this would call real installer.Install() via channel
		go func() {
			for i := 0; i < len(installSteps); i++ {
				time.Sleep(600 * time.Millisecond)
				// This is a simplified simulation.
				// In production, progressCh := make(chan installer.ProgressUpdate)
				// go installer.Install(progressCh) and read from channel
			}
		}()
		return installProgressMsg{
			Percent: 0,
			Message: "Starting installation...",
			Done:    false,
		}
	}
}

func (m InstallModel) Update(msg tea.Msg) (InstallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case installProgressMsg:
		m.currentStep++
		m.config.ProgressPercent = msg.Percent
		m.config.ProgressMessage = msg.Message

		if msg.Done {
			m.completed = true
			m.config.InstallComplete = true
			return m, nil
		}

		// Advance to next step after a short delay
		return m, func() tea.Msg {
			time.Sleep(600 * time.Millisecond)
			step := m.currentStep
			if step >= len(installSteps) {
				return installProgressMsg{
					Percent: 100,
					Message: "Installation complete!",
					Done:    true,
				}
			}
			return installProgressMsg{
				Percent: float64(step+1) / float64(len(installSteps)) * 100,
				Message: installSteps[step],
				Done:    step == len(installSteps)-1,
			}
		}

	case tea.KeyMsg:
		if m.completed {
			switch msg.String() {
			case "enter", " ", "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m InstallModel) View() string {
	if m.completed {
		return m.completedView()
	}
	return m.inProgressView()
}

func (m InstallModel) inProgressView() string {
	logo := MiniArchLogo()
	title := TitleStyle.Render("Installing Arch Linux Server...")

	// Spinner animation
	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinnerFrames[m.spinnerFrame]
	m.spinnerFrame = (m.spinnerFrame + 1) % len(spinnerFrames)

	spinnerView := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render(spinner + " " + m.config.ProgressMessage)

	// Progress bar
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
	bar = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render(bar)

	percentText := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render(fmt.Sprintf("%.0f%%", m.config.ProgressPercent))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		title,
		"",
		BoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				spinnerView,
				"",
				bar,
				"",
				percentText,
			),
		),
	)

	return Screen(TotalSteps, content, "Installation in progress... Press Ctrl+C to cancel (not recommended)")
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