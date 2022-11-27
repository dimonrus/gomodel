package gomodel

import (
	"database/sql"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
	"github.com/dimonrus/porterr"
	"net/http"
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
	ve := reflect.ValueOf(model)
	if ve.IsNil() {
		return nil
	}
	ve = ve.Elem()
	te := reflect.TypeOf(model).Elem()
	var k int
	meta := &MetaModel{
		TableName: model.Table(),
		Fields:    make([]ModelFiledTag, ve.NumField()),
	}
	for i := 0; i < ve.NumField(); i++ {
		ParseModelFiledTag(te.Field(i).Tag.Get("db"), &meta.Fields[i])
		if !meta.Fields[i].IsIgnored {
			meta.Fields[i].Value = ve.Field(i).Addr().Interface()
			k++
		}
	}
	meta.Fields = meta.Fields[:k]
	return meta
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
			tField.Clear()
			if ve.Field(i).Addr().Pointer() == cte.Elem().Addr().Pointer() {
				ParseModelFiledTag(te.Field(i).Tag.Get("db"), &tField)
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
		var tField ModelFiledTag
		for i := 0; i < ve.NumField(); i++ {
			tField.Clear()
			ParseModelFiledTag(te.Field(i).Tag.Get("db"), &tField)
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
	var tField ModelFiledTag
	for i := 0; i < len(modelValues); i++ {
		tField.Clear()
		ParseModelFiledTag(te.Field(i).Tag.Get("db"), &tField)
		if gohelp.ExistsInArray(tField.Column, columns) {
			values[j] = modelValues[i]
			j++
		}
	}
	values = values[:j]
	return
}

// Do exec query on model
func Do(q godb.Queryer, isql gosql.ISQL) (e porterr.IError) {
	if isql == nil {
		e = porterr.New(porterr.PortErrorLoad, "ISQL is empty. Check your logic")
		return
	}
	var err error
	query, params, returning := isql.SQL()
	if len(returning) > 0 {
		err = q.QueryRow(query, params...).Scan(returning...)
	} else {
		_, err = q.Exec(query, params...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			e = porterr.New(porterr.PortErrorSearch, "No record found. Check params or model already deleted").HTTP(http.StatusNotFound)
		} else {
			e = porterr.New(porterr.PortErrorIO, err.Error())
		}
	}
	return
}

// Load get isql and load model
func Load(q godb.Queryer, model IModel) porterr.IError {
	return Do(q, GetLoadSQL(model))
}

// Save get isql and save model
func Save(q godb.Queryer, model IModel) porterr.IError {
	return Do(q, GetSaveSQL(model))
}

// Delete get isql and delete model
func Delete(q godb.Queryer, model IModel) porterr.IError {
	return Do(q, GetDeleteSQL(model))
}
