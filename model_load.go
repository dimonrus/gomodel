package gomodel

import (
	"github.com/dimonrus/gosql"
	"reflect"
)

// GetLoadSQL return sql query fot load model
func GetLoadSQL(model IModel) gosql.ISQL {
	ve := reflect.ValueOf(model)
	te := reflect.TypeOf(model).Elem()
	if ve.IsNil() {
		return nil
	}
	ve = ve.Elem()
	selectSql := gosql.NewSelect()
	selectSql.From(model.Table())
	cond := gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	var tField ModelFiledTag
	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		tField.Clear()
		ParseModelFiledTag(te.Field(i).Tag.Get("db"), &tField)
		if tField.IsIgnored || tField.Column == "" {
			continue
		}
		if tField.IsPrimaryKey && !field.IsNil() {
			cond.AddExpression(tField.Column+" = ?", field.Interface())
		} else if tField.IsUnique && !field.IsNil() {
			if cond.IsEmpty() {
				cond.AddExpression(tField.Column+" = ?", field.Interface())
			}
		} else if tField.IsDeletedAt {
			cond.AddExpression(tField.Column + " IS NULL")
		}
		selectSql.Columns().Append(tField.Column, field.Addr().Interface())
	}
	if !cond.IsEmpty() {
		selectSql.Where().Replace(cond)
	}
	return selectSql
}
