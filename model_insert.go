package gomodel

import (
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
	"reflect"
)

// GetInsertSQL model insert query
func GetInsertSQL(model IModel, fields ...any) gosql.ISQL {
	me := reflect.ValueOf(model)
	te := reflect.TypeOf(model).Elem()
	if me.IsNil() {
		return nil
	}
	me = me.Elem()
	var values []any
	var insert = gosql.NewInsert()
	var conflict = gosql.NewConflict()
	var hasPrimaryKey bool
	var isConflict bool
	var tField ModelFiledTag

	// use when fields directly passed in args
	var direct = len(fields) > 0
	// define values for iteration
	if direct {
		values = fields
	} else {
		values = model.Values()
	}
	for j := 0; j < me.NumField(); j++ {
		field := me.Field(j)
		tField.Clear()
		ParseModelFiledTag(te.Field(j).Tag.Get("db"), &tField)
		for i := range values {
			fv := reflect.ValueOf(values[i])
			if fv.Kind() != reflect.Ptr {
				continue
			}
			if field.Addr().Pointer() == fv.Elem().Addr().Pointer() {
				if tField.IsPrimaryKey {
					hasPrimaryKey = true
					if !field.IsNil() {
						if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
							insert.Columns().Append(tField.Column, pq.Array(field.Interface()))
						} else {
							insert.Columns().Append(tField.Column, field.Interface())
						}
						if !tField.IsSequence {
							isConflict = true
							conflict.Object(tField.Column)
						}
					} else {
						if !tField.IsSequence {
							if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
								insert.Columns().Append(tField.Column, pq.Array(field.Interface()))
							} else {
								insert.Columns().Append(tField.Column, field.Interface())
							}
						} else {
							insert.Returning().Append(tField.Column, field.Addr().Interface())
						}
					}
				} else if tField.IsUnique && !hasPrimaryKey {
					if !field.IsNil() {
						if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
							insert.Columns().Append(tField.Column, pq.Array(field.Interface()))
						} else {
							insert.Columns().Append(tField.Column, field.Interface())
						}
						if !tField.IsSequence {
							isConflict = true
							conflict.Object(tField.Column)
						}
					} else {
						if !tField.IsSequence {
							if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
								insert.Columns().Append(tField.Column, pq.Array(field.Interface()))
							} else {
								insert.Columns().Append(tField.Column, field.Interface())
							}
						} else {
							insert.Returning().Append(tField.Column, field.Addr().Interface())
						}
					}
				} else if !tField.IsIgnored {
					if tField.IsCreatedAt {
						insert.Returning().Append(tField.Column, field.Addr().Interface())
					} else if tField.IsUpdatedAt {
						insert.Returning().Append(tField.Column, field.Addr().Interface())
					} else if tField.IsDeletedAt {
						insert.Returning().Append(tField.Column, field.Addr().Interface())
					} else if tField.IsSequence {
						insert.Returning().Append(tField.Column, field.Addr().Interface())
					} else {
						if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
							insert.Columns().Append(tField.Column, pq.Array(field.Interface()))
							conflict.Set().Append(tField.Column+" = ?", pq.Array(field.Interface()))
						} else {
							conflict.Set().Append(tField.Column+" = ?", field.Interface())
							insert.Columns().Append(tField.Column, field.Interface())
						}
					}
				}
			}
		}
	}
	if !insert.IsEmpty() {
		insert.Into(gohelp.ToUnderscore(te.Name()))
		if isConflict {
			insert.Conflict().Object(conflict.GetObject())
			insert.Conflict().Action(gosql.ConflictActionUpdate)
			insert.Conflict().Set().Add(conflict.Set().Split()...)
			insert.Conflict().Set().Arg(conflict.Set().GetArguments()...)
		}
	}
	return insert
}
