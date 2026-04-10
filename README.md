# mssh

SSH Host & Key Manager for macOS — a single Bash script, zero dependencies.

## Features

- View, add, edit, delete SSH hosts in `~/.ssh/config`
- List and generate SSH keys (ed25519 / rsa / ecdsa)
- Associate keys with hosts during creation

## Requirements

- macOS (bash, ssh-keygen, awk, sed — all built-in)
- No Go, no Node, no pip — nothing to install

## Install

**One-line install:**
```bash
curl -fsSL https://raw.githubusercontent.com/odrwz/mssh/main/install.sh | bash
```

**Manual:**
```bash
curl -fsSL https://raw.githubusercontent.com/odrwz/mssh/main/mssh \
  -o /usr/local/bin/mssh && chmod +x /usr/local/bin/mssh
```

**Homebrew:**
```bash
brew install odrwz/tap/mssh
```

## Usage

```bash
mssh
```

An interactive menu will guide you through all options.

### Color Themes

`mssh` supports 7 built-in color themes plus a no-color mode:

```bash
mssh --theme ocean       # Use a specific theme for this session
mssh -t forest           # Short flag
mssh --set-theme ocean   # Set ocean as the permanent default
mssh --no-color          # Disable colors for this session
```

Available themes: `default`, `dark`, `solarized`, `minimal`, `ocean`, `forest`, `sunset`.

**Set a permanent default theme:**

```bash
mssh --set-theme ocean
```

This writes the preference to `~/.config/mssh/theme`. You can override it temporarily with `--theme` or the `MSSH_THEME` environment variable.

**Environment variables:**

```bash
export MSSH_THEME="ocean"   # Default theme
export NO_COLOR=1           # Disable colors globally
mssh
```

## Uninstall

```bash
rm /usr/local/bin/mssh
```
