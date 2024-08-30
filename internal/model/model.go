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

// logger represents an instance of the logrus.Logger struct used for logging.
// It is used to log various events and errors in the ModelManager struct and its methods.
// Example usage can be found in the loadModels() method of the ModelManager struct.
var logger = logrus.New()

// Model represents a basic model structure for database entities.
// It includes the following fields:
//   - ID: The unique identifier for the model.
//   - CreatedAt: The timestamp of when the model was created.
//   - UpdatedAt: The timestamp of when the model was last updated.
type Model struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    // Ensure the field is exported
}

// ModelInterface is an interface that represents a model in a database.
// It defines methods for retrieving and manipulating data from the model's corresponding table.
//   - `TableName()` returns the name of the database table associated with the model.
//   - `PrimaryKey()` returns the name of the primary key column in the model's table.
//   - `BeforeCreate()` is called before creating a new record in the database.
//     It allows the model to perform any necessary operations or validations before the record is created.
//     It returns an error if any error occurs during the operation.
//   - `AfterCreate()` is called after a new record is created in the database.
//     It allows the model to perform any necessary operations or validations after the record is created.
//     It returns an error if any error occurs during the operation.
//   - `BeforeUpdate()` is called before updating an existing record in the database.
//     It allows the model to perform any necessary operations or validations before the record is updated.
//     It returns an error if any error occurs during the operation.
//   - `AfterUpdate()` is called after an existing record is updated in the database.
//     It allows the model to perform any necessary operations or validations after the record is updated.
//     It returns an error if any error occurs during the operation.
//   - `BeforeDelete()` is called before deleting a record from the database.
//     It allows the model to perform any necessary operations or validations before the record is deleted.
//     It returns an error if any error occurs during the operation.
//   - `AfterDelete()` is called after a record is deleted from the database.
//     It allows the model to perform any necessary operations or validations after the record is deleted.
//     It returns an error if any error occurs during the operation.
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

// DefaultModel represents a default implementation of a model that includes common fields like ID, CreatedAt, and UpdatedAt.
type DefaultModel struct {
	Model
}

// TableName returns the name of the table associated with the DefaultModel struct.
func (m *DefaultModel) TableName() string {
	return ""
}

// PrimaryKey returns the name of the primary key field for the DefaultModel.
func (m *DefaultModel) PrimaryKey() string {
	return "ID"
}

// BeforeCreate is a method that is called before a new instance of DefaultModel is created.
// This method is executed immediately before the model is saved to the database.
// It sets the CreatedAt and UpdatedAt fields of the model to the current time.
// The method signature should be: func (m *DefaultModel) BeforeCreate() error.
// This method does not return any error.
// The BeforeCreate method can be overridden in custom models to define custom behavior
// or perform any required actions before creating a new record.
func (m *DefaultModel) BeforeCreate() error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// AfterCreate is a method that is called after a new instance of DefaultModel is created.
// This method is executed immediately after the model is saved to the database.
// It can be overridden to perform any custom logic or actions after the model is created.
// The method signature should be: func (m *DefaultModel) AfterCreate() error.
// This method does not return any error.
func (m *DefaultModel) AfterCreate() error {
	return nil
}

// BeforeUpdate updates the 'UpdatedAt' field of the 'DefaultModel' instance
// with the current time.
// It is called automatically by the ORM before updating the model in the
// database.
// It returns an error if any error occurs during the update process.
func (m *DefaultModel) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}

// AfterUpdate is a method that is called after an instance of DefaultModel is updated in the database.
// This method can be overridden in custom models to define custom behavior or perform any required actions
// after the model is updated.
// The method signature should be: func (m *DefaultModel) AfterUpdate() error.
// This method does not return any error.
func (m *DefaultModel) AfterUpdate() error {
	return nil
}

// BeforeDelete is a method implemented by the DefaultModel struct.
// It is called before a record is deleted from the database.
// This method can be overridden in custom models to define custom behavior
// or perform any required actions before deleting a record.
// It returns an error if any errors occur during the deletion process.
func (m *DefaultModel) BeforeDelete() error {
	return nil
}

// AfterDelete is a method called after a record is deleted.
// It is used to perform any necessary cleanup or additional operations
// after the deletion of a record.
// This method can be overridden by embedding structs to provide
// custom behavior.
//
// The method returns an error if an error occurs during the cleanup
// or additional operations. If no error occurs, it returns nil.
// Any error returned will be handled by the calling code.
func (m *DefaultModel) AfterDelete() error {
	return nil
}

// Field represents a database field in a model.
type Field struct {
	Name      string
	Type      string
	Tag       string
	IsNull    bool
	IsPrimary bool
}

// NewField creates a new instance of the Field struct with the provided name, fieldType, tag,
// isNull, and isPrimary values. It returns the created Field.
func NewField(name, fieldType, tag string, isNull, isPrimary bool) Field {
	return Field{
		Name:      name,
		Type:      fieldType,
		Tag:       tag,
		IsNull:    isNull,
		IsPrimary: isPrimary,
	}
}

// ModelDefinition represents the definition of a model with its name, fields, and output directory.
type ModelDefinition struct {
	Name      string
	Fields    []Field
	OutputDir string
}

// NewModelDefinition creates a new instance of ModelDefinition with the specified name and fields.
// It returns a pointer to the newly created ModelDefinition.
func NewModelDefinition(name string, fields []Field) *ModelDefinition {
	return &ModelDefinition{
		Name:   name,
		Fields: fields,
	}
}

// ModelManager is responsible for managing model definitions. It provides functionalities to create, update, delete,
// retrieve, and list models. It also supports field validation and generating SQL migration scripts based on a model's
// definition. The manager uses a map to store the models, where the key is the model's name and the value is a pointer
// to a ModelDefinition struct. The manager can save and load models from a JSON file.
type ModelManager struct {
	models map[string]*ModelDefinition
}

// NewModelManager returns a new instance of ModelManager. It initializes the models map and loads the models from storage.
func NewModelManager() *ModelManager {
	mm := &ModelManager{
		models: make(map[string]*ModelDefinition),
	}
	mm.loadModels()
	return mm
}

// CreateModel creates a new model with the given name and fields. It checks if a model with the same name
// already exists and returns an error in that case. Otherwise, it creates a new model definition with the
// provided fields and adds it to the model manager's models map. It then saves the models to the storage file.
//
// Parameters:
// - name: The name of the model to create.
// - fields: The fields of the model as an array of Field structs.
//
// Returns:
// - error: An error if the model already exists or there is an error saving the models to the storage file.
func (mm *ModelManager) CreateModel(name string, fields []Field) error {
	if _, exists := mm.models[name]; exists {
		return fmt.Errorf("model %s already exists", name)
	}

	mm.models[name] = NewModelDefinition(name, fields)
	return mm.saveModels()
}

// UpdateModel updates the fields of an existing model. It first checks if the model exists in the model manager's
// models map. If the model does not exist, an error is returned. Otherwise, the model's fields are updated with the
// provided fields.
func (mm *ModelManager) UpdateModel(name string, fields []Field) error {
	if _, exists := mm.models[name]; !exists {
		return fmt.Errorf("model %s does not exist", name)
	}

	mm.models[name] = NewModelDefinition(name, fields)
	return nil
}

// DeleteModel deletes a model from the ModelManager's models collection.
// It takes the name of the model to be deleted as a parameter.
// If the model does not exist in the collection, it returns an error.
// Otherwise, the model is deleted from the collection.
// It returns nil if the deletion is successful.
func (mm *ModelManager) DeleteModel(name string) error {
	if _, exists := mm.models[name]; !exists {
		return fmt.Errorf("model %s does not exist", name)
	}

	delete(mm.models, name)
	return nil
}

// GetModel retrieves a model definition by name from the ModelManager. It returns the model definition
// and an error if the model does not exist.
func (mm *ModelManager) GetModel(name string) (*ModelDefinition, error) {
	model, exists := mm.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s does not exist", name)
	}

	return model, nil
}

// ListModels returns a sorted list of model names in the ModelManager.
func (mm *ModelManager) ListModels() []string {
	var modelNames []string
	for name := range mm.models {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)
	return modelNames
}

// ValidateField validates the type of a field.
// It checks if the field type is one of the valid types: string, int, bool, time.Time, float64, []byte.
// If the field type is not valid, it returns an error indicating the invalid field type.
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

// GenerateMigration generates a SQL migration statement for creating a table based on a given ModelDefinition.
// The generated migration includes the table name, field names, data types, and any additional constraints (e.g., primary key, not null).
// The resulting migration statement is returned as a string.
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

// getSQLType returns the SQL data type corresponding to a given Go type. It maps the following Go types to their SQL equivalents:
// - string: VARCHAR(255)
// - int: INTEGER
// - bool: BOOLEAN
// - time.Time: TIMESTAMP
// - float64: DOUBLE PRECISION
// - []byte: BYTEA
// If the given Go type does not match any of the above, it returns "VARCHAR(255)" as the default SQL type.
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

// modelStorageFile is the file name of the JSON file used to store the models.
const modelStorageFile = "models.json"

// saveModels saves the models in the ModelManager to a JSON file.
// It marshals the models into JSON format and writes the data to the
// specified modelStorageFile. The file is created with read and write
// permissions for the owner and readable by others. If there is an
// error during the process, it is returned as an error.
//
// This method is called by CreateModel after adding a new model, UpdateModel
// after updating a model, and DeleteModel after deleting a model.
//
// Example usage:
//
//	mm := &ModelManager{}
//	// ... code to populate mm.models ...
//	err := mm.saveModels()
//	if err != nil {
//	  fmt.Println("Failed to save models:", err)
//	}
//
// Note: This method is not intended to be called directly by users of
// this package.
func (mm *ModelManager) saveModels() error {
	data, err := json.Marshal(mm.models)
	if err != nil {
		return err
	}
	return os.WriteFile(modelStorageFile, data, 0644)
}

// loadModels reads the content of the models file, if it exists, and
// unmarshals the data into the ModelManager's models map. If the file
// does not exist, it logs a message and returns. If there is an error
// while reading or unmarshaling the data, it logs an error message.
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

// SetOutputDir sets the output directory for the ModelDefinition.
func (m *ModelDefinition) SetOutputDir(dir string) {
	m.OutputDir = dir
}
