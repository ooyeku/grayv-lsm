package model

import (
	"time"
)

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
	Name   string
	Fields []Field
}

// NewModelDefinition creates a new ModelDefinition instance
func NewModelDefinition(name string, fields []Field) *ModelDefinition {
	return &ModelDefinition{
		Name:   name,
		Fields: fields,
	}
}
