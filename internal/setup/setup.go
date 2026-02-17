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
	"github.com/joho/godotenv"
)

const serverName = "imperial-construct.tail64150e.ts.net"

//go:embed dendrite.yaml.tmpl
var dendriteYAMLTmpl string

const wellKnownClientJSON = `{"m.homeserver":{"base_url":"https://{{ .ServerName }}"},"org.matrix.msc3575.proxy":{"url":"https://{{ .ServerName }}"}}`

const wellKnownServerJSON = `{"m.server":"{{ .ServerName }}:443"}`

type templateData struct {
	ServerName       string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	SharedSecret     string
}

// Run generates all configuration files for Dendrite (host mode, writes to ./config/).
func Run() error {
	// Load .env file if present (does not override existing env vars)
	_ = godotenv.Load()

	data := newTemplateData()
	data.SharedSecret = generateSecret(32)

	files := []struct {
		path string
		tmpl string
	}{
		{"config/matrix_key.pem", ""},
		{"config/dendrite.yaml", dendriteYAMLTmpl},
		{"caddy/well-known/matrix/client", wellKnownClientJSON},
		{"caddy/well-known/matrix/server", wellKnownServerJSON},
	}

	if err := generateFiles(files, data); err != nil {
		return err
	}

	fmt.Println("\nvox-loop init complete. Generated files:")
	for _, f := range files {
		fmt.Printf("  %s\n", f.path)
	}
	fmt.Printf("\nRegistration shared secret: %s\n", data.SharedSecret)
	fmt.Println("Save this secret — you'll need it for 'vox-loop admin create-account'.")
	return nil
}

// GenerateContainerConfig generates dendrite.yaml at /etc/dendrite/ for container entrypoint use.
func GenerateContainerConfig() error {
	data := newTemplateData()
	data.SharedSecret = envOrDefault("REGISTRATION_SHARED_SECRET", generateSecret(32))

	files := []struct {
		path string
		tmpl string
	}{
		{"/etc/dendrite/dendrite.yaml", dendriteYAMLTmpl},
	}

	if err := generateFiles(files, data); err != nil {
		return err
	}

	fmt.Printf("registration shared secret: %s\n", data.SharedSecret)
	return nil
}

func newTemplateData() templateData {
	return templateData{
		ServerName:       serverName,
		PostgresUser:     envOrDefault("POSTGRES_USER", "dendrite"),
		PostgresPassword: envOrDefault("POSTGRES_PASSWORD", "changeme"),
		PostgresDB:       envOrDefault("POSTGRES_DB", "dendrite"),
		PostgresHost:     envOrDefault("POSTGRES_HOST", "postgres"),
	}
}

func generateFiles(files []struct {
	path string
	tmpl string
}, data templateData) error {
	for _, f := range files {
		if f.tmpl == "" {
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
