package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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

	modelsDir := "models"
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("error creating models directory: %w", err)
	}

	fileName := filepath.Join(modelsDir, strings.ToLower(modelDef.Name)+".go")
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

func LoadModelDefinition(modelName string) (*ModelDefinition, error) {
	// This is a placeholder implementation. In a real-world scenario,
	// you would parse the existing model file and create a ModelDefinition from it.
	return &ModelDefinition{
		Name:   modelName,
		Fields: []Field{},
	}, nil
}
