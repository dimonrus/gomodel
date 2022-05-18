package gomodel

import (
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
	"github.com/dimonrus/porterr"
	"reflect"
)

const (
	// DefaultSchema default database schema
	DefaultSchema = "public"
)

// IModel DB model interface
type IModel interface {
	// Table Returns table name
	Table() string
	// Columns returns all columns
	Columns() []string
	// Values returns all model values
	Values() []any
}

// GetColumn model column in table
func GetColumn(model IModel, field any) string {
	columns := GetColumns(model, field)
	if len(columns) > 0 {
		return columns[0]
	}
	return ""
}

// GetColumns model columns by fields
func GetColumns(model IModel, field ...any) []string {
	if model == nil {
		return nil
	}
	ve := reflect.ValueOf(model).Elem()
	te := reflect.TypeOf(model).Elem()
	columns := make([]string, len(field))
	var k int
	var tField ModelFiledTag
	for j := range field {
		cte := reflect.ValueOf(field[j])
		if cte.Kind() != reflect.Ptr {
			continue
		}
		for i := 0; i < ve.NumField(); i++ {
			if ve.Field(i).Addr().Pointer() == cte.Elem().Addr().Pointer() {
				tField = ParseModelFiledTag(te.Field(i).Tag.Get("db"))
				columns[k] = tField.Column
				k++
			}
		}
	}
	return columns[:k]
}

// GetValues model values by columns
func GetValues(model IModel, columns ...string) (values []any) {
	if model == nil {
		return nil
	}
	te := reflect.TypeOf(model).Elem()
	modelValues := model.Values()
	values = make([]any, len(columns))
	var j int
	for i := 0; i < len(modelValues); i++ {
		tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
		if gohelp.ExistsInArray(tField.Column, columns) {
			values[j] = modelValues[i]
			j++
		}
	}
	values = values[:j]
	return
}

// GetUpdate model update query
// model - target model
// condition - one or more where condition
// fields - list of fields that you want to update
func GetUpdate(model IModel, condition *gosql.Condition, fields ...any) (update *gosql.Update, e porterr.IError) {
	if model == nil {
		e = porterr.New(porterr.PortErrorArgument, "model is nil, check your logic")
		return
	}
	if fields == nil {
		fields = model.Values()
		return
	}
	ve := reflect.ValueOf(model).Elem()
	te := reflect.TypeOf(model).Elem()
	update = gosql.NewUpdate()
	for i := 0; i < ve.NumField(); i++ {
		for _, v := range fields {
			cte := reflect.ValueOf(v)
			if cte.Kind() != reflect.Ptr {
				e = porterr.New(porterr.PortErrorArgument, "Fields must be an interfaces")
				return
			}
			if ve.Field(i).Addr().Pointer() == cte.Elem().Addr().Pointer() {
				tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
				update.Set().Append(tField.Column+" = ?", v)
			}
		}
	}
	if update.IsEmpty() {
		e = porterr.New(porterr.PortErrorArgument, "no columns found in model")
		return
	}
	update.Table(model.Table())
	if !condition.IsEmpty() {
		update.Where().Replace(condition)
	}
	return update, e
}

// GetDelete model delete query
func GetDelete(model IModel, condition *gosql.Condition) (delete *gosql.Delete, e porterr.IError) {
	if model == nil {
		e = porterr.New(porterr.PortErrorArgument, "model is nil, check your logic")
		return
	}
	delete = gosql.NewDelete().From(model.Table())
	if !condition.IsEmpty() {
		delete.Where().Replace(condition)
	}
	return delete, e
}
