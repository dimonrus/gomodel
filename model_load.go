package gomodel

import (
	"github.com/dimonrus/gosql"
	"reflect"
)

// GetLoadSQL return sql query fot load model
func GetLoadSQL(model IModel) gosql.ISQL {
	if model == nil {
		return nil
	}
	ve := reflect.ValueOf(model).Elem()
	te := reflect.TypeOf(model).Elem()
	selectSql := gosql.NewSelect()
	selectSql.From(model.Table())
	cond := gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
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
			cond.AddExpression(tField.Column + " IS NOT NULL")
		}
		selectSql.Columns().Append(tField.Column, field.Interface())
	}
	if !cond.IsEmpty() {
		selectSql.Where().Replace(cond)
	}
	return selectSql
}
