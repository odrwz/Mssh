#!/usr/bin/env bash
# CLImssh installer
set -e
curl -fsSL https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh \
  -o /usr/local/bin/climssh && chmod +x /usr/local/bin/climssh
echo "climssh installed. Run: climssh"
