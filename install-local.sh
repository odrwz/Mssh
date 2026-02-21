#!/usr/bin/env bash
# 本地测试安装脚本（直接复制当前目录的 climssh）
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET="/usr/local/bin/climssh"

echo "Installing climssh from $SCRIPT_DIR..."
cp "$SCRIPT_DIR/climssh" "$TARGET"
chmod +x "$TARGET"
echo "✅ climssh installed to $TARGET"
echo "Run: climssh"
