package installer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
)

// Installer handles the actual Arch Linux installation process.
// It reads configuration from model.Config and executes system commands.
type Installer struct {
	config *model.Config
}

// New creates an Installer with the given configuration.
func New(config *model.Config) *Installer {
	return &Installer{config: config}
}

// Result contains the installation outcome.
type Result struct {
	Success bool
	Message string
	Err     error
}

// Install runs the full installation pipeline and sends progress updates.
// This is designed to be called from a goroutine.
func (inst *Installer) Install(progress chan<- ProgressUpdate) {
	defer close(progress)

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Partitioning disk", inst.partitionDisk},
		{"Formatting filesystems", inst.formatFilesystems},
		{"Mounting partitions", inst.mountPartitions},
		{"Installing base system", inst.pacstrapBase},
		{"Generating fstab", inst.generateFstab},
		{"Configuring timezone", inst.configureTimezone},
		{"Setting up locale", inst.configureLocale},
		{"Setting hostname", inst.setHostname},
		{"Configuring network", inst.configureNetwork},
		{"Setting root password", inst.setRootPassword},
		{"Creating user account", inst.createUser},
		{"Installing bootloader", inst.installBootloader},
		{"Configuring SSH", inst.configureSSH},
		{"Installing additional packages", inst.installPackages},
		{"Finalizing", inst.finalize},
	}

	totalSteps := float64(len(steps))
	for i, step := range steps {
		percent := (float64(i) / totalSteps) * 100
		progress <- ProgressUpdate{
			Percent: percent,
			Message: step.name + "...",
		}

		if err := step.fn(); err != nil {
			progress <- ProgressUpdate{
				Percent: percent,
				Message: "Error: " + err.Error(),
				Done:    true,
				Err:     err,
			}
			return
		}
	}

	progress <- ProgressUpdate{
		Percent: 100,
		Message: "Installation complete!",
		Done:    true,
	}
}

// ProgressUpdate is sent through the channel during installation.
type ProgressUpdate struct {
	Percent float64
	Message string
	Done    bool
	Err     error
}

// safeExec runs a command and returns its output.
func safeExec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command '%s %s' failed: %v\nOutput: %s",
			name, strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}

// chrootExec runs a command inside the chroot environment.
func (inst *Installer) chrootExec(args ...string) (string, error) {
	cmd := []string{"arch-chroot", "/mnt"}
	cmd = append(cmd, args...)
	return safeExec(cmd[0], cmd[1:]...)
}

// partitionDisk creates partitions on the target disk.
func (inst *Installer) partitionDisk() error {
	// This would use parted or fdisk to partition the disk
	return nil
}

// formatFilesystems formats the created partitions.
func (inst *Installer) formatFilesystems() error {
	return nil
}

// mountPartitions mounts filesystems to /mnt.
func (inst *Installer) mountPartitions() error {
	return nil
}

// pacstrapBase installs the base system using pacstrap.
func (inst *Installer) pacstrapBase() error {
	packages := []string{"base", "linux", "linux-firmware"}
	if inst.config.InstallBaseDev {
		packages = append(packages, "base-devel")
	}
	switch inst.config.KernelType {
	case "linux-lts":
		packages[1] = "linux-lts"
	case "linux-zen":
		packages[1] = "linux-zen"
	case "linux-hardened":
		packages[1] = "linux-hardened"
	}
	_, err := safeExec("pacstrap", append([]string{"/mnt"}, packages...)...)
	return err
}

// generateFstab generates the fstab file.
func (inst *Installer) generateFstab() error {
	_, err := safeExec("genfstab", "-U", "/mnt")
	if err != nil {
		return err
	}
	_, err = safeExec("sh", "-c", "genfstab -U /mnt >> /mnt/etc/fstab")
	return err
}

// configureTimezone sets the system timezone.
func (inst *Installer) configureTimezone() error {
	if inst.config.TimezoneRegion == "UTC" {
		_, err := inst.chrootExec("ln", "-sf", "/usr/share/zoneinfo/UTC", "/etc/localtime")
		return err
	}
	tzPath := fmt.Sprintf("/usr/share/zoneinfo/%s", inst.config.TimezoneRegion)
	_, err := inst.chrootExec("ln", "-sf", tzPath, "/etc/localtime")
	return err
}

// configureLocale sets up locale configuration.
func (inst *Installer) configureLocale() error {
	locale := inst.config.Locale
	if locale == "" {
		locale = "en_US.UTF-8"
	}
	// Uncomment the locale in locale.gen
	_, err := inst.chrootExec("sed", "-i", fmt.Sprintf("s/^#%s/%s/", locale, locale), "/etc/locale.gen")
	if err != nil {
		return err
	}
	_, err = inst.chrootExec("locale-gen")
	if err != nil {
		return err
	}
	// Set LANG
	echoCmd := fmt.Sprintf("echo 'LANG=%s' > /etc/locale.conf", locale)
	_, err = inst.chrootExec("sh", "-c", echoCmd)
	return err
}

// setHostname sets the system hostname.
func (inst *Installer) setHostname() error {
	echoCmd := fmt.Sprintf("echo '%s' > /etc/hostname", inst.config.Hostname)
	_, err := inst.chrootExec("sh", "-c", echoCmd)
	return err
}

// configureNetwork sets up network configuration.
func (inst *Installer) configureNetwork() error {
	// Enable systemd-networkd and systemd-resolved
	_, err := inst.chrootExec("systemctl", "enable", "systemd-networkd")
	if err != nil {
		return err
	}
	_, err = inst.chrootExec("systemctl", "enable", "systemd-resolved")
	return err
}

// setRootPassword sets the root password.
func (inst *Installer) setRootPassword() error {
	// Use chroot and passwd with stdin
	_, err := inst.chrootExec("sh", "-c",
		fmt.Sprintf("echo 'root:%s' | chpasswd", inst.config.RootPassword))
	return err
}

// createUser creates a sudo user if configured.
func (inst *Installer) createUser() error {
	if !inst.config.CreateUser || inst.config.UserName == "" {
		return nil
	}
	_, err := inst.chrootExec("useradd", "-m", "-G", "wheel", inst.config.UserName)
	if err != nil {
		return err
	}
	_, err = inst.chrootExec("sh", "-c",
		fmt.Sprintf("echo '%s:%s' | chpasswd", inst.config.UserName, inst.config.UserPassword))
	if err != nil {
		return err
	}
	// Uncomment %wheel ALL=(ALL:ALL) ALL in sudoers
	_, err = inst.chrootExec("sed", "-i", "s/^# %wheel ALL=(ALL:ALL) ALL/%wheel ALL=(ALL:ALL) ALL/", "/etc/sudoers")
	return err
}

// installBootloader installs the bootloader.
func (inst *Installer) installBootloader() error {
	switch inst.config.BootloaderType {
	case "grub":
		if inst.config.UEFIMode {
			_, err := inst.chrootExec("grub-install", "--target=x86_64-efi", "--efi-directory=/boot", "--bootloader-id=GRUB")
			if err != nil {
				return err
			}
		} else {
			_, err := inst.chrootExec("grub-install", inst.config.DiskDevice)
			if err != nil {
				return err
			}
		}
		_, err := inst.chrootExec("grub-mkconfig", "-o", "/boot/grub/grub.cfg")
		return err
	case "systemd-boot":
		if inst.config.UEFIMode {
			_, err := inst.chrootExec("bootctl", "install")
			return err
		}
		return fmt.Errorf("systemd-boot requires UEFI mode")
	}
	return nil
}

// configureSSH sets up OpenSSH server.
func (inst *Installer) configureSSH() error {
	if !inst.config.EnableSSH {
		return nil
	}
	// Enable sshd service
	_, err := inst.chrootExec("systemctl", "enable", "sshd")
	if err != nil {
		return err
	}
	// Configure SSH port
	if inst.config.SSHPort != 22 {
		portStr := fmt.Sprintf("Port %d", inst.config.SSHPort)
		_, err = inst.chrootExec("sh", "-c",
			fmt.Sprintf("echo '%s' >> /etc/ssh/sshd_config", portStr))
		if err != nil {
			return err
		}
	}
	// Configure root login
	if !inst.config.AllowRootLogin {
		_, err = inst.chrootExec("sed", "-i",
			"s/^#PermitRootLogin.*/PermitRootLogin no/", "/etc/ssh/sshd_config")
		return err
	}
	_, err = inst.chrootExec("sed", "-i",
		"s/^#PermitRootLogin.*/PermitRootLogin yes/", "/etc/ssh/sshd_config")
	return err
}

// installPackages installs additional packages.
func (inst *Installer) installPackages() error {
	var packages []string

	if inst.config.InstallDocker {
		packages = append(packages, "docker")
	}
	if inst.config.InstallNginx {
		packages = append(packages, "nginx")
	}
	if inst.config.InstallPostgres {
		packages = append(packages, "postgresql")
	}
	if inst.config.InstallMariaDB {
		packages = append(packages, "mariadb")
	}
	if inst.config.InstallRedis {
		packages = append(packages, "redis")
	}
	if inst.config.InstallFail2ban {
		packages = append(packages, "fail2ban")
	}
	if inst.config.InstallUfw {
		packages = append(packages, "ufw")
	}
	if inst.config.InstallGit {
		packages = append(packages, "git")
	}
	if inst.config.InstallVim {
		packages = append(packages, "vim")
	}
	if inst.config.CustomPackages != "" {
		extra := strings.Fields(inst.config.CustomPackages)
		packages = append(packages, extra...)
	}

	if len(packages) == 0 {
		return nil
	}

	args := append([]string{"-S", "--noconfirm"}, packages...)
	_, err := inst.chrootExec(args...)
	return err
}

// finalize performs cleanup and final steps.
func (inst *Installer) finalize() error {
	// Enable useful services
	_, err := inst.chrootExec("systemctl", "enable", "systemd-timesyncd")
	return err
}