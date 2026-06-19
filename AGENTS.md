# AGENTS.md — AI Context for Arch Server Installation TUI

## Project Overview

An interactive TUI (Terminal User Interface) tool written in Go for installing Arch Linux as a production-ready server. Uses Bubble Tea framework for the terminal UI, Lip Gloss for styling.

## Key Architecture

```
arch-server-installation-tui/
├── cmd/installer/main.go     # Entry point
├── internal/
│   ├── model/config.go       # Shared config struct (all user choices)
│   ├── tui/                  # All TUI screens (Bubble Tea models)
│   │   ├── root.go           # Navigation orchestrator
│   │   ├── theme.go          # Lip Gloss style definitions
│   │   ├── logo.go           # ASCII art + step indicator
│   │   ├── layout.go         # Screen layout helpers
│   │   └── *.go             # Step 1-13 models
│   ├── installer/installer.go  # Real installation pipeline
│   ├── mirror/mirrors.go       # Mirror definitions + filters
│   └── utils/                  # Utility packages (env, exec, validate)
├── .github/workflows/ci.yml    # GitHub Actions CI/CD
├── .golangci.yml               # Linter config
├── go.mod / go.sum
├── .gitignore
├── LICENSE (MIT)
├── README.md
├── AGENTS.md
└── CONTRIBUTING.md
```

## TUI Architecture Pattern

Every screen is a Bubble Tea model with this interface:

```go
type StepModel struct {
    config *model.Config  // Shared pointer to config
    Next   bool           // Set true to advance to next step
    cursor int            // Current selection index
}

func NewStepModel(cfg *model.Config) StepModel
func (m StepModel) Init() tea.Cmd
func (m StepModel) Update(msg tea.Msg) (StepModel, tea.Cmd)
func (m StepModel) View() string
```

The root model (`root.go`) delegates `Update`/`View` to the current step model based on `m.step` (1-indexed, 13 total steps).

## Important Conventions

1. **Config is shared by pointer** — All step models receive `*model.Config` and mutate it directly
2. **Navigation via `Next` bool** — When a step completes, set `m.Next = true`; root model advances
3. **Theme is global** — `theme.go` exports all Lip Gloss styles as package-level vars
4. **No external deps beyond Bubble Tea ecosystem** — Uses only `bubbletea`, `lipgloss`, `bubbles`
5. **Step numbering** — Steps are 1-based: 1=Welcome ... 13=Install

## Step Flow

1 → 2 → 3 → 4 → 5 → 6 → 7 → 8 → 9 → 10 → 11 → 12 → 13

- `Esc` goes back one step (except step 1)
- `Tab`/`Enter` advances to next step (after validation)
- `Ctrl+C`/`q` quits (except during installation)

## Code Style

- Standard Go formatting (`gofmt`)
- Error handling: validate at step transition, show errors in UI footer
- Tests: `model/config_test.go` with table-driven tests
- Go 1.22 minimum
- `go vet ./...` must pass

## CI/CD

- Lint + Test run on PRs only
- Build runs on push to main and manual `workflow_dispatch`
- Release builds on tag publish

## Utility Packages

- `internal/utils/env.go` — Environment detection (Arch ISO, internet, disks, memory)
- `internal/utils/exec.go` — Safe command execution wrappers
- `internal/utils/validate.go` — Input validation (hostname, IP, port, password)
