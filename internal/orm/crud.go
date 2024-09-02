package orm

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/ooyeku/grayv-lsm/internal/model"
)

// CRUD provides basic CRUD operations for models
type CRUD struct {
	conn *Connection
}

// NewCRUD creates a new CRUD instance
func NewCRUD(conn *Connection) *CRUD {
	return &CRUD{conn: conn}
}

// Create inserts a new record into the database
func (c *CRUD) Create(m model.ModelInterface) error {
	v := reflect.ValueOf(m).Elem()
	t := v.Type()

	var fields []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "Model" {
			fields = append(fields, field.Name)
			values = append(values, v.Field(i).Interface())
		}
	}

	q := NewQuery(m.TableName()).Insert(fields...)
	query, _ := q.Build()

	_, err := c.conn.db.Exec(query, values...)
	return err
}

// Read retrieves a record from the database
func (c *CRUD) Read(m model.ModelInterface, id interface{}) error {
	q := NewQuery(m.TableName()).Where(fmt.Sprintf("%s = ?", m.PrimaryKey()), id)
	query, params := q.Build()

	row := c.conn.db.QueryRow(query, params...)

	v := reflect.ValueOf(m).Elem()
	fields := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		fields[i] = v.Field(i).Addr().Interface()
	}

	return row.Scan(fields...)
}

// Update updates a record in the database
func (c *CRUD) Update(m model.ModelInterface) error {
	v := reflect.ValueOf(m).Elem()
	t := v.Type()

	var fields []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "Model" && field.Name != m.PrimaryKey() {
			fields = append(fields, field.Name)
			values = append(values, v.Field(i).Interface())
		}
	}

	id := v.FieldByName(m.PrimaryKey()).Interface()
	q := NewQuery(m.TableName()).Update(fields...).Where(fmt.Sprintf("%s = ?", m.PrimaryKey()), id)
	query, _ := q.Build()

	values = append(values, id)
	_, err := c.conn.db.Exec(query, values...)
	return err
}

// Delete removes a record from the database
func (c *CRUD) Delete(m model.ModelInterface, id interface{}) error {
	q := NewQuery(m.TableName()).Delete().Where(fmt.Sprintf("%s = ?", m.PrimaryKey()), id)
	query, params := q.Build()

	_, err := c.conn.db.Exec(query, params...)
	return err
}

// Query executes a custom query and returns the rows
func (c *CRUD) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.conn.db.Query(query, args...)
}

// Exec executes a custom query without returning any rows
func (c *CRUD) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.conn.db.Exec(query, args...)
}
