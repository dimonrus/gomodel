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

// MetaModel Meta model contain full information about model and fields
type MetaModel struct {
	// Table name
	TableName string
	// Fields
	Fields ModelFiledTagList
}

// PrepareMetaModel Prepare Meta Model definition
func PrepareMetaModel(model IModel) *MetaModel {
	if model == nil {
		return nil
	}
	var k int
	ve := reflect.ValueOf(model).Elem()
	te := reflect.TypeOf(model).Elem()
	meta := MetaModel{
		TableName: gohelp.ToUnderscore(model.Table()),
		Fields:    make([]ModelFiledTag, ve.NumField()),
	}
	for i := 0; i < ve.NumField(); i++ {
		tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
		if !tField.IsIgnored {
			tField.Value = ve.Field(i).UnsafePointer()
			meta.Fields[k] = tField
			k++
		}
	}
	meta.Fields = meta.Fields[:k]
	return &meta
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
	var k int
	var tField ModelFiledTag
	if len(field) == 0 {
		field = model.Values()
	}
	columns := make([]string, len(field))
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

// extract model. Get name, columns, values
func extract(model IModel) (table string, columns []string, values []any) {
	if model != nil {
		ve := reflect.ValueOf(model).Elem()
		te := reflect.TypeOf(model).Elem()
		table = gohelp.ToUnderscore(te.Name())
		columns = make([]string, ve.NumField())
		values = make([]any, ve.NumField())
		var k int
		for i := 0; i < ve.NumField(); i++ {
			tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
			if !tField.IsIgnored {
				columns[k] = tField.Column
				values[k] = ve.Interface()
				k++
			}
		}
	}
	return
}

// GetValues model values by columns
func GetValues(model IModel, columns ...string) (values []any) {
	if model == nil {
		return
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
