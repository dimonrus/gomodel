package gomodel

import (
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
	"reflect"
)

// GetSaveSQL prepare save query
// it can be insert or update or upsert
// some popular scenario was implemented. not all
func GetSaveSQL(model IModel) gosql.ISQL {
	if model == nil {
		return nil
	}
	ve := reflect.ValueOf(model).Elem()
	te := reflect.TypeOf(model).Elem()
	var columnsInsert = gosql.NewExpression()
	var columnsUpdate = gosql.NewExpression()
	var conflictObject = gosql.NewExpression()
	var returning = gosql.NewExpression()
	var condition = gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	var conflict = gosql.NewConflict().Action(gosql.ConflictActionUpdate)
	var insert, update, upsert, hasPrimaryKey bool
	var result gosql.ISQL

	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
		if tField.IsPrimaryKey {
			hasPrimaryKey = true
			if !field.IsNil() {
				if tField.IsSequence {
					update = true
					condition.AddExpression(tField.Column+" = ?", field.Interface())
				} else {
					upsert = true
					if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
						columnsInsert.Append(tField.Column, pq.Array(field.Interface()))
					} else {
						columnsInsert.Append(tField.Column, field.Interface())
					}
					conflictObject.Add(tField.Column)
				}
			} else {
				if tField.IsSequence {
					insert = true
					returning.Append(tField.Column, field.Addr().Interface())
				} else {
					upsert = true
					if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
						columnsInsert.Append(tField.Column, pq.Array(field.Interface()))
					} else {
						columnsInsert.Append(tField.Column, field.Interface())
					}
					conflictObject.Add(tField.Column)
				}
			}
		} else if tField.IsUnique && !hasPrimaryKey {
			if !field.IsNil() {
				if tField.IsSequence {
					update = true
					condition.AddExpression(tField.Column+" = ?", field.Interface())
				} else if !insert && !update {
					upsert = true
					if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
						columnsInsert.Append(tField.Column, pq.Array(field.Interface()))
					} else {
						columnsInsert.Append(tField.Column, field.Interface())
					}
					conflictObject.Add(tField.Column)
				}
			} else {
				if tField.IsSequence {
					insert = true
					returning.Append(tField.Column, field.Addr().Interface())
				} else if !insert && !update {
					upsert = true
					if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
						columnsInsert.Append(tField.Column, pq.Array(field.Interface()))
					} else {
						columnsInsert.Append(tField.Column, field.Interface())
					}
					conflictObject.Add(tField.Column)
				}
			}
		} else if !tField.IsIgnored {
			if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
				columnsInsert.Append(tField.Column, pq.Array(field.Interface()))
				columnsUpdate.Append(tField.Column+" = ?", pq.Array(field.Interface()))
			} else {
				if tField.IsCreatedAt {
					returning.Append(tField.Column, field.Addr().Interface())
				} else if tField.IsUpdatedAt {
					returning.Append(tField.Column, field.Addr().Interface())
				} else if tField.IsDeletedAt {
					returning.Append(tField.Column, field.Addr().Interface())
				} else if tField.IsSequence {
					returning.Append(tField.Column, field.Addr().Interface())
				} else {
					columnsInsert.Append(tField.Column, field.Interface())
					columnsUpdate.Append(tField.Column+" = ?", field.Interface())
				}
			}
		}
	}
	if update {
		uQuery := gosql.NewUpdate()
		if condition != nil {
			uQuery.Where().Replace(condition)
		}
		uQuery.Table(model.Table())
		if columnsUpdate.Len() > 0 {
			uQuery.Set().Append(columnsUpdate.String(", "), columnsUpdate.GetArguments()...)
		}
		if returning.Len() > 0 {
			uQuery.Returning().Append(returning.String(", "), returning.GetArguments()...)
		}
		result = uQuery
	} else if insert {
		insertQuery := gosql.NewInsert()
		insertQuery.Into(model.Table())
		insertQuery.Columns().Add(columnsInsert.Split()...)
		insertQuery.Columns().Arg(columnsInsert.GetArguments()...)
		if returning.Len() > 0 {
			insertQuery.Returning().Append(returning.String(", "), returning.GetArguments()...)
		}
		result = insertQuery
	} else if upsert {
		upsertQuery := gosql.NewInsert()
		upsertQuery.Into(model.Table())
		upsertQuery.Columns().Add(columnsInsert.Split()...)
		upsertQuery.Columns().Arg(columnsInsert.GetArguments()...)
		if returning.Len() > 0 {
			upsertQuery.Returning().Append(returning.String(", "), returning.GetArguments()...)
		}
		if columnsUpdate.Len() > 0 {
			conflict.Set().Add(columnsUpdate.Split()...)
			conflict.Set().Arg(columnsUpdate.GetArguments()...)
		}
		conflict.Object(conflictObject.String(", "))
		upsertQuery.SetConflict(*conflict)
		result = upsertQuery
	}
	return result
}

// getSaveScenario check model for save scenario
// Helper method to understand how to get right save scenario
func getSaveScenario(model IModel) (insert, update, upsert bool) {
	if model == nil {
		return
	}
	ve := reflect.ValueOf(model).Elem()
	te := reflect.TypeOf(model).Elem()

	var hasPrimaryKey bool
	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		tField := ParseModelFiledTag(te.Field(i).Tag.Get("db"))
		if tField.IsPrimaryKey {
			hasPrimaryKey = true
			if !field.IsNil() {
				if tField.IsSequence {
					update = true
				} else {
					upsert = true
				}
			} else {
				if tField.IsSequence {
					insert = true
				} else {
					// conflict situation. Has primary, no value, no seq
					upsert = true
				}
			}
		} else if tField.IsUnique && !hasPrimaryKey {
			if !field.IsNil() {
				if tField.IsSequence {
					update = true
				} else if !insert && !update {
					upsert = true
				}
			} else {
				if tField.IsSequence {
					insert = true
				} else if !insert && !update {
					// conflict situation. Has unique, no value, no seq
					upsert = true
				}
			}
		}
	}
	return
}
