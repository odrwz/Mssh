# CLImssh

`climssh` is an interactive command-line tool designed for Mac users to easily manage their `~/.ssh/config` file and SSH keys.

## Features
- **Interactive UI**: Navigate menus to view, create, and manage your SSH connections.
- **SSH Config Management**: Parse and write hosts safely to `~/.ssh/config`.
- **Key Generation**: Interactively generate RSA, Ed25519, or ECDSA keys and link them directly to your connections.

## Installation (macOS)
You can install this via Homebrew. (Note: Replace `odrwz` with your GitHub username once the repository is pushed).

```bash
brew tap odrwz/climssh
brew install climssh
```

Alternatively, build from source using Go:
```bash
go build -o climssh .
sudo mv climssh /usr/local/bin/
```

## Usage
Simply run the program in your terminal:
```bash
climssh
```
It will present you with an interactive menu. Use the Arrow Keys to navigate and Enter to select.
