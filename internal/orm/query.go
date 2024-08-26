package orm

import (
	"fmt"
	"strings"
)

// Query represents a database query
type Query struct {
	table     string
	operation string
	fields    []string
	where     []string
	params    []interface{}
	limit     int
	offset    int
}

// NewQuery creates a new Query instance
func NewQuery(table string) *Query {
	return &Query{
		table:  table,
		fields: []string{"*"},
	}
}

// Select specifies the fields to select
func (q *Query) Select(fields ...string) *Query {
	q.operation = "SELECT"
	q.fields = fields
	return q
}

// Where adds a WHERE condition
func (q *Query) Where(condition string, params ...interface{}) *Query {
	q.where = append(q.where, condition)
	q.params = append(q.params, params...)
	return q
}

// Limit sets the LIMIT clause
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Offset sets the OFFSET clause
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// Insert prepares an INSERT query
func (q *Query) Insert(fields ...string) *Query {
	q.operation = "INSERT"
	q.fields = fields
	return q
}

// Update prepares an UPDATE query
func (q *Query) Update(fields ...string) *Query {
	q.operation = "UPDATE"
	q.fields = fields
	return q
}

// Delete prepares a DELETE query
func (q *Query) Delete() *Query {
	q.operation = "DELETE"
	return q
}

// Build constructs the SQL query
func (q *Query) Build() (string, []interface{}) {
	var query strings.Builder
	var params []interface{}

	switch q.operation {
	case "SELECT":
		query.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.fields, ", "), q.table))
	case "INSERT":
		placeholders := make([]string, len(q.fields))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			q.table, strings.Join(q.fields, ", "), strings.Join(placeholders, ", ")))
	case "UPDATE":
		query.WriteString(fmt.Sprintf("UPDATE %s SET ", q.table))
		for i, field := range q.fields {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString(fmt.Sprintf("%s = ?", field))
		}
	case "DELETE":
		query.WriteString(fmt.Sprintf("DELETE FROM %s", q.table))
	}

	if len(q.where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(q.where, " AND "))
		params = append(params, q.params...)
	}

	if q.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", q.limit))
	}

	if q.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", q.offset))
	}

	return query.String(), params
}
