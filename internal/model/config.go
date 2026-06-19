package model

// Config holds all user choices throughout the installation wizard.
// It is shared by reference across all TUI screens.
type Config struct {
	// Keyboard
	KeyboardLayout string

	// Network
	NetworkDHCP  bool
	Hostname     string
	IPAddress    string
	Netmask      string
	Gateway      string
	DNSServers   string
	NetworkIface string // detected interface name

	// Mirror
	MirrorURL     string
	MirrorCountry string
	CustomMirror  string // Custom mirror URL input
	EnableArchCN  bool   // Toggle Arch Linux CN repository
	ArchCNMirror  string // Arch Linux CN mirror URL

	// Disk
	DiskDevice          string
	DiskSize            string // detected size
	UseAutoPartitioning bool
	RootPartitionSize   string // e.g., "20G" or "100%"
	EncryptDisk         bool
	LVMEnabled          bool

	// Filesystem
	FilesystemType string // ext4, btrfs, xfs, f2fs

	// Bootloader
	BootloaderType string // grub, systemd-boot
	UEFIMode       bool

	// Timezone & Locale
	TimezoneRegion string
	TimezoneCity   string
	Locales        []string // multiple locales supported

	// Users
	RootPassword string
	UserName     string
	UserPassword string
	CreateUser   bool

	// SSH
	EnableSSH         bool
	SSHPort           int
	AllowRootLogin    bool
	ImportSSHKeys     bool
	SSHAuthorizedKeys string

	// Packages
	KernelType      string // linux, linux-lts, linux-zen, linux-hardened
	InstallBaseDev  bool
	InstallDocker   bool
	InstallNginx    bool
	InstallPostgres bool
	InstallMariaDB  bool
	InstallRedis    bool
	InstallFail2ban bool
	InstallUfw      bool
	InstallGit      bool
	InstallVim      bool
	CustomPackages  string // Space-separated extra packages

	// Installation state
	InstallStarted  bool
	InstallComplete bool
	InstallError    string
	ProgressPercent float64
	ProgressMessage string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		KeyboardLayout:      "us",
		NetworkDHCP:         true,
		Hostname:            "arch-server",
		MirrorCountry:       "auto",
		UseAutoPartitioning: true,
		FilesystemType:      "ext4",
		BootloaderType:      "grub",
		UEFIMode:            true,
		Locales:             []string{"en_US.UTF-8"},
		TimezoneRegion:      "UTC",
		TimezoneCity:        "",
		KernelType:          "linux",
		CreateUser:          true,
		EnableSSH:           true,
		SSHPort:             22,
		AllowRootLogin:      false,
		InstallBaseDev:      false,
		EnableArchCN:        false,
		ArchCNMirror:        "https://mirrors.tuna.tsinghua.edu.cn/archlinuxcn",
	}
}

// ValidateHostname checks if the hostname is valid.
func ValidateHostname(hostname string) bool {
	if len(hostname) < 1 || len(hostname) > 63 {
		return false
	}
	for _, c := range hostname {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}
	return true
}

// ValidateIP checks if the given string is a valid IPv4 address.
func ValidateIP(ip string) bool {
	if ip == "" {
		return false
	}
	segments := 0
	current := 0
	digitCount := 0
	for i := 0; i < len(ip); i++ {
		c := ip[i]
		if c >= '0' && c <= '9' {
			current = current*10 + int(c-'0')
			digitCount++
			if digitCount > 3 || current > 255 {
				return false
			}
		} else if c == '.' {
			if digitCount == 0 {
				return false
			}
			segments++
			current = 0
			digitCount = 0
		} else {
			return false
		}
	}
	if digitCount == 0 {
		return false
	}
	segments++
	return segments == 4
}
