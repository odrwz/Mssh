package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ListSSHKeys finds SSH private keys in ~/.ssh by checking for a matching .pub file
func ListSSHKeys() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// Build a set of all .pub file stems for fast lookup
	pubSet := make(map[string]bool)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".pub") {
			stem := strings.TrimSuffix(e.Name(), ".pub")
			pubSet[stem] = true
		}
	}

	var keys []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		// Skip known non-key files
		if strings.HasSuffix(name, ".pub") ||
			name == "config" ||
			name == "known_hosts" ||
			name == "authorized_keys" ||
			strings.HasSuffix(name, ".known_hosts") {
			continue
		}
		// Only include files that have a corresponding .pub (reliable heuristic)
		if pubSet[name] {
			keys = append(keys, filepath.Join(sshDir, name))
		}
	}

	return keys, nil
}

// GenerateSSHKey creates a new SSH key using ssh-keygen.
// For RSA/ECDSA, bits specifies the key length. Pass "" for defaults.
func GenerateSSHKey(keyType string, bits string, comment string, filename string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if filename == "" {
		filename = filepath.Join(home, ".ssh", "id_"+keyType)
	} else if !filepath.IsAbs(filename) {
		filename = filepath.Join(home, ".ssh", filename)
	}

	// Warn user if the key file already exists to avoid ssh-keygen interactively blocking
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("key file already exists: %s\n(delete it first or choose a different filename)", filename)
	}

	args := []string{"-t", keyType, "-f", filename, "-q", "-N", ""}
	if bits != "" && keyType != "ed25519" {
		args = append(args, "-b", bits)
	}
	if comment != "" {
		args = append(args, "-C", comment)
	}

	cmd := exec.Command("ssh-keygen", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ssh-keygen failed: %w", err)
	}
	return nil
}
