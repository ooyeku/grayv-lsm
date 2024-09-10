package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ooyeku/grayv-lsm/pkg/logging"
)

// AppCreator is a type that represents an application creator. It has a logger property of type *logging.ColorfulLogger.
type AppCreator struct {
	logger *logging.ColorfulLogger
}

// NewAppCreator is a function that creates and returns a new instance of the AppCreator struct.
// It initializes the logger of the AppCreator with a new instance of the ColorfulLogger struct.
// This function does not take any arguments.
// It returns a pointer to the created AppCreator instance.
// Example usage:
//
//	appCreator := NewAppCreator()
//	appName := "myapp"
//	err := appCreator.CreateApp(appName)
//	if err != nil {
//	    // handle error
//	}
//	gravApps, err := appCreator.ListApps()
//	if err != nil {
//	    // handle error
//	}
//	for _, app := range gravApps {
//	    fmt.Println(app)
//	}
//	err := appCreator.DeleteApp(appName)
//	if err != nil {
//	    // handle error
//	}
func NewAppCreator() *AppCreator {
	return &AppCreator{logger: logging.NewColorfulLogger()}
}

// CreateApp creates a new Grav app with the specified name. It appends "_grav" to the app name,
// creates the main app directory, and creates several subdirectories. It also creates a main.go file
// and initializes a Go module for the app. The app name and other relevant information are logged.
// If any step fails, an error is returned.
//
// Parameters:
// - name: the name of the app to be created.
//
// Returns:
// - error: an error if the app creation fails.
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

// createMainFile creates the main.go file for the Grav app.
func (ac *AppCreator) createMainFile(appName string) error {
	mainTemplate := `package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to %s!", appName)
    })

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
`
	return ac.createFileFromTemplate(filepath.Join(appName, "cmd", "main.go"), mainTemplate, appName)
}

// createGoMod initializes a new Go module for the specified app name.
// It executes the `go mod init` command in the directory of the app,
// sets the app name as the module name, and creates the go.mod file.
// It returns an error if the initialization fails along with any output from the command.
// It logs a message if the Go module is successfully initialized.
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

// createFileFromTemplate creates a new file at the given filePath using the provided templateContent and data.
// It returns an error if file creation or template parsing fails.
// This method is used by the AppCreator to generate specific files for an app.
//
// Parameters:
// - filePath: the path where the file should be created.
// - templateContent: the content of the template to be used for file generation.
// - data: the data to be passed to the template for rendering.
//
// Returns:
// - error: an error if file creation or template parsing fails.
func (ac *AppCreator) createFileFromTemplate(filePath, templateContent string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ac.logger.Error("failed to close file: " + err.Error())
		}
	}(file)

	tmpl, err := template.New("file").Parse(templateContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, data)
}

// ListApps returns a list of Grav apps in the current directory. It searches for directories
// that have names ending with "_grav". The method returns the list of app names and an error
// if there was an issue reading the directory.
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

// DeleteApp deletes the Grav app with the specified name. It first appends "_grav" to the app name
// to form the directory name. Then, it removes the entire app directory using the os.RemoveAll function.
// If the deletion fails, an error is returned along with a descriptive error message. Successful deletion
// is logged using the logger's Info method.
func (ac *AppCreator) DeleteApp(name string) error {
	appName := name + "_grav"
	err := os.RemoveAll(appName)
	if err != nil {
		return fmt.Errorf("failed to delete app directory %s: %w", appName, err)
	}
	ac.logger.Info("Grav app '" + appName + "' deleted successfully")
	return nil
}
