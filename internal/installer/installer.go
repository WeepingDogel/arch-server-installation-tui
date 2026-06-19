package installer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/model"
)

// Installer handles the actual Arch Linux installation process.
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
	Percent   float64
	Message   string
	LogOutput string
	Done      bool
	Err       error
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

// partitionDisk creates partitions based on user config.
func (inst *Installer) partitionDisk() error {
	dev := inst.config.DiskDevice
	if dev == "" {
		return fmt.Errorf("no disk device selected")
	}

	if inst.config.PartitionScheme == "gpt" {
		_, err := safeExec("parted", "-s", dev, "mklabel", "gpt")
		if err != nil {
			return err
		}
		_, err = safeExec("parted", "-s", dev, "mkpart", "primary", "fat32", "1M", "513M")
		if err != nil {
			return err
		}
		_, err = safeExec("parted", "-s", dev, "set", "1", "esp", "on")
		if err != nil {
			return err
		}
		if inst.config.SwapSize != "" {
			efiEnd := unitToMB("513M")
			swapEnd := efiEnd + unitToMB(inst.config.SwapSize)
			_, err = safeExec("parted", "-s", dev, "mkpart", "primary", "linux-swap", fmt.Sprintf("%dM", efiEnd), fmt.Sprintf("%dM", swapEnd))
			if err != nil {
				return err
			}
			_, err = safeExec("parted", "-s", dev, "mkpart", "primary", "ext4", fmt.Sprintf("%dM", swapEnd), "100%")
			return err
		}
		_, err = safeExec("parted", "-s", dev, "mkpart", "primary", "ext4", "513M", "100%")
		return err
	}

	_, err := safeExec("parted", "-s", dev, "mklabel", "msdos")
	if err != nil {
		return err
	}
	_, err = safeExec("parted", "-s", dev, "mkpart", "primary", "ext4", "1M", "100%")
	if err != nil {
		return err
	}
	_, err = safeExec("parted", "-s", dev, "set", "1", "boot", "on")
	return err
}

// unitToMB parses a size string like "512M", "4G" and returns MB.
func unitToMB(size string) int {
	if strings.HasSuffix(size, "M") || strings.HasSuffix(size, "m") {
		var n int
		_, _ = fmt.Sscanf(size, "%d", &n)
		return n
	}
	if strings.HasSuffix(size, "G") || strings.HasSuffix(size, "g") {
		var n float64
		_, _ = fmt.Sscanf(size, "%f", &n)
		return int(n * 1024)
	}
	return 0
}

// formatFilesystems formats the created partitions.
func (inst *Installer) formatFilesystems() error {
	dev := inst.config.DiskDevice
	fsType := inst.config.FilesystemType

	if inst.config.PartitionScheme == "gpt" {
		_, err := safeExec("mkfs.fat", "-F32", devPart(dev, 1))
		if err != nil {
			return err
		}
		rootPart := 2
		if inst.config.SwapSize != "" {
			rootPart = 3
		}
		return formatPartition(devPart(dev, rootPart), fsType)
	}

	return formatPartition(devPart(dev, 1), fsType)
}

func devPart(dev string, part int) string {
	if strings.Contains(dev, "nvme") || strings.Contains(dev, "mmcblk") {
		return fmt.Sprintf("%sp%d", dev, part)
	}
	return fmt.Sprintf("%s%d", dev, part)
}

func formatPartition(part, fsType string) error {
	switch fsType {
	case "btrfs":
		_, err := safeExec("mkfs.btrfs", "-f", part)
		return err
	case "xfs":
		_, err := safeExec("mkfs.xfs", "-f", part)
		return err
	case "f2fs":
		_, err := safeExec("mkfs.f2fs", "-f", part)
		return err
	default:
		_, err := safeExec("mkfs.ext4", "-F", part)
		return err
	}
}

// mountPartitions mounts filesystems to /mnt.
func (inst *Installer) mountPartitions() error {
	dev := inst.config.DiskDevice
	rootPart := 2
	if inst.config.SwapSize != "" {
		rootPart = 3
	}
	if inst.config.PartitionScheme == "mbr" {
		rootPart = 1
	}

	_, err := safeExec("mount", devPart(dev, rootPart), "/mnt")
	if err != nil {
		return err
	}

	if inst.config.PartitionScheme == "gpt" {
		_, err = safeExec("mkdir", "-p", "/mnt/boot")
		if err != nil {
			return err
		}
		_, err = safeExec("mount", devPart(dev, 1), "/mnt/boot")
		if err != nil {
			return err
		}
	}

	if inst.config.SwapSize != "" {
		_, err = safeExec("swapon", devPart(dev, 2))
		if err != nil {
			return err
		}
	}

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

// configureLocale sets up locale configuration (supports multiple locales).
func (inst *Installer) configureLocale() error {
	locales := inst.config.Locales
	if len(locales) == 0 {
		locales = []string{"en_US.UTF-8"}
	}

	for _, locale := range locales {
		_, err := inst.chrootExec("sed", "-i", fmt.Sprintf("s/^#%s/%s/", locale, locale), "/etc/locale.gen")
		if err != nil {
			return err
		}
	}

	_, err := inst.chrootExec("locale-gen")
	if err != nil {
		return err
	}
	echoCmd := fmt.Sprintf("echo 'LANG=%s' > /etc/locale.conf", locales[0])
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
	_, err := inst.chrootExec("systemctl", "enable", "systemd-networkd")
	if err != nil {
		return err
	}
	_, err = inst.chrootExec("systemctl", "enable", "systemd-resolved")
	return err
}

// setRootPassword sets the root password.
func (inst *Installer) setRootPassword() error {
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
	_, err := inst.chrootExec("systemctl", "enable", "sshd")
	if err != nil {
		return err
	}
	if inst.config.SSHPort != 22 {
		portStr := fmt.Sprintf("Port %d", inst.config.SSHPort)
		_, err = inst.chrootExec("sh", "-c",
			fmt.Sprintf("echo '%s' >> /etc/ssh/sshd_config", portStr))
		if err != nil {
			return err
		}
	}
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
	if inst.config.EnableArchCN {
		packages = append(packages, "archlinuxcn-keyring")
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
	_, err := inst.chrootExec("systemctl", "enable", "systemd-timesyncd")
	if err != nil {
		return err
	}
	if inst.config.EnableArchCN && inst.config.ArchCNMirror != "" {
		repoLine := fmt.Sprintf("\n[archlinuxcn]\nServer = %s/$arch\n", inst.config.ArchCNMirror)
		_, err = inst.chrootExec("sh", "-c",
			fmt.Sprintf("echo '%s' >> /etc/pacman.conf", repoLine))
	}
	return err
}
