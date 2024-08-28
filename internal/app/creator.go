package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ooyeku/grav-lsm/pkg/logging"
)

type AppCreator struct {
	logger *logging.ColorfulLogger
}

func NewAppCreator() *AppCreator {
	return &AppCreator{logger: logging.NewColorfulLogger()}
}

func (ac *AppCreator) CreateApp(name string) error {
	// Append _grav to the app name
	appName := name + "_grav"

	// Create the main app directory
	if err := os.Mkdir(appName, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{"cmd", "internal/models", "internal/handlers", "config"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(appName, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create main.go
	if err := ac.createMainFile(appName); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}

	// Create go.mod
	if err := ac.createGoMod(appName); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	ac.logger.Info("Grav app '" + appName + "' created successfully")
	return nil
}

func (ac *AppCreator) createMainFile(appName string) error {
	mainTemplate := `package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to %s!")
    })

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
`
	return ac.createFileFromTemplate(filepath.Join(appName, "cmd", "main.go"), mainTemplate, appName)
}

func (ac *AppCreator) createGoMod(appName string) error {
	cmd := exec.Command("go", "mod", "init", appName)
	cmd.Dir = appName
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize go module: %w\n%s", err, output)
	}
	ac.logger.Info("Go module initialized for " + appName)
	return nil
}

func (ac *AppCreator) createFileFromTemplate(filePath, templateContent string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("file").Parse(templateContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, data)
}

func (ac *AppCreator) ListApps() ([]string, error) {
	var gravApps []string

	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), "_grav") {
			gravApps = append(gravApps, entry.Name())
		}
	}

	return gravApps, nil
}

func (ac *AppCreator) DeleteApp(name string) error {
	appName := name + "_grav"
	err := os.RemoveAll(appName)
	if err != nil {
		return fmt.Errorf("failed to delete app directory %s: %w", appName, err)
	}
	ac.logger.Info("Grav app '" + appName + "' deleted successfully")
	return nil
}
