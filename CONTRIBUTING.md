# Contributing to Arch Linux Server Installer

Thank you for considering contributing! This project aims to provide a polished, production-ready Arch Linux installation experience.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone git@github.com:YOUR_USERNAME/arch-server-installation-tui.git`
3. Create a feature branch: `git checkout -b feat/my-feature`
4. Make your changes
5. Run tests: `go test ./... -v`
6. Run vet: `go vet ./...`
7. Commit and push, then open a PR

## Development Guidelines

### Code Style

- Run `gofmt` before committing — standard Go formatting is enforced
- All exported types and functions must have Go-style doc comments
- Keep functions focused and small — each step model should handle only its own logic
- Use table-driven tests for validation functions

### Adding a New Step

1. Create a new file `internal/tui/yourscreen.go`
2. Implement the `Init()`, `Update()`, and `View()` methods with a `Next bool` field
3. Add the model to `RootModel` in `root.go`
4. Add the step name to `StepNames` in `logo.go`
5. Increment `TotalSteps` in `logo.go`
6. Wire the Update/View delegation in `root.go`

### Adding Mirrors

Edit `internal/mirror/mirrors.go` and add entries to the `DefaultMirrors()` function.

### Testing

- All validation functions must have tests in `internal/model/config_test.go`
- Run `go test ./... -race -cover` before submitting a PR
- Aim for at least 80% coverage on the model package

### CI/CD

The GitHub Actions workflow runs:

- **Lint** — `golangci-lint` (PRs only)
- **Test** — `go test -race` (PRs only)
- **Build** — cross-compile for linux/amd64 + linux/arm64 (push to main or manual dispatch)

## Pull Request Process

1. Ensure all CI checks pass (lint, test)
2. Update README.md if your change affects the UI, configuration, or workflow
3. Keep PRs focused on a single concern
4. Reference any related issues

## Code of Conduct

Be respectful, constructive, and inclusive. This is a learning-friendly project.