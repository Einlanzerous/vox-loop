package preflight

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"syscall"
)

const (
	keyPath    = "/etc/dendrite/matrix_key.pem"
	configPath = "/etc/dendrite/dendrite.yaml"
	dendriteBin = "/usr/bin/dendrite"
)

// RunAndExec performs pre-flight checks and execs into the Dendrite binary.
func RunAndExec() error {
	if err := ensureMatrixKey(); err != nil {
		return fmt.Errorf("matrix key check: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("dendrite config not found at %s — run 'vox-loop init' first", configPath)
	}

	fmt.Println("pre-flight checks passed, starting Dendrite...")
	return execDendrite()
}

func ensureMatrixKey() error {
	if _, err := os.Stat(keyPath); err == nil {
		fmt.Println("matrix_key.pem found")
		return nil
	}

	fmt.Println("matrix_key.pem not found, generating...")
	return GenerateMatrixKey(keyPath)
}

// GenerateMatrixKey creates an ed25519 Matrix signing key in PEM format.
func GenerateMatrixKey(path string) error {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("generating ed25519 key: %w", err)
	}

	block := &pem.Block{
		Type:  "MATRIX PRIVATE KEY",
		Headers: map[string]string{
			"Key-ID": "ed25519:auto",
		},
		Bytes: priv.Seed(),
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("creating key file: %w", err)
	}
	defer f.Close()

	if err := pem.Encode(f, block); err != nil {
		return fmt.Errorf("encoding PEM: %w", err)
	}

	fmt.Printf("generated %s\n", path)
	return nil
}

func execDendrite() error {
	argv := []string{dendriteBin, "--config", configPath}
	return syscall.Exec(dendriteBin, argv, os.Environ())
}
