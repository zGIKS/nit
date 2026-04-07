# nit

A keyboard-driven terminal UI for Git. Stage files, write commits, switch and create branches, push, pull, and fetch — all without leaving your terminal.

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea).

---

## Features

- **Staging area** — stage or unstage individual files, or stage/unstage everything at once
- **Commit** — write and submit a commit message from inside the TUI
- **Branch management** — switch branches and create new ones from any source
- **Push / Pull / Fetch** — run common remote operations quickly from keys and menu
- **Commit graph** — visual branch graph rendered with Unicode box-drawing characters
- **Clipboard support** — copy text via OSC 52 (works over SSH/tmux), system clipboard (`pbcopy`, `wl-copy`, `xclip`, `xsel`), or both
- **Fully configurable** — key bindings, clipboard mode, and UI labels via a TOML file or environment variables
- **Mouse support** — optional mouse navigation in addition to the keyboard

---

## Requirements

- Go 1.22+ (to build from source)
- Git

For system clipboard support on Linux one of the following must be installed:
`wl-clipboard` (Wayland), `xclip`, or `xsel` (X11).  
macOS uses `pbcopy`/`pbpaste` which are built-in.

---

## Installation

### From source

```bash
go install github.com/zGIKS/nit/cmd/nit@latest
```

Or clone and build manually:

```bash
git clone https://github.com/zGIKS/nit.git
cd nit
go build -o nit ./cmd/nit
```

### Arch Linux (AUR)

```bash
yay -S nit-bin
# or
paru -S nit-bin
```

### Debian / Ubuntu

Download the binary from the [releases page](https://github.com/zGIKS/nit/releases) and place it in your `$PATH`:

```bash
sudo install -m755 nit /usr/local/bin/nit
```

---

## Usage

Run `nit` from inside any Git repository:

```bash
cd /path/to/your/repo
nit
```

### Default key bindings

| Key | Action |
|-----|--------|
| `Tab` | Switch panel (changes ↔ branches) |
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Stage / unstage selected file · Select branch |
| `s` | Stage all changes |
| `u` | Unstage all changes |
| `c` | Focus the commit message input |
| `f` | Fetch from remote |
| `p` / `Ctrl+P` | Push to remote |
| `q` / `Ctrl+C` | Quit |

#### Inside the commit input

| Key | Action |
|-----|--------|
| `Enter` | Commit |
| `Esc` | Cancel / close |
| `Ctrl+C` / `Ctrl+X` | Cut to clipboard |
| `Ctrl+V` | Paste from clipboard |
| `Ctrl+A` | Move cursor to beginning of line |
| `Ctrl+E` | Move cursor to end of line |

#### Branch creation dialog

| Key | Action |
|-----|--------|
| `Enter` | Create branch and push to origin |
| `Esc` | Cancel |

---

## Configuration

`nit` looks for a configuration file in the following order:

1. Path set by the `NIT_CONFIG_FILE` environment variable
2. `~/.config/nit/nit.toml` (Linux, respects `$XDG_CONFIG_HOME`)
3. `~/Library/Application Support/nit/nit.toml` (macOS)
4. `./nit.toml` in the current working directory (fallback)

Copy the example config to get started:

```bash
# Linux
mkdir -p ~/.config/nit
cp nit.example.toml ~/.config/nit/nit.toml

# macOS
mkdir -p ~/Library/Application\ Support/nit
cp nit.example.toml ~/Library/Application\ Support/nit/nit.toml
```

### Clipboard modes

| Mode | Behaviour |
|------|-----------|
| `only_copy` | Try OSC 52, then system clipboard. No paste. **(default)** |
| `osc52` | OSC 52 only (works over SSH and inside tmux/screen) |
| `system` | System clipboard only (`pbcopy`, `wl-copy`, `xclip`, `xsel`) |
| `auto` | System clipboard for both copy and paste |
| `internal` | No clipboard integration |

```toml
[clipboard]
mode = "only_copy"
# Optionally override the commands used:
# copy_cmd  = "wl-copy"
# paste_cmd = "wl-paste -n"
```

### Environment variables

| Variable | Description |
|----------|-------------|
| `NIT_CONFIG_FILE` | Override the config file path |
| `NIT_CLIPBOARD_MODE` | Override the clipboard mode |
| `NIT_CLIPBOARD_COPY_CMD` | Override the copy command |
| `NIT_CLIPBOARD_PASTE_CMD` | Override the paste command |
| `NIT_MOUSE_MODE` | Mouse mode: `cell` (default), `all`, or `off` |

### Custom key bindings

All bindings can be overridden in the config file:

```toml
[keys.quit]
keys = ["ctrl+c", "q"]

[keys.push]
keys = ["p", "ctrl+p"]
```

See [`nit.example.toml`](nit.example.toml) for the full list of options.

---

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
