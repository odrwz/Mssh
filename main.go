package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

// helpers wired to stdlib so ui.go can call them without importing os/path
func userHomeDir() (string, error) { return os.UserHomeDir() }
func isAbs(p string) bool          { return filepath.IsAbs(p) }
func joinPath(parts ...string) string { return filepath.Join(parts...) }

func main() {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║    CLImssh — SSH Host & Key Manager  ║")
	fmt.Println("╚══════════════════════════════════════╝")

	for {
		var action string
		if err := survey.AskOne(&survey.Select{
			Message: "Main Menu:",
			Options: []string{
				"Manage SSH Hosts",
				"New SSH Host",
				"Manage SSH Keys",
				"Exit",
			},
		}, &action); err != nil {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		switch action {
		case "Manage SSH Hosts":
			manageHosts()
		case "New SSH Host":
			createHost()
		case "Manage SSH Keys":
			manageKeys()
		case "Exit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		}
	}
}

// manageHosts re-reads hosts each entry so list stays fresh.
func manageHosts() {
	for {
		hosts, err := ListSSHHosts()
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}

		if len(hosts) == 0 {
			fmt.Println("\nNo SSH hosts found in ~/.ssh/config.\n")
			return
		}

		options := make([]string, 0, len(hosts)+1)
		for _, h := range hosts {
			port := h.Port
			if port == "" {
				port = "22"
			}
			options = append(options, fmt.Sprintf("%-20s  %s@%s:%s", h.Alias, h.User, h.HostName, port))
		}
		options = append(options, "← Back")

		var selected string
		if err := survey.AskOne(&survey.Select{
			Message: "Select a host:",
			Options: options,
		}, &selected); err != nil || selected == "← Back" {
			return
		}

		// Find the selected host struct
		idx := -1
		for i, o := range options {
			if o == selected {
				idx = i
				break
			}
		}
		if idx < 0 || idx >= len(hosts) {
			return
		}
		manageOneHost(hosts[idx])
	}
}

// manageOneHost shows details and action menu for a single host.
func manageOneHost(h SSHHost) {
	for {
		printHostDetails(h)

		var action string
		if err := survey.AskOne(&survey.Select{
			Message: fmt.Sprintf("Actions for [%s]:", h.Alias),
			Options: []string{"Edit", "Delete", "← Back"},
		}, &action); err != nil || action == "← Back" {
			return
		}

		switch action {
		case "Edit":
			updated := editHostPrompt(h)
			if updated == nil {
				continue
			}
			if err := UpdateSSHHost(h, *updated); err != nil {
				fmt.Println("Error updating host:", err)
			} else {
				fmt.Println("Host updated successfully!")
				h = *updated
			}
		case "Delete":
			var confirm bool
			survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("Are you sure you want to delete host [%s]?", h.Alias),
				Default: false,
			}, &confirm)
			if confirm {
				if err := DeleteSSHHost(h.Alias); err != nil {
					fmt.Println("Error deleting host:", err)
				} else {
					fmt.Println("Host deleted.")
					return
				}
			}
		}
	}
}

// editHostPrompt prompts the user to change host fields.
func editHostPrompt(h SSHHost) *SSHHost {
	var answers struct {
		Alias        string
		HostName     string
		User         string
		Port         string
		IdentityFile string
	}
	answers.Alias = h.Alias
	answers.HostName = h.HostName
	answers.User = h.User
	answers.Port = h.Port
	answers.IdentityFile = h.IdentityFile

	if err := survey.Ask([]*survey.Question{
		{Name: "Alias", Prompt: &survey.Input{Message: "Alias:", Default: h.Alias}},
		{Name: "HostName", Prompt: &survey.Input{Message: "HostName:", Default: h.HostName}},
		{Name: "User", Prompt: &survey.Input{Message: "User:", Default: h.User}},
		{Name: "Port", Prompt: &survey.Input{Message: "Port:", Default: defaultStr(h.Port, "22")}},
		{Name: "IdentityFile", Prompt: &survey.Input{Message: "IdentityFile:", Default: h.IdentityFile}},
	}, &answers); err != nil {
		return nil
	}

	return &SSHHost{
		Alias:        answers.Alias,
		HostName:     answers.HostName,
		User:         answers.User,
		Port:         answers.Port,
		IdentityFile: answers.IdentityFile,
	}
}

// createHost interactively builds and saves a new SSH host entry.
func createHost() {
	var answers struct {
		Alias    string
		HostName string
		User     string
		Port     string
	}

	if err := survey.Ask([]*survey.Question{
		{Name: "Alias", Prompt: &survey.Input{Message: "Host Alias (e.g., myserver):"}, Validate: survey.Required},
		{Name: "HostName", Prompt: &survey.Input{Message: "HostName or IP:"}, Validate: survey.Required},
		{Name: "User", Prompt: &survey.Input{Message: "User (e.g., ubuntu, root):"}},
		{Name: "Port", Prompt: &survey.Input{Message: "Port:", Default: "22"}},
	}, &answers); err != nil {
		fmt.Println("Cancelled.")
		return
	}

	// Key selection
	var keyAction string
	if err := survey.AskOne(&survey.Select{
		Message: "IdentityFile (SSH Key):",
		Options: []string{"Select Existing Key", "Generate New Key", "Skip"},
	}, &keyAction); err != nil {
		fmt.Println("Cancelled.")
		return
	}

	identityFile := ""
	switch keyAction {
	case "Select Existing Key":
		keys, _ := ListSSHKeys()
		if len(keys) == 0 {
			fmt.Println("No existing keys found in ~/.ssh/.")
		} else {
			survey.AskOne(&survey.Select{Message: "Select Key:", Options: keys}, &identityFile)
		}
	case "Generate New Key":
		generated := promptKeyGeneration()
		if generated != "" {
			identityFile = generated
			fmt.Printf("New key will be linked: %s\n", identityFile)
		}
	}

	if err := AddSSHHost(SSHHost{
		Alias:        answers.Alias,
		HostName:     answers.HostName,
		User:         answers.User,
		Port:         answers.Port,
		IdentityFile: identityFile,
	}); err != nil {
		fmt.Println("Failed to save host:", err)
	} else {
		fmt.Printf("\nHost [%s] added to ~/.ssh/config successfully!\n", answers.Alias)
	}
}

// manageKeys shows a key management submenu.
func manageKeys() {
	for {
		var action string
		if err := survey.AskOne(&survey.Select{
			Message: "Key Management:",
			Options: []string{"List Keys", "Generate New Key", "← Back"},
		}, &action); err != nil || action == "← Back" {
			return
		}

		switch action {
		case "List Keys":
			keys, err := ListSSHKeys()
			if err != nil {
				fmt.Println("Error listing keys:", err)
				continue
			}
			if len(keys) == 0 {
				fmt.Println("\nNo SSH private keys found in ~/.ssh/.\n")
			} else {
				fmt.Println("\nSSH Private Keys:")
				for _, k := range keys {
					fmt.Println(" •", k)
				}
				fmt.Println()
			}
		case "Generate New Key":
			promptKeyGeneration()
		}
	}
}
