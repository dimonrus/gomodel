package gomodel

import (
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
	"reflect"
)

// GetUpdateSQL model update query
// model - target model
// fields - list of fields that you want to update
func GetUpdateSQL(model IModel, fields ...any) gosql.ISQL {
	var ve = reflect.ValueOf(model)
	var te = reflect.TypeOf(model).Elem()
	if ve.IsNil() {
		return nil
	}
	ve = ve.Elem()
	if fields == nil {
		fields = model.Values()
	}
	var hasPrimaryKey bool
	var condition = gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	var update = gosql.NewUpdate()
	var tField ModelFiledTag
	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		for _, v := range fields {
			tField.Clear()
			cte := reflect.ValueOf(v)
			if cte.Kind() != reflect.Ptr {
				return nil
			}
			if ve.Field(i).Addr().Pointer() == cte.Elem().Addr().Pointer() {
				ParseModelFiledTag(te.Field(i).Tag.Get("db"), &tField)
				if tField.IsPrimaryKey {
					hasPrimaryKey = true
					if !field.IsNil() {
						if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
							condition.AddExpression(tField.Column+" = ?", pq.Array(field.Interface()))
						} else {
							condition.AddExpression(tField.Column+" = ?", field.Interface())
						}
					}
				} else if tField.IsUnique && !hasPrimaryKey {
					if !field.IsNil() {
						if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
							condition.AddExpression(tField.Column+" = ?", pq.Array(field.Interface()))
						} else {
							condition.AddExpression(tField.Column+" = ?", field.Interface())
						}
					}
				} else if !tField.IsIgnored {
					if tField.IsCreatedAt {
						update.Returning().Append(tField.Column, field.Addr().Interface())
					} else if tField.IsUpdatedAt {
						if !field.IsNil() {
							update.Set().Append(tField.Column+" = ?", field.Interface())
						} else {
							update.Set().Append(tField.Column + " = NOW()")
						}
						update.Returning().Append(tField.Column, field.Addr().Interface())
					} else if tField.IsDeletedAt {
						update.Returning().Append(tField.Column, field.Addr().Interface())
					} else if tField.IsSequence {
						if !field.IsNil() {
							if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
								condition.AddExpression(tField.Column+" = ?", pq.Array(field.Interface()))
							} else {
								condition.AddExpression(tField.Column+" = ?", field.Interface())
							}
						}
						update.Returning().Append(tField.Column, field.Addr().Interface())
					} else {
						if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
							update.Set().Append(tField.Column+" = ?", pq.Array(field.Interface()))
						} else {
							update.Set().Append(tField.Column+" = ?", field.Interface())
						}
					}
				}
			}
		}
	}
	if update.IsEmpty() || condition.IsEmpty() {
		return nil
	}
	update.Table(model.Table())
	update.Where().Replace(condition)
	return update
}
