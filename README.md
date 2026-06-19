# Arch Linux Server Installer

```
                    /\
                   /  \
                  /    \
                 _\     \
                /        \
               /          \
              /     __   \_\
             /     /  \     \
            /__,--'    '--,__\
```

An interactive TUI (Terminal User Interface) tool for installing Arch Linux as a production-ready server. Inspired by Ubuntu Server's installation experience.

Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea) — elegant, fast, and maintainable.

**by WeepingDogel**

---

## Features

- **13-Step Installation Wizard** — Guided setup from keyboard layout to final installation
- **Elegant Terminal UI** — Styled with Lip Gloss, featuring the Arch Linux diamond logo, step indicators, animated spinners, and progress bars
- **30+ Mirrors** — Including special Arch Linux CN mirrors (TUNA, USTC, 163, Aliyun, Huawei Cloud, Tencent Cloud, SJTU, and more)
- **Comprehensive Server Configuration** — Network (DHCP/Static), disk partitioning, filesystem (ext4/btrfs/xfs/f2fs), bootloader (GRUB/systemd-boot), timezone, locale, users, SSH, and package selection
- **Real Installation Engine** — Production-ready `arch-chroot`, `pacstrap`, `genfstab`, bootloader installation commands
- **Input Validation** — Hostname, IP address, and password strength validation
- **CI/CD Pipeline** — GitHub Actions automation for lint, test, security scan, cross-compilation, and release
- **23 Unit Tests** — All passing

---

## Quick Start

### Prerequisites

- An Arch Linux live ISO environment (booted)
- Active internet connection
- Go 1.21+ (for building)

### Build

```bash
git clone https://github.com/WeepingDogel/arch-server-installation-tui.git
cd arch-server-installation-tui
go build -o arch-installer ./cmd/installer/
```

### Run

```bash
sudo ./arch-installer
```

> **Note:** This tool must be run as root inside an Arch Linux live environment to perform the actual installation. Running it elsewhere will display the UI but the installation steps will simulate.

---

## Installation Wizard

| Step | Screen | Description |
|------|--------|-------------|
| 1 | **Welcome** | ASCII logo, requirements checklist, begin prompt |
| 2 | **Keyboard Layout** | Choose from 29 layouts (us, uk, de, fr, jp, cn, etc.) |
| 3 | **Network** | DHCP/Static toggle, hostname, IP, netmask, gateway, DNS |
| 4 | **Mirror Selection** | 30+ mirrors with special Arch Linux CN list (press `s`/`a` to toggle) |
| 5 | **Disk** | Select target disk device (/dev/sda, /dev/nvme0n1, etc.) |
| 6 | **Filesystem** | ext4, btrfs, xfs, or f2fs |
| 7 | **Bootloader** | GRUB (BIOS/UEFI) or systemd-boot |
| 8 | **Timezone & Locale** | 24 timezones + 24 locales |
| 9 | **Users** | Root password + optional sudo user |
| 10 | **SSH** | Enable/disable, port, root login toggle |
| 11 | **Packages** | Kernel (linux/linux-lts/linux-zen/linux-hardened), Docker, Nginx, PostgreSQL, MariaDB, Redis, Fail2ban, UFW, Git, Vim |
| 12 | **Summary** | Full configuration review + confirm/back buttons |
| 13 | **Installation** | Animated progress bar with step-by-step status |

---

## Mirrors

The tool includes 30 mirrors across global regions, with special emphasis on Arch Linux CN mirrors:

| Region | Mirrors |
|--------|---------|
| **China ★** | TUNA (Tsinghua), USTC, 163 (NetEase), Aliyun, Huawei Cloud, Tencent Cloud, SJTU, Nanjing University, Chongqing University, BFSU, Neusoft, Xiyou Linux, plus Arch Linux CN repos |
| **Global** | Arch Linux Official, Kernel.org, Princeton |
| **Asia Pacific** | Japan (JAIST), Singapore (0x), Taiwan (NCHC), South Korea (KAIST), India (IIT Bombay) |
| **Europe** | Netherlands (NLUUG), France (IRCAM), UK (UKFast), Sweden (Lysator), Poland (ICM), Russia (Yandex), Austria (VUM) |

Press `s` to toggle **Special Mirrors** view, `a` for **All Mirrors**.

---

## Tech Stack

- **[Go](https://go.dev/)** 1.21+ — Type-safe, compiled, performant
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — Elm-architecture TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** — CSS-like terminal styling
- **[Bubbles](https://github.com/charmbracelet/bubbles)** — UI components (textinput, etc.)

---

## Project Structure

```
arch-server-installation-tui/
├── main.go                                  # Root entry
├── cmd/installer/main.go                    # Alternative entry point
├── internal/
│   ├── model/
│   │   ├── config.go                        # Shared configuration struct
│   │   └── config_test.go                   # 23 unit tests
│   ├── tui/
│   │   ├── root.go                          # Navigation orchestrator
│   │   ├── theme.go                         # 30+ Lip Gloss style definitions
│   │   ├── logo.go                          # ASCII art + step indicator
│   │   ├── layout.go                        # Screen layout helpers
│   │   ├── welcome.go                       # Step 1
│   │   ├── keyboard.go                      # Step 2
│   │   ├── network.go                       # Step 3
│   │   ├── mirror.go                        # Step 4
│   │   ├── disk.go                          # Step 5
│   │   ├── filesystem.go                    # Step 6
│   │   ├── bootloader.go                    # Step 7
│   │   ├── timezone.go                      # Step 8
│   │   ├── users.go                         # Step 9
│   │   ├── ssh.go                           # Step 10
│   │   ├── packages.go                      # Step 11
│   │   ├── summary.go                       # Step 12
│   │   ├── install.go                       # Step 13 (progress + completion)
│   │   └── components/                      # Reusable UI components
│   ├── installer/
│   │   └── installer.go                     # Real installation pipeline (15 steps)
│   └── mirror/
│       └── mirrors.go                       # Mirror definitions + filters
├── .github/workflows/ci.yml                 # GitHub Actions CI/CD
├── .golangci.yml                            # Linter configuration
├── go.mod / go.sum
└── README.md
```

---

## CI/CD Pipeline

The GitHub Actions workflow automates:

| Job | Description |
|-----|-------------|
| **lint** | `golangci-lint` — code quality checks |
| **test** | `go test -race -cover` — 23 tests with coverage |
| **security** | `govulncheck` — vulnerability scanning |
| **build** | Cross-compilation for `linux/amd64` + `linux/arm64` |
| **release** | On tag push: build, package, upload to GitHub Releases |

---

## Development

### Run tests

```bash
go test ./... -v
```

### Lint

```bash
golangci-lint run ./...
```

### Build for production

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o arch-installer ./cmd/installer/
```

---

## License

MIT

---

## Author

**WeepingDogel**