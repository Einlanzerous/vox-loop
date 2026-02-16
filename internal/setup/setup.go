package setup

import (
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Einlanzerous/vox-loop/internal/preflight"
)

const serverName = "imperial-construct.tail64150e.ts.net"

//go:embed dendrite.yaml.tmpl
var dendriteYAMLTmpl string

const wellKnownClientJSON = `{"m.homeserver":{"base_url":"https://{{ .ServerName }}"},"org.matrix.msc3575.proxy":{"url":"https://{{ .ServerName }}"}}`

const wellKnownServerJSON = `{"m.server":"{{ .ServerName }}:443"}`

type templateData struct {
	ServerName       string
	PostgresPassword string
	SharedSecret     string
}

// Run generates all configuration files for Dendrite.
func Run() error {
	data := templateData{
		ServerName:       serverName,
		PostgresPassword: envOrDefault("POSTGRES_PASSWORD", "changeme"),
		SharedSecret:     generateSecret(32),
	}

	files := []struct {
		path string
		tmpl string
	}{
		{"config/matrix_key.pem", ""},
		{"config/dendrite.yaml", dendriteYAMLTmpl},
		{"caddy/well-known/matrix/client", wellKnownClientJSON},
		{"caddy/well-known/matrix/server", wellKnownServerJSON},
	}

	for _, f := range files {
		if f.tmpl == "" {
			// Special case: generate matrix key
			dir := filepath.Dir(f.path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", dir, err)
			}
			if err := preflight.GenerateMatrixKey(f.path); err != nil {
				return err
			}
			continue
		}

		if err := renderTemplate(f.path, f.tmpl, data); err != nil {
			return err
		}
	}

	fmt.Println("\nvox-loop init complete. Generated files:")
	for _, f := range files {
		fmt.Printf("  %s\n", f.path)
	}
	fmt.Printf("\nRegistration shared secret: %s\n", data.SharedSecret)
	fmt.Println("Save this secret — you'll need it for 'vox-loop admin create-account'.")
	return nil
}

func renderTemplate(path, tmplStr string, data templateData) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	t, err := template.New(filepath.Base(path)).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parsing template for %s: %w", path, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("rendering %s: %w", path, err)
	}

	fmt.Printf("wrote %s\n", path)
	return nil
}

func generateSecret(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate random bytes: %v", err))
	}
	return hex.EncodeToString(b)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
