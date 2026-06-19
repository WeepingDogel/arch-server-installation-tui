package installer

import (
	"bufio"
	"fmt"
	"io"
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

// ProgressUpdate is sent through the channel during installation.
type ProgressUpdate struct {
	Percent   float64
	Message   string
	LogOutput string
	Stage     string // step name for grouping
	Done      bool
	Err       error
}

// Install runs the full installation pipeline.
func (inst *Installer) Install(progress chan<- ProgressUpdate) {
	defer close(progress)

	totalSteps := 15.0
	step := 0

	// Helper to run a step with real-time log output
	run := func(name string, fn func(chan<- string) error) {
		percent := (float64(step) / totalSteps) * 100
		step++

		// Send initial message
		progress <- ProgressUpdate{
			Percent: percent,
			Message: name + "...",
			Stage:   name,
		}

		// Create a channel for log lines
		logCh := make(chan string, 100)
		done := make(chan error, 1)

		go func() {
			done <- fn(logCh)
			close(logCh)
		}()

		// Stream log lines until step finishes
		for line := range logCh {
			progress <- ProgressUpdate{
				Percent:   percent,
				Message:   name,
				LogOutput: line,
				Stage:     name,
			}
		}

		err := <-done
		if err != nil {
			progress <- ProgressUpdate{
				Percent:   percent,
				Message:   "Error: " + err.Error(),
				LogOutput: err.Error(),
				Stage:     name + " [FAILED]",
				Done:      true,
				Err:       err,
			}
			return
		}

		progress <- ProgressUpdate{
			Percent:   percent + (100.0/totalSteps)*0.5,
			Message:   name + " ✓",
			Stage:     name + " [OK]",
			LogOutput: "",
		}
	}

	run("Partitioning disk", inst.partitionDisk)
	run("Formatting filesystems", inst.formatFilesystems)
	run("Mounting partitions", inst.mountPartitions)
	run("Installing base system", inst.pacstrapBase)
	run("Generating fstab", inst.generateFstab)
	run("Configuring timezone", inst.configureTimezone)
	run("Setting up locale", inst.configureLocale)
	run("Setting hostname", inst.setHostname)
	run("Configuring network", inst.configureNetwork)
	run("Setting root password", inst.setRootPassword)
	run("Creating user account", inst.createUser)
	run("Installing bootloader", inst.installBootloader)
	run("Configuring SSH", inst.configureSSH)
	run("Installing additional packages", inst.installPackages)
	run("Finalizing", inst.finalize)

	progress <- ProgressUpdate{
		Percent: 100,
		Message: "Installation complete!",
		Done:    true,
	}
}

// streamExec runs a command and streams each line of output to the channel.
func streamExec(logCh chan<- string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Read stdout line by line
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 1024*64), 1024*64)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				logCh <- line
			}
		}
	}()

	// Read stderr line by line
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				logCh <- line
			}
		}
	}()

	// Wait for stderr/stdout to finish
	_, _ = io.Copy(io.Discard, stderr)
	_, _ = io.Copy(io.Discard, stdout)

	return cmd.Wait()
}

// chrootExecStream runs a command inside chroot and streams output.
func (inst *Installer) chrootExecStream(logCh chan<- string, args ...string) error {
	cmd := []string{"arch-chroot", "/mnt"}
	cmd = append(cmd, args...)
	return streamExec(logCh, cmd[0], cmd[1:]...)
}

// partitionDisk creates partitions based on user config.
func (inst *Installer) partitionDisk(logCh chan<- string) error {
	dev := inst.config.DiskDevice
	if dev == "" {
		return fmt.Errorf("no disk device selected")
	}

	logCh <- fmt.Sprintf("Partitioning %s (%s scheme)...", dev, strings.ToUpper(inst.config.PartitionScheme))

	if inst.config.PartitionScheme == "gpt" {
		if err := streamExec(logCh, "parted", "-s", dev, "mklabel", "gpt"); err != nil {
			return fmt.Errorf("failed to create GPT label: %w", err)
		}
		logCh <- "Created GPT partition table."

		if err := streamExec(logCh, "parted", "-s", dev, "mkpart", "primary", "fat32", "1M", "513M"); err != nil {
			return fmt.Errorf("failed to create EFI partition: %w", err)
		}
		logCh <- "Created EFI partition (512MB)."

		if inst.config.UEFIMode {
			if err := streamExec(logCh, "parted", "-s", dev, "set", "1", "esp", "on"); err != nil {
				return fmt.Errorf("failed to set ESP flag: %w", err)
			}
			logCh <- "Set ESP flag on partition 1."
		}

		if inst.config.SwapSize != "" {
			swapEnd := 513 + unitToMB(inst.config.SwapSize)
			if err := streamExec(logCh, "parted", "-s", dev, "mkpart", "primary", "linux-swap", "513M", fmt.Sprintf("%dM", swapEnd)); err != nil {
				return fmt.Errorf("failed to create swap partition: %w", err)
			}
			logCh <- fmt.Sprintf("Created swap partition (%s).", inst.config.SwapSize)

			if err := streamExec(logCh, "parted", "-s", dev, "mkpart", "primary", "ext4", fmt.Sprintf("%dM", swapEnd), "100%"); err != nil {
				return fmt.Errorf("failed to create root partition: %w", err)
			}
			logCh <- "Created root partition."
			return nil
		}

		if err := streamExec(logCh, "parted", "-s", dev, "mkpart", "primary", "ext4", "513M", "100%"); err != nil {
			return fmt.Errorf("failed to create root partition: %w", err)
		}
		logCh <- "Created root partition."
		return nil
	}

	// MBR
	if err := streamExec(logCh, "parted", "-s", dev, "mklabel", "msdos"); err != nil {
		return fmt.Errorf("failed to create MBR label: %w", err)
	}
	logCh <- "Created MBR partition table."

	if err := streamExec(logCh, "parted", "-s", dev, "mkpart", "primary", "ext4", "1M", "100%"); err != nil {
		return fmt.Errorf("failed to create root partition: %w", err)
	}
	logCh <- "Created root partition."

	if err := streamExec(logCh, "parted", "-s", dev, "set", "1", "boot", "on"); err != nil {
		return fmt.Errorf("failed to set boot flag: %w", err)
	}
	logCh <- "Set boot flag on partition 1."
	return nil
}

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

func devPart(dev string, part int) string {
	if strings.Contains(dev, "nvme") || strings.Contains(dev, "mmcblk") {
		return fmt.Sprintf("%sp%d", dev, part)
	}
	return fmt.Sprintf("%s%d", dev, part)
}

// formatFilesystems formats the created partitions.
func (inst *Installer) formatFilesystems(logCh chan<- string) error {
	dev := inst.config.DiskDevice
	fsType := inst.config.FilesystemType

	switch inst.config.PartitionScheme {
	case "gpt":
		logCh <- fmt.Sprintf("Formatting EFI partition %s...", devPart(dev, 1))
		if err := streamExec(logCh, "mkfs.fat", "-F32", devPart(dev, 1)); err != nil {
			return fmt.Errorf("failed to format EFI: %w", err)
		}
		rootPart := 2
		if inst.config.SwapSize != "" {
			rootPart = 3
		}
		logCh <- fmt.Sprintf("Formatting root partition %s as %s...", devPart(dev, rootPart), fsType)
		return formatPartitionStream(logCh, devPart(dev, rootPart), fsType)

	default: // MBR
		return formatPartitionStream(logCh, devPart(dev, 1), fsType)
	}
}

func formatPartitionStream(logCh chan<- string, part, fsType string) error {
	switch fsType {
	case "btrfs":
		logCh <- fmt.Sprintf("Running mkfs.btrfs on %s...", part)
		return streamExec(logCh, "mkfs.btrfs", "-f", part)
	case "xfs":
		logCh <- fmt.Sprintf("Running mkfs.xfs on %s...", part)
		return streamExec(logCh, "mkfs.xfs", "-f", part)
	case "f2fs":
		logCh <- fmt.Sprintf("Running mkfs.f2fs on %s...", part)
		return streamExec(logCh, "mkfs.f2fs", "-f", part)
	default:
		logCh <- fmt.Sprintf("Running mkfs.ext4 on %s...", part)
		return streamExec(logCh, "mkfs.ext4", "-F", part)
	}
}

// mountPartitions mounts filesystems to /mnt.
func (inst *Installer) mountPartitions(logCh chan<- string) error {
	dev := inst.config.DiskDevice
	rootPart := 2
	if inst.config.SwapSize != "" {
		rootPart = 3
	}
	if inst.config.PartitionScheme == "mbr" {
		rootPart = 1
	}

	logCh <- fmt.Sprintf("Mounting %s to /mnt...", devPart(dev, rootPart))
	if err := streamExec(logCh, "mount", devPart(dev, rootPart), "/mnt"); err != nil {
		return fmt.Errorf("failed to mount root: %w", err)
	}

	if inst.config.PartitionScheme == "gpt" {
		logCh <- "Creating /mnt/boot..."
		if err := streamExec(logCh, "mkdir", "-p", "/mnt/boot"); err != nil {
			return err
		}
		logCh <- fmt.Sprintf("Mounting %s to /mnt/boot...", devPart(dev, 1))
		if err := streamExec(logCh, "mount", devPart(dev, 1), "/mnt/boot"); err != nil {
			return fmt.Errorf("failed to mount EFI: %w", err)
		}
	}

	if inst.config.SwapSize != "" {
		logCh <- "Enabling swap..."
		if err := streamExec(logCh, "swapon", devPart(dev, 2)); err != nil {
			return fmt.Errorf("failed to enable swap: %w", err)
		}
	}

	return nil
}

// pacstrapBase installs the base system using pacstrap.
func (inst *Installer) pacstrapBase(logCh chan<- string) error {
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
	// Bootloader packages
	if inst.config.BootloaderType == "grub" {
		packages = append(packages, "grub")
		if inst.config.UEFIMode {
			packages = append(packages, "efibootmgr")
		}
	}

	logCh <- fmt.Sprintf("Installing base system via pacstrap (%d packages)...", len(packages))
	logCh <- "Packages: " + strings.Join(packages, ", ")
	logCh <- "This may take a while depending on your internet speed..."

	args := append([]string{"/mnt"}, packages...)
	if err := streamExec(logCh, "pacstrap", args...); err != nil {
		return fmt.Errorf("pacstrap failed: %w", err)
	}
	logCh <- "Base system installed successfully."
	return nil
}

// generateFstab generates the fstab file.
func (inst *Installer) generateFstab(logCh chan<- string) error {
	logCh <- "Generating fstab..."
	return streamExec(logCh, "sh", "-c", "genfstab -U /mnt >> /mnt/etc/fstab")
}

// configureTimezone sets the system timezone.
func (inst *Installer) configureTimezone(logCh chan<- string) error {
	if inst.config.TimezoneRegion == "UTC" {
		logCh <- "Setting timezone to UTC..."
		return inst.chrootExecStream(logCh, "ln", "-sf", "/usr/share/zoneinfo/UTC", "/etc/localtime")
	}
	tzPath := fmt.Sprintf("/usr/share/zoneinfo/%s", inst.config.TimezoneRegion)
	logCh <- fmt.Sprintf("Setting timezone to %s...", inst.config.TimezoneRegion)
	return inst.chrootExecStream(logCh, "ln", "-sf", tzPath, "/etc/localtime")
}

// configureLocale sets up locale configuration.
func (inst *Installer) configureLocale(logCh chan<- string) error {
	locales := inst.config.Locales
	if len(locales) == 0 {
		locales = []string{"en_US.UTF-8"}
	}
	for _, locale := range locales {
		logCh <- fmt.Sprintf("Enabling locale: %s", locale)
		if err := inst.chrootExecStream(logCh, "sed", "-i", fmt.Sprintf("s/^#%s/%s/", locale, locale), "/etc/locale.gen"); err != nil {
			return fmt.Errorf("failed to enable locale %s: %w", locale, err)
		}
	}
	logCh <- "Running locale-gen..."
	if err := inst.chrootExecStream(logCh, "locale-gen"); err != nil {
		return fmt.Errorf("locale-gen failed: %w", err)
	}
	logCh <- fmt.Sprintf("Setting LANG=%s...", locales[0])
	return inst.chrootExecStream(logCh, "sh", "-c", fmt.Sprintf("echo 'LANG=%s' > /etc/locale.conf", locales[0]))
}

// setHostname sets the system hostname.
func (inst *Installer) setHostname(logCh chan<- string) error {
	logCh <- fmt.Sprintf("Setting hostname to %s...", inst.config.Hostname)
	return inst.chrootExecStream(logCh, "sh", "-c", fmt.Sprintf("echo '%s' > /etc/hostname", inst.config.Hostname))
}

// configureNetwork sets up network configuration.
func (inst *Installer) configureNetwork(logCh chan<- string) error {
	logCh <- "Enabling systemd-networkd..."
	if err := inst.chrootExecStream(logCh, "systemctl", "enable", "systemd-networkd"); err != nil {
		return err
	}
	logCh <- "Enabling systemd-resolved..."
	return inst.chrootExecStream(logCh, "systemctl", "enable", "systemd-resolved")
}

// setRootPassword sets the root password.
func (inst *Installer) setRootPassword(logCh chan<- string) error {
	logCh <- "Setting root password..."
	return inst.chrootExecStream(logCh, "sh", "-c",
		fmt.Sprintf("echo 'root:%s' | chpasswd", inst.config.RootPassword))
}

// createUser creates a sudo user if configured.
func (inst *Installer) createUser(logCh chan<- string) error {
	if !inst.config.CreateUser || inst.config.UserName == "" {
		logCh <- "Skipping user creation."
		return nil
	}
	logCh <- fmt.Sprintf("Creating user %s...", inst.config.UserName)
	if err := inst.chrootExecStream(logCh, "useradd", "-m", "-G", "wheel", inst.config.UserName); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	if err := inst.chrootExecStream(logCh, "sh", "-c",
		fmt.Sprintf("echo '%s:%s' | chpasswd", inst.config.UserName, inst.config.UserPassword)); err != nil {
		return fmt.Errorf("failed to set user password: %w", err)
	}
	logCh <- "Granting sudo access..."
	return inst.chrootExecStream(logCh, "sed", "-i", "s/^# %wheel ALL=(ALL:ALL) ALL/%wheel ALL=(ALL:ALL) ALL/", "/etc/sudoers")
}

// installBootloader installs the bootloader.
func (inst *Installer) installBootloader(logCh chan<- string) error {
	switch inst.config.BootloaderType {
	case "grub":
		logCh <- "Installing GRUB bootloader..."
		if inst.config.UEFIMode {
			logCh <- "UEFI mode detected, installing GRUB for x86_64-efi..."
			// Mount efivarfs (optional, may fail on non-UEFI)
			_ = streamExec(logCh, "mount", "-t", "efivarfs", "efivarfs", "/sys/firmware/efi/efivars")

			// Install GRUB to the ESP mounted at /boot inside chroot
			if err := inst.chrootExecStream(logCh, "grub-install", "--target=x86_64-efi", "--efi-directory=/boot", "--bootloader-id=GRUB"); err != nil {
				return fmt.Errorf("grub-install failed: %w", err)
			}
		} else {
			logCh <- "BIOS mode detected, installing GRUB to MBR..."
			if err := inst.chrootExecStream(logCh, "grub-install", inst.config.DiskDevice); err != nil {
				return fmt.Errorf("grub-install failed: %w", err)
			}
		}
		logCh <- "Generating GRUB configuration..."
		if err := inst.chrootExecStream(logCh, "grub-mkconfig", "-o", "/boot/grub/grub.cfg"); err != nil {
			return fmt.Errorf("grub-mkconfig failed: %w", err)
		}
		logCh <- "GRUB installed successfully."

	case "systemd-boot":
		if !inst.config.UEFIMode {
			return fmt.Errorf("systemd-boot requires UEFI mode")
		}
		logCh <- "Installing systemd-boot..."
		if err := inst.chrootExecStream(logCh, "bootctl", "install"); err != nil {
			return fmt.Errorf("bootctl install failed: %w", err)
		}
		logCh <- "systemd-boot installed successfully."
	}
	return nil
}

// configureSSH sets up OpenSSH server.
func (inst *Installer) configureSSH(logCh chan<- string) error {
	if !inst.config.EnableSSH {
		logCh <- "SSH is disabled, skipping."
		return nil
	}
	logCh <- "Enabling sshd service..."
	if err := inst.chrootExecStream(logCh, "systemctl", "enable", "sshd"); err != nil {
		return err
	}
	if inst.config.SSHPort != 22 {
		logCh <- fmt.Sprintf("Setting SSH port to %d...", inst.config.SSHPort)
		if err := inst.chrootExecStream(logCh, "sh", "-c",
			fmt.Sprintf("echo 'Port %d' >> /etc/ssh/sshd_config", inst.config.SSHPort)); err != nil {
			return err
		}
	}
	rootLogin := "yes"
	if !inst.config.AllowRootLogin {
		rootLogin = "no"
	}
	logCh <- fmt.Sprintf("Setting PermitRootLogin to %s...", rootLogin)
	return inst.chrootExecStream(logCh, "sed", "-i",
		fmt.Sprintf("s/^#PermitRootLogin.*/PermitRootLogin %s/", rootLogin), "/etc/ssh/sshd_config")
}

// installPackages installs additional packages.
func (inst *Installer) installPackages(logCh chan<- string) error {
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
		logCh <- "No additional packages selected."
		return nil
	}

	logCh <- fmt.Sprintf("Installing %d additional package(s)...", len(packages))
	logCh <- "Packages: " + strings.Join(packages, ", ")
	logCh <- "This may take a while..."

	args := append([]string{"-S", "--noconfirm"}, packages...)
	if err := inst.chrootExecStream(logCh, args...); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}
	logCh <- "Additional packages installed."
	return nil
}

// finalize performs cleanup and final steps.
func (inst *Installer) finalize(logCh chan<- string) error {
	logCh <- "Enabling systemd-timesyncd..."
	if err := inst.chrootExecStream(logCh, "systemctl", "enable", "systemd-timesyncd"); err != nil {
		return err
	}

	if inst.config.EnableArchCN && inst.config.ArchCNMirror != "" {
		logCh <- fmt.Sprintf("Adding Arch Linux CN repository from %s...", inst.config.ArchCNMirror)
		repoLine := fmt.Sprintf("\n[archlinuxcn]\nServer = %s/$arch\n", inst.config.ArchCNMirror)
		if err := inst.chrootExecStream(logCh, "sh", "-c",
			fmt.Sprintf("echo '%s' >> /etc/pacman.conf", repoLine)); err != nil {
			return err
		}
	}

	logCh <- "Syncing disks..."
	if err := streamExec(logCh, "sync"); err != nil {
		return err
	}
	logCh <- "Installation complete! You can now reboot."
	return nil
}
