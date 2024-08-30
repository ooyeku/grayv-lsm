package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// modelTemplate is a constant that holds the template for generating a model file based on a `ModelDefinition`.
// The template includes the necessary import statements and defines the struct fields using the provided `ModelDefinition` fields.
// The `{{.Name}}` placeholder is replaced with the name of the model. The field names are transformed to title case using the `title` function.
// The `json` struct tag is generated using the field name transformed to lowercase.
// The `TableName` method is defined to return the lowercase plural form of the model name followed by "s".
const modelTemplate = `package models

import (
	"github.com/ooyeku/grav-lsm/internal/model"
)

type {{.Name}} struct {
	model.DefaultModel
	{{- range .Fields}}
	{{.Name | title}} {{.Type}} ` + "`json:\"{{.Name | toLower}}\"`" + `
	{{- end}}
}

func ({{.Name | firstLetter}} *{{.Name}}) TableName() string {
	return "{{.Name | toLower}}s"
}
`

// GenerateModelFile generates a model file based on the provided model definition.
// The function uses a template to define the structure and fields of the model.
// The template includes necessary import statements and generates the necessary struct tags for JSON serialization.
// The generated model file is saved in the specified output directory, or in the default "models" directory if no output directory is provided.
// Returns an error if there is any issue parsing the template, creating the output directory, creating the file, executing the template, or any other related error.
func GenerateModelFile(modelDef *ModelDefinition) error {
	tmpl, err := template.New("model").Funcs(template.FuncMap{
		"toLower": strings.ToLower,
		"firstLetter": func(s string) string {
			return strings.ToLower(s[:1])
		},
		"title": strings.Title,
	}).Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	outputDir := modelDef.OutputDir
	if outputDir == "" {
		outputDir = "models"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	fileName := filepath.Join(outputDir, strings.ToLower(modelDef.Name)+".go")
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, modelDef); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// LoadModelDefinition loads the definition of a model with the given name. It returns
// a pointer to a ModelDefinition struct and an error. The function currently has a placeholder
// implementation and returns a ModelDefinition with the provided modelName and an empty Fields slice.
// In a real-world scenario, you would parse an existing model file and populate the ModelDefinition
// struct accordingly.
func LoadModelDefinition(modelName string) (*ModelDefinition, error) {
	// This is a placeholder implementation. In a real-world scenario,
	// you would parse the existing model file and create a ModelDefinition from it.
	return &ModelDefinition{
		Name:   modelName,
		Fields: []Field{},
	}, nil
}
