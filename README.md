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

## Uninstall

```bash
rm /usr/local/bin/mssh
```
