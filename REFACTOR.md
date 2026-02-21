# CLImssh 重构指南

> **目标**：将现有 Go 项目重构为纯 Bash Shell 脚本，做到零依赖、体积极小、macOS 开箱即用。

---

## 现状问题

| 问题 | 详情 |
|------|------|
| 仓库体积臃肿 | `go1.22.1.linux-amd64.tar.gz`（68.9 MB）被错误提交到 Git |
| 编译产物入库 | `climssh` 二进制（5.1 MB）也被提交 |
| 依赖复杂 | 用户需要安装 Go 工具链才能从源码构建 |
| 第三方库过度引入 | `survey/v2`（TUI）和 `ssh_config` 对于此工具的需求而言完全可以用标准工具替代 |
| 跨平台问题 | 仓库内二进制为 Linux amd64，macOS 用户无法运行 |

---

## 重构方向：Go → Bash Shell Script

macOS 自带 `bash`、`ssh-keygen`、`awk`、`grep`、`sed`，足以实现所有功能，无需安装任何东西。

| 对比项 | 重构前 | 重构后 |
|--------|--------|--------|
| 语言 | Go 1.22 | Bash |
| 安装依赖 | Go 工具链 | 无 |
| 核心文件大小 | 5.1 MB 二进制 | < 15 KB 脚本 |
| 仓库总大小 | ~75 MB | < 50 KB |
| 安装方式 | GitHub Release 下载 / `go install` | `curl` 一行命令 |

---

## 新目录结构

重构完成后，仓库只保留以下文件：

```
CLImssh/
├── climssh        # 主脚本（单文件，包含所有功能，可直接执行）
├── install.sh     # 简化安装脚本
├── mssh.rb        # Homebrew formula（更新为安装脚本而非二进制）
├── README.md      # 更新安装说明
└── .gitignore     # 排除二进制和临时文件
```

---

## 执行步骤

### Step 1 — 清理仓库垃圾文件

```bash
# 移除不应入库的文件
git rm --cached go1.22.1.linux-amd64.tar.gz
git rm --cached climssh climssh-manager
git rm go.sum go.mod
git rm main.go sshconfig.go sshkey.go ui.go
```

更新 `.gitignore`，添加以下内容：

```gitignore
# 编译产物
climssh-darwin-*
climssh-linux-*
*.tar.gz

# Go 相关
go.sum
vendor/
```

---

### Step 2 — 创建主脚本 `climssh`

新建文件 `climssh`，内容结构如下（共约 250 行）：

```bash
#!/usr/bin/env bash
# CLImssh — SSH Host & Key Manager
# 依赖：bash, ssh-keygen, awk, grep, sed（macOS 内置）

SSH_CONFIG="$HOME/.ssh/config"

# ── 工具函数 ────────────────────────────────────────────

# 读取 ~/.ssh/config，输出所有 Host 条目（每行格式：alias|hostname|user|port|identityfile）
list_hosts() {
  [[ ! -f "$SSH_CONFIG" ]] && return
  awk '
    /^Host / && $2 != "*" {
      if (alias != "") print alias"|"hostname"|"user"|"port"|"idf
      alias=$2; hostname=""; user=""; port="22"; idf=""
    }
    /^[[:space:]]+HostName /     { hostname=$2 }
    /^[[:space:]]+User /         { user=$2 }
    /^[[:space:]]+Port /         { port=$2 }
    /^[[:space:]]+IdentityFile / { idf=$2 }
    END { if (alias != "") print alias"|"hostname"|"user"|"port"|"idf }
  ' "$SSH_CONFIG"
}

# 追加一个新 Host 块到 config
add_host() {
  local alias=$1 hostname=$2 user=$3 port=$4 idf=$5
  mkdir -p "$HOME/.ssh" && chmod 700 "$HOME/.ssh"
  [[ ! -f "$SSH_CONFIG" ]] && touch "$SSH_CONFIG" && chmod 600 "$SSH_CONFIG"
  {
    echo ""
    echo "Host $alias"
    [[ -n "$hostname" ]] && echo "  HostName $hostname"
    [[ -n "$user" ]]     && echo "  User $user"
    [[ "$port" != "22" && -n "$port" ]] && echo "  Port $port"
    [[ -n "$idf" ]]      && echo "  IdentityFile $idf"
  } >> "$SSH_CONFIG"
}

# 删除指定 alias 的 Host 块
delete_host() {
  local alias=$1
  awk -v target="Host $alias" '
    /^Host / { skip = ($0 == target) }
    !skip
  ' "$SSH_CONFIG" > /tmp/ssh_config_tmp && mv /tmp/ssh_config_tmp "$SSH_CONFIG"
}

# 列出 ~/.ssh 下的私钥（有对应 .pub 文件的）
list_keys() {
  for f in "$HOME/.ssh"/*; do
    [[ -f "$f" && -f "$f.pub" ]] && echo "$f"
  done
}

# ── 菜单模块 ────────────────────────────────────────────

manage_hosts() {
  while true; do
    echo ""
    mapfile -t hosts < <(list_hosts)
    if [[ ${#hosts[@]} -eq 0 ]]; then
      echo "No SSH hosts found in ~/.ssh/config."
      return
    fi
    echo "── SSH Hosts ──────────────────────"
    local i=1
    for h in "${hosts[@]}"; do
      IFS='|' read -r alias hostname user port idf <<< "$h"
      printf "  %d) %-20s %s@%s:%s\n" $i "$alias" "$user" "$hostname" "$port"
      ((i++))
    done
    echo "  0) ← Back"
    read -rp "Select: " sel
    [[ "$sel" == "0" || -z "$sel" ]] && return
    local idx=$(( sel - 1 ))
    IFS='|' read -r alias hostname user port idf <<< "${hosts[$idx]}"
    manage_one_host "$alias" "$hostname" "$user" "$port" "$idf"
  done
}

manage_one_host() {
  local alias=$1 hostname=$2 user=$3 port=$4 idf=$5
  while true; do
    echo ""
    printf "  Alias       : %s\n  HostName    : %s\n  User        : %s\n  Port        : %s\n  IdentityFile: %s\n" \
      "$alias" "$hostname" "$user" "${port:-22}" "${idf:-(none)}"
    echo ""
    echo "  1) Edit   2) Delete   0) ← Back"
    read -rp "Action: " act
    case $act in
      1) edit_host "$alias" "$hostname" "$user" "$port" "$idf"; return ;;
      2)
        read -rp "Delete [$alias]? (y/N): " confirm
        [[ "$confirm" == "y" ]] && delete_host "$alias" && echo "Host deleted." && return
        ;;
      0|"") return ;;
    esac
  done
}

edit_host() {
  local old_alias=$1
  read -rp "Alias [$1]: "        new_alias;    new_alias=${new_alias:-$1}
  read -rp "HostName [$2]: "     new_hostname; new_hostname=${new_hostname:-$2}
  read -rp "User [$3]: "         new_user;     new_user=${new_user:-$3}
  read -rp "Port [${4:-22}]: "   new_port;     new_port=${new_port:-${4:-22}}
  read -rp "IdentityFile [$5]: " new_idf;      new_idf=${new_idf:-$5}
  delete_host "$old_alias"
  add_host "$new_alias" "$new_hostname" "$new_user" "$new_port" "$new_idf"
  echo "Host updated."
}

create_host() {
  echo ""
  read -rp "Host Alias (e.g. myserver): " alias
  [[ -z "$alias" ]] && echo "Cancelled." && return
  read -rp "HostName or IP: " hostname
  [[ -z "$hostname" ]] && echo "Cancelled." && return
  read -rp "User (e.g. ubuntu): " user
  read -rp "Port [22]: " port; port=${port:-22}

  echo "  1) Select Existing Key  2) Generate New Key  3) Skip"
  read -rp "IdentityFile: " key_opt
  local idf=""
  case $key_opt in
    1)
      mapfile -t keys < <(list_keys)
      if [[ ${#keys[@]} -eq 0 ]]; then
        echo "No existing keys found."
      else
        local i=1
        for k in "${keys[@]}"; do echo "  $i) $k"; ((i++)); done
        read -rp "Select key: " ksel
        idf="${keys[$((ksel-1))]}"
      fi ;;
    2) idf=$(generate_key) ;;
  esac

  add_host "$alias" "$hostname" "$user" "$port" "$idf"
  echo "Host [$alias] added to ~/.ssh/config."
}

manage_keys() {
  while true; do
    echo ""
    echo "  1) List Keys  2) Generate New Key  0) ← Back"
    read -rp "Select: " opt
    case $opt in
      1)
        mapfile -t keys < <(list_keys)
        if [[ ${#keys[@]} -eq 0 ]]; then
          echo "No SSH private keys found in ~/.ssh/."
        else
          echo "SSH Private Keys:"; for k in "${keys[@]}"; do echo "  • $k"; done
        fi ;;
      2) generate_key ;;
      0|"") return ;;
    esac
  done
}

generate_key() {
  echo "  1) ed25519 (recommended)  2) rsa  3) ecdsa"
  read -rp "Key Type: " t
  case $t in
    1) key_type="ed25519"; bits="" ;;
    2) key_type="rsa";
       echo "  1) 2048  2) 3072  3) 4096"
       read -rp "Bits [3]=4096: " b
       case $b in 1) bits=2048;; 2) bits=3072;; *) bits=4096;; esac ;;
    3) key_type="ecdsa"
       echo "  1) 256  2) 384  3) 521"
       read -rp "Curve [1]=256: " b
       case $b in 2) bits=384;; 3) bits=521;; *) bits=256;; esac ;;
    *) echo "Cancelled."; return ;;
  esac

  read -rp "Comment (optional): " comment
  read -rp "Filename [id_${key_type}]: " filename
  filename=${filename:-id_${key_type}}
  local filepath="$HOME/.ssh/$filename"

  if [[ -f "$filepath" ]]; then
    echo "Error: $filepath already exists. Delete it first or choose a different name."
    return
  fi

  local args=(-t "$key_type" -f "$filepath" -q -N "")
  [[ -n "$bits" ]]    && args+=(-b "$bits")
  [[ -n "$comment" ]] && args+=(-C "$comment")
  ssh-keygen "${args[@]}" && echo "Key generated: $filepath"
  echo "$filepath"
}

# ── 主入口 ──────────────────────────────────────────────

main() {
  echo "╔══════════════════════════════════════╗"
  echo "║    CLImssh — SSH Host & Key Manager  ║"
  echo "╚══════════════════════════════════════╝"
  while true; do
    echo ""
    echo "  1) Manage SSH Hosts"
    echo "  2) New SSH Host"
    echo "  3) Manage SSH Keys"
    echo "  0) Exit"
    read -rp "Select: " opt
    case $opt in
      1) manage_hosts ;;
      2) create_host  ;;
      3) manage_keys  ;;
      0|"") echo "Goodbye!"; exit 0 ;;
    esac
  done
}

main
```

赋予执行权限：

```bash
chmod +x climssh
```

---

### Step 3 — 简化 `install.sh`

将 `install.sh` 替换为以下三行脚本：

```bash
#!/usr/bin/env bash
curl -fsSL https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh \
  -o /usr/local/bin/climssh && chmod +x /usr/local/bin/climssh
echo "✅ climssh installed. Run: climssh"
```

---

### Step 4 — 更新 Homebrew Formula `mssh.rb`

```ruby
class Mssh < Formula
  desc "SSH Host & Key Manager CLI"
  homepage "https://github.com/odrwz/CLImssh"
  url "https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh"
  # sha256 替换为脚本文件的实际 sha256sum 输出
  version "2.0.0"

  def install
    bin.install "climssh"
  end

  test do
    assert_match "CLImssh", shell_output("#{bin}/climssh --help 2>&1", 1)
  end
end
```

> [!NOTE]
> 需要在发布 tag 后用 `sha256sum climssh` 获取哈希值填入 formula。

---

### Step 5 — 更新 `README.md`

将安装说明更新为：

```markdown
## 安装

**方式一（推荐）：一行命令**
```bash
curl -fsSL https://raw.githubusercontent.com/odrwz/CLImssh/main/install.sh | bash
```

**方式二：Homebrew**
```bash
brew install odrwz/tap/mssh
```

**方式三：手动下载**
```bash
curl -fsSL https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh \
  -o /usr/local/bin/climssh && chmod +x /usr/local/bin/climssh
```
```

---

### Step 6 — 提交

```bash
git add climssh install.sh mssh.rb README.md .gitignore
git commit -m "refactor: rewrite in bash, remove Go dependencies"
git tag v2.0.0
git push && git push --tags
```

---

## 验证清单

在 macOS 上完成以下验证后，重构即告完成：

- [ ] `bash climssh` 可正常启动主菜单
- [ ] 可新增、编辑、删除 SSH Host，`~/.ssh/config` 内容修改正确
- [ ] 可列出和生成三种类型的 SSH 密钥（ed25519 / rsa / ecdsa）
- [ ] `curl | bash` 安装方式在 macOS 14+ 上可用
- [ ] `wc -c climssh` 输出小于 15360（15 KB）
- [ ] 仓库无二进制文件、无 `.tar.gz` 文件
