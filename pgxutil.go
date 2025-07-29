package dbx

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	_pgx "github.com/jackc/pgx/v5"
)

// Get selects a single row and scans it into a struct
func Get[T any](ctx context.Context, sql string, args ...any) (*T, error) {
	rows, err := Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, _pgx.ErrNoRows
	}
	
	result := new(T)
	err = scanRowToStruct(rows, result)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

// Select selects multiple rows and scans them into a slice of structs
func Select[T any](ctx context.Context, sql string, args ...any) ([]*T, error) {
	rows, err := Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []*T
	
	for rows.Next() {
		item := new(T)
		err = scanRowToStruct(rows, item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return results, nil
}

// scanRowToStruct scans a row into a struct, ignoring columns that don't have corresponding struct fields
func scanRowToStruct(rows _pgx.Rows, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}
	
	structValue := v.Elem()
	structType := structValue.Type()
	
	// Build a map of field names to struct field info
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		
		// Get column name from db tag or field name
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = strings.ToLower(field.Name)
		}
		
		// Skip if tag is "-"
		if columnName == "-" {
			continue
		}
		
		if fieldValue.CanSet() {
			fieldMap[columnName] = fieldValue
		}
	}
	
	// Get column descriptions from the result
	fieldDescriptions := rows.FieldDescriptions()
	scanTargets := make([]any, len(fieldDescriptions))
	
	for i, col := range fieldDescriptions {
		columnName := string(col.Name)
		if fieldValue, exists := fieldMap[columnName]; exists {
			// Field exists in struct - scan into it
			scanTargets[i] = fieldValue.Addr().Interface()
		} else {
			// Field doesn't exist in struct - scan into a dummy variable
			var dummy any
			scanTargets[i] = &dummy
		}
	}
	
	return rows.Scan(scanTargets...)
}

// InsertStruct inserts a struct into the specified table
// Always returns the inserted row using RETURNING *
func InsertStruct[T any](ctx context.Context, tableName string, data T) (*T, error) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	
	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	
	var columns []string
	var placeholders []string
	var values []any
	
	// Build columns and values from struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		
		// Get column name from db tag or field name
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = strings.ToLower(field.Name)
		}
		
		// Skip if tag is "-"
		if columnName == "-" {
			continue
		}
		
		columns = append(columns, columnName)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
		values = append(values, fieldValue.Interface())
	}
	
	// Build the INSERT query with RETURNING
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	
	// Execute with RETURNING
	rows, err := Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	
	// Collect the returned row
	result, err := _pgx.CollectOneRow(rows, _pgx.RowToAddrOfStructByName[T])
	if err != nil {
		return nil, err
	}
	
	return result, nil
}
