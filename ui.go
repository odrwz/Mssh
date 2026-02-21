package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// promptKeyGeneration interactively asks the user for key parameters and generates the key.
// Returns the absolute path to the generated private key, or "" on cancel/error.
func promptKeyGeneration() string {
	var keyType string
	if err := survey.AskOne(&survey.Select{
		Message: "Select Encryption Type:",
		Options: []string{"ed25519 (recommended)", "rsa", "ecdsa"},
		Default: "ed25519 (recommended)",
	}, &keyType); err != nil {
		fmt.Println("Cancelled.")
		return ""
	}

	// Normalise: strip the "(recommended)" suffix
	switch keyType {
	case "ed25519 (recommended)":
		keyType = "ed25519"
	}

	bits := ""
	switch keyType {
	case "rsa":
		if err := survey.AskOne(&survey.Select{
			Message: "RSA Key Length:",
			Options: []string{"2048", "3072", "4096"},
			Default: "4096",
		}, &bits); err != nil {
			fmt.Println("Cancelled.")
			return ""
		}
	case "ecdsa":
		if err := survey.AskOne(&survey.Select{
			Message: "ECDSA Curve (bits):",
			Options: []string{"256", "384", "521"},
			Default: "256",
		}, &bits); err != nil {
			fmt.Println("Cancelled.")
			return ""
		}
	}

	var comment, filename string
	survey.AskOne(&survey.Input{Message: "Comment (optional, e.g. your email):"}, &comment)
	survey.AskOne(&survey.Input{Message: fmt.Sprintf("Filename (leave blank for id_%s):", keyType)}, &filename)

	fmt.Printf("\nGenerating %s key...\n", keyType)
	if err := GenerateSSHKey(keyType, bits, comment, filename); err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	fmt.Println("Key generated successfully!")

	// Return the resolved filename so callers can link it
	if filename == "" {
		filename = "id_" + keyType
	}
	home, _ := userHomeDir()
	if !isAbs(filename) {
		filename = joinPath(home, ".ssh", filename)
	}
	return filename
}

// printHostDetails prints a formatted detail block for a host.
func printHostDetails(h SSHHost) {
	fmt.Println()
	fmt.Printf("  Alias       : %s\n", h.Alias)
	fmt.Printf("  HostName    : %s\n", h.HostName)
	fmt.Printf("  User        : %s\n", h.User)
	fmt.Printf("  Port        : %s\n", defaultStr(h.Port, "22"))
	fmt.Printf("  IdentityFile: %s\n", defaultStr(h.IdentityFile, "(none)"))
	fmt.Println()
}

func defaultStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
