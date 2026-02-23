# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-02-23

### Added
- GitHub Actions CI workflow running `go vet` and `go test` on Linux, macOS, and Windows.
- GitHub Actions release workflow using GoReleaser.
- GoReleaser config for `linux/darwin` on `amd64/arm64`.
- `nit --version` support with build-time version injection via `-ldflags`.
- Configurable commit editor keybindings in `nit.toml` (`[keys.commit_editor.*]`).
