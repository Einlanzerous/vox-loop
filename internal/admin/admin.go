package admin

import (
	"fmt"
	"os"
	"os/exec"
)

const createAccountBin = "/usr/bin/create-account"

// CreateAccount creates a Matrix user account via the Dendrite create-account tool.
func CreateAccount(username string, isAdmin bool) error {
	configPath := os.Getenv("DENDRITE_CONFIG")
	if configPath == "" {
		configPath = "/etc/dendrite/dendrite.yaml"
	}

	args := []string{
		"--config", configPath,
		"--username", username,
	}

	if isAdmin {
		args = append(args, "--admin")
	}

	fmt.Printf("creating account %q (admin=%v)...\n", username, isAdmin)

	cmd := exec.Command(createAccountBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create-account failed: %w", err)
	}

	fmt.Printf("account %q created successfully\n", username)
	return nil
}
