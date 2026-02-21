#!/usr/bin/env bash
# mssh installer
set -e
curl -fsSL https://raw.githubusercontent.com/odrwz/mssh/main/mssh \
  -o /usr/local/bin/mssh && chmod +x /usr/local/bin/mssh
echo "mssh installed. Run: mssh"
