package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinburke/ssh_config"
)

// SSHHost represents a parsed SSH host entry
type SSHHost struct {
	Alias        string
	HostName     string
	User         string
	Port         string
	IdentityFile string
}

// GetSSHConfigPath returns the path to the user's ~/.ssh/config
func GetSSHConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
		return ""
	}
	return filepath.Join(home, ".ssh", "config")
}

// ListSSHHosts parses the SSH config and returns a list of hosts
func ListSSHHosts() ([]SSHHost, error) {
	configPath := GetSSHConfigPath()
	f, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []SSHHost{}, nil
		}
		return nil, err
	}
	defer f.Close()

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil, err
	}

	var hosts []SSHHost
	for _, host := range cfg.Hosts {
		// Skip global fallback wildcard (*)
		if len(host.Patterns) > 0 && host.Patterns[0].String() == "*" {
			continue
		}

		alias := host.Patterns[0].String()
		if alias == "" {
			continue
		}

		h := SSHHost{Alias: alias}
		for _, node := range host.Nodes {
			kv, ok := node.(*ssh_config.KV)
			if !ok {
				continue
			}
			switch strings.ToLower(kv.Key) {
			case "hostname":
				h.HostName = kv.Value
			case "user":
				h.User = kv.Value
			case "port":
				h.Port = kv.Value
			case "identityfile":
				h.IdentityFile = kv.Value
			}
		}
		hosts = append(hosts, h)
	}

	return hosts, nil
}

// EnsureSSHConfigExists creates ~/.ssh/config if it doesn't exist
func EnsureSSHConfigExists() error {
	path := GetSSHConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		// os.WriteFile replaces the deprecated ioutil.WriteFile
		return os.WriteFile(path, []byte(""), 0600)
	}
	return nil
}

// AddSSHHost appends a new host block to ~/.ssh/config
func AddSSHHost(h SSHHost) error {
	if err := EnsureSSHConfigExists(); err != nil {
		return err
	}

	configPath := GetSSHConfigPath()
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nHost %s\n", h.Alias))
	if h.HostName != "" {
		sb.WriteString(fmt.Sprintf("  HostName %s\n", h.HostName))
	}
	if h.User != "" {
		sb.WriteString(fmt.Sprintf("  User %s\n", h.User))
	}
	if h.Port != "" && h.Port != "22" {
		// Only write port if non-default to keep config clean
		sb.WriteString(fmt.Sprintf("  Port %s\n", h.Port))
	}
	if h.IdentityFile != "" {
		sb.WriteString(fmt.Sprintf("  IdentityFile %s\n", h.IdentityFile))
	}

	_, err = f.WriteString(sb.String())
	return err
}

// UpdateSSHHost rewrites a host entry in-place in the config file
func UpdateSSHHost(original SSHHost, updated SSHHost) error {
	configPath := GetSSHConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var out []string
	inTarget := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Detect the start of our target host block
		if strings.EqualFold(trimmed, "Host "+original.Alias) {
			inTarget = true
			// Write updated block instead
			out = append(out, fmt.Sprintf("Host %s", updated.Alias))
			if updated.HostName != "" {
				out = append(out, fmt.Sprintf("  HostName %s", updated.HostName))
			}
			if updated.User != "" {
				out = append(out, fmt.Sprintf("  User %s", updated.User))
			}
			if updated.Port != "" && updated.Port != "22" {
				out = append(out, fmt.Sprintf("  Port %s", updated.Port))
			}
			if updated.IdentityFile != "" {
				out = append(out, fmt.Sprintf("  IdentityFile %s", updated.IdentityFile))
			}
			continue
		}
		// Detect end of target block (next Host or end of file)
		if inTarget {
			if strings.HasPrefix(strings.TrimSpace(line), "Host ") {
				inTarget = false
			} else {
				continue // skip old block lines
			}
		}
		out = append(out, line)
	}

	return os.WriteFile(configPath, []byte(strings.Join(out, "\n")), 0600)
}

// DeleteSSHHost removes a host block from the config file
func DeleteSSHHost(alias string) error {
	configPath := GetSSHConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var out []string
	inTarget := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, "Host "+alias) {
			inTarget = true
			continue
		}
		if inTarget {
			if strings.HasPrefix(trimmed, "Host ") {
				inTarget = false
			} else {
				continue
			}
		}
		out = append(out, line)
	}

	return os.WriteFile(configPath, []byte(strings.Join(out, "\n")), 0600)
}
