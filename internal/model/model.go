package model

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
	"time"
)

var logger = logrus.New()

// Model represents a basic model structure with common fields
type Model struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ModelInterface defines the methods that should be implemented by models
type ModelInterface interface {
	TableName() string
	PrimaryKey() string
	BeforeCreate() error
	AfterCreate() error
	BeforeUpdate() error
	AfterUpdate() error
	BeforeDelete() error
	AfterDelete() error
}

// DefaultModel provides a default implementation of ModelInterface
type DefaultModel struct {
	Model
}

// TableName returns the default table name for the model
func (m *DefaultModel) TableName() string {
	return ""
}

// PrimaryKey returns the primary key field name
func (m *DefaultModel) PrimaryKey() string {
	return "ID"
}

// BeforeCreate is called before creating a new record
func (m *DefaultModel) BeforeCreate() error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// AfterCreate is called after creating a new record
func (m *DefaultModel) AfterCreate() error {
	return nil
}

// BeforeUpdate is called before updating a record
func (m *DefaultModel) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}

// AfterUpdate is called after updating a record
func (m *DefaultModel) AfterUpdate() error {
	return nil
}

// BeforeDelete is called before deleting a record
func (m *DefaultModel) BeforeDelete() error {
	return nil
}

// AfterDelete is called after deleting a record
func (m *DefaultModel) AfterDelete() error {
	return nil
}

// Field represents a database field
type Field struct {
	Name      string
	Type      string
	Tag       string
	IsNull    bool
	IsPrimary bool
}

// NewField creates a new Field instance
func NewField(name, fieldType, tag string, isNull, isPrimary bool) Field {
	return Field{
		Name:      name,
		Type:      fieldType,
		Tag:       tag,
		IsNull:    isNull,
		IsPrimary: isPrimary,
	}
}

// ModelDefinition represents the structure of a model
type ModelDefinition struct {
	Name      string
	Fields    []Field
	OutputDir string
}

// NewModelDefinition creates a new ModelDefinition instance
func NewModelDefinition(name string, fields []Field) *ModelDefinition {
	return &ModelDefinition{
		Name:   name,
		Fields: fields,
	}
}

// ModelManager handles model-related operations
type ModelManager struct {
	models map[string]*ModelDefinition
}

// NewModelManager creates a new ModelManager instance
func NewModelManager() *ModelManager {
	mm := &ModelManager{
		models: make(map[string]*ModelDefinition),
	}
	mm.loadModels()
	return mm
}

// CreateModel creates a new model
func (mm *ModelManager) CreateModel(name string, fields []Field) error {
	if _, exists := mm.models[name]; exists {
		return fmt.Errorf("model %s already exists", name)
	}

	mm.models[name] = NewModelDefinition(name, fields)
	return mm.saveModels()
}

// UpdateModel updates an existing model
func (mm *ModelManager) UpdateModel(name string, fields []Field) error {
	if _, exists := mm.models[name]; !exists {
		return fmt.Errorf("model %s does not exist", name)
	}

	mm.models[name] = NewModelDefinition(name, fields)
	return nil
}

// DeleteModel deletes a model
func (mm *ModelManager) DeleteModel(name string) error {
	if _, exists := mm.models[name]; !exists {
		return fmt.Errorf("model %s does not exist", name)
	}

	delete(mm.models, name)
	return nil
}

// GetModel retrieves a model by name
func (mm *ModelManager) GetModel(name string) (*ModelDefinition, error) {
	model, exists := mm.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s does not exist", name)
	}

	return model, nil
}

// ListModels returns a list of all model names
func (mm *ModelManager) ListModels() []string {
	var modelNames []string
	for name := range mm.models {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)
	return modelNames
}

// ValidateField validates a field's type
func (mm *ModelManager) ValidateField(field Field) error {
	validTypes := map[string]bool{
		"string": true, "int": true, "bool": true, "time.Time": true,
		"float64": true, "[]byte": true,
	}

	if !validTypes[field.Type] {
		return fmt.Errorf("invalid field type: %s", field.Type)
	}

	return nil
}

// GenerateMigration generates a migration for a model
func (mm *ModelManager) GenerateMigration(model *ModelDefinition) string {
	var migration strings.Builder

	migration.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", strings.ToLower(model.Name)))

	for _, field := range model.Fields {
		migration.WriteString(fmt.Sprintf("  %s %s", strings.ToLower(field.Name), getSQLType(field.Type)))
		if field.IsPrimary {
			migration.WriteString(" PRIMARY KEY")
		}
		if !field.IsNull {
			migration.WriteString(" NOT NULL")
		}
		migration.WriteString(",\n")
	}

	migration.WriteString(");\n")

	return migration.String()
}

// getSQLType converts a Go type to a SQL type
func getSQLType(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "int":
		return "INTEGER"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMP"
	case "float64":
		return "DOUBLE PRECISION"
	case "[]byte":
		return "BYTEA"
	default:
		return "VARCHAR(255)"
	}
}

const modelStorageFile = "models.json"

func (mm *ModelManager) saveModels() error {
	data, err := json.Marshal(mm.models)
	if err != nil {
		return err
	}
	return os.WriteFile(modelStorageFile, data, 0644)
}

func (mm *ModelManager) loadModels() {
	data, err := os.ReadFile(modelStorageFile)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.WithError(err).Error("Failed to read models file")
		}
		return
	}

	err = json.Unmarshal(data, &mm.models)
	if err != nil {
		logger.WithError(err).Error("Failed to unmarshal models")
	}
}

// Add this method to your ModelDefinition struct
func (m *ModelDefinition) SetOutputDir(dir string) {
	m.OutputDir = dir
}
