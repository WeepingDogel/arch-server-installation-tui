package model

// Config holds all user choices throughout the installation wizard.
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
	NetworkIface string

	// Mirror
	MirrorURL     string
	MirrorCountry string
	CustomMirror  string
	EnableArchCN  bool
	ArchCNMirror  string

	// Disk
	DiskDevice      string
	DiskSize        string
	PartitionScheme string // "gpt" or "mbr"
	PartitionMode   string // "auto" or "manual"
	EncryptDisk     bool
	LVMEnabled      bool
	EfiSize         string // e.g. "512M"
	SwapSize        string // e.g. "4G" or "" to skip
	HomeSize        string // e.g. "50G" or "" to merge into root
	RootSize        string // e.g. "100%" for auto remainder

	// Filesystem
	FilesystemType string // ext4, btrfs, xfs, f2fs

	// Bootloader
	BootloaderType string // grub, systemd-boot
	UEFIMode       bool

	// Timezone & Locale
	TimezoneRegion string
	TimezoneCity   string
	Locales        []string

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
	KernelType      string
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
	CustomPackages  string

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
		KeyboardLayout:  "us",
		NetworkDHCP:     true,
		Hostname:        "arch-server",
		MirrorCountry:   "auto",
		FilesystemType:  "ext4",
		BootloaderType:  "grub",
		UEFIMode:        true,
		Locales:         []string{"en_US.UTF-8"},
		TimezoneRegion:  "UTC",
		TimezoneCity:    "",
		KernelType:      "linux",
		CreateUser:      true,
		EnableSSH:       true,
		SSHPort:         22,
		AllowRootLogin:  false,
		InstallBaseDev:  false,
		EnableArchCN:    false,
		ArchCNMirror:    "https://mirrors.tuna.tsinghua.edu.cn/archlinuxcn",
		PartitionScheme: "gpt",
		PartitionMode:   "auto",
		EfiSize:         "512M",
		SwapSize:        "",
		HomeSize:        "",
		RootSize:        "100%",
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
