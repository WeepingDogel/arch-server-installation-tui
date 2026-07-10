package tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/installer"
	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InstallModel shows the installation progress with live streaming logs.
type InstallModel struct {
	config       *model.Config
	spinnerFrame int
	completed    bool
	installer    *installer.Installer
	logs         []string
	currentStage string
	progressCh   chan installer.ProgressUpdate
}

// NewInstallModel creates the installation progress screen.
func NewInstallModel(config *model.Config) InstallModel {
	return InstallModel{
		config:       config,
		spinnerFrame: 0,
		installer:    installer.New(config),
		logs:         make([]string, 0, 500),
	}
}

func (m InstallModel) Init() tea.Cmd { return nil }

// StartInstall begins installation and streams progress to the TUI.
func (m InstallModel) StartInstall() tea.Cmd {
	m.config.InstallStarted = true
	m.logs = append(m.logs, "◆ Installation started")
	m.progressCh = make(chan installer.ProgressUpdate)
	go m.installer.Install(m.progressCh)
	return m.pollNext()
}

// pollNext reads one message from the progress channel.
func (m InstallModel) pollNext() tea.Cmd {
	return func() tea.Msg {
		p, ok := <-m.progressCh
		if !ok {
			// Channel closed - installation finished
			return installProgressMsg{
				Percent: 100,
				Message: "Installation complete!",
				Done:    true,
			}
		}
		return installProgressMsg{
			Percent:   p.Percent,
			Message:   p.Message,
			LogOutput: p.LogOutput,
			Stage:     p.Stage,
			Done:      p.Done,
			Err:       p.Err,
		}
	}
}

func (m InstallModel) Update(msg tea.Msg) (InstallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case installProgressMsg:
		m.spinnerFrame = (m.spinnerFrame + 1) % 10
		m.config.ProgressPercent = msg.Percent
		m.config.ProgressMessage = msg.Message

		if msg.Stage != "" {
			m.currentStage = msg.Stage
		}
		if msg.LogOutput != "" {
			line := msg.LogOutput
			if len(line) > 200 {
				line = line[:200] + "..."
			}
			m.logs = append(m.logs, line)
			if len(m.logs) > 500 {
				m.logs = m.logs[len(m.logs)-500:]
			}
		}

		if msg.Done {
			m.completed = true
			m.config.InstallComplete = true
			m.logs = append(m.logs, "")
			m.logs = append(m.logs, "✓ Installation complete!")
			m.logs = append(m.logs, "")
			m.logs = append(m.logs, "Remove the installation media and press ENTER to reboot.")
			return m, nil
		}

		// Schedule next poll
		return m, func() tea.Msg {
			time.Sleep(50 * time.Millisecond)
			return pollMsg{}
		}

	case pollMsg:
		return m, m.pollNext()

	case tea.KeyMsg:
		if m.completed && (msg.String() == "enter" || msg.String() == "r") {
			return m, m.rebootCmd()
		}
		if m.completed && (msg.String() == "q" || msg.String() == "ctrl+c") {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m InstallModel) rebootCmd() tea.Cmd {
	return func() tea.Msg {
		// Sync disks to flush all pending writes
		_ = exec.Command("sync").Run()
		// Unmount /mnt recursively to clean up the installed system
		_ = exec.Command("umount", "-R", "/mnt").Run()
		// Reboot into the newly installed system
		_ = exec.Command("reboot").Start()
		return nil
	}
}

type pollMsg struct{}

func (m InstallModel) View() string {
	if m.completed {
		return m.completedView()
	}
	return m.inProgressView()
}

func (m InstallModel) inProgressView() string {
	// Progress bar at top
	stepIndicator := StepIndicator(min(int(m.config.ProgressPercent/100.0*13)+1, 13), TotalSteps)

	barWidth := 40
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

	// Current stage display
	stageText := lipgloss.NewStyle().Foreground(ColorWhite).Render("▶ " + m.config.ProgressMessage)

	progressBox := lipgloss.JoinVertical(
		lipgloss.Left,
		stepIndicator,
		"",
		bar,
		percentText,
		"",
		stageText,
	)

	// Live log view (last 15 lines)
	var logContent string
	start := 0
	if len(m.logs) > 15 {
		start = len(m.logs) - 15
	}
	for i, line := range m.logs[start:] {
		num := start + i + 1
		// Colorize based on content
		style := lipgloss.NewStyle().Foreground(ColorGray)
		if strings.Contains(line, "Error") || strings.Contains(line, "FAILED") {
			style = ErrorStyle
		} else if strings.Contains(line, "✓") || strings.Contains(line, "done") {
			style = SuccessStyle
		} else if strings.Contains(line, "◆") || strings.Contains(line, "▶") {
			style = lipgloss.NewStyle().Foreground(ColorPrimary)
		}
		logContent += style.Render(fmt.Sprintf("%3d ", num)) + line + "\n"
	}

	if m.currentStage != "" {
		logContent = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("["+m.currentStage+"]\n") + logContent
	}

	logBox := BoxStyle.MaxWidth(60).Render(logContent)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		"",
		progressBox,
		"",
		logBox,
	)
}

func (m InstallModel) completedView() string {
	logo := ArchLogo()
	title := lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Render("✓ Installation Complete!")
	subtitle := SubtitleStyle.Render("Arch Linux Server is ready.")

	summaryInfo := BoxStyle.MaxWidth(56).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf("Hostname: %s", m.config.Hostname),
			fmt.Sprintf("Disk: %s (%s)", m.config.DiskDevice, m.config.DiskSize),
			fmt.Sprintf("Boot: %s / %s", m.config.PartitionScheme, m.config.BootloaderType),
			fmt.Sprintf("Users: root + %s", ifElse(m.config.CreateUser, m.config.UserName, "(none)")),
			fmt.Sprintf("SSH: Port %d", m.config.SSHPort),
			"",
			InfoBox("Remove the installation media before rebooting!"),
		),
	)

	// Reboot button
	rebootBtn := lipgloss.NewStyle().
		Background(ColorSuccess).
		Foreground(ColorWhite).
		Bold(true).
		Padding(0, 6).
		Render("⟳  Reboot Now")

	quitBtn := lipgloss.NewStyle().
		Background(ColorGray).
		Foreground(ColorWhite).
		Padding(0, 6).
		Render("Quit")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, rebootBtn, "  ", quitBtn)
	help := HelpStyle.Render("ENTER=Reboot  Q=Quit  Ctrl+C=Quit")

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
		buttons,
		"",
		help,
	)
}
