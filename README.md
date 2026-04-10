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
mssh --theme ocean       # Use a specific theme
mssh -t forest           # Short flag
mssh --no-color          # Disable colors completely
```

Available themes: `default`, `dark`, `solarized`, `minimal`, `ocean`, `forest`, `sunset`.

You can also set a default theme via environment variable:

```bash
export MSSH_THEME="ocean"
mssh
```

Or disable colors permanently:

```bash
export NO_COLOR=1
mssh
```

## Uninstall

```bash
rm /usr/local/bin/mssh
```
