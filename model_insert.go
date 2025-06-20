package gomodel

import (
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
	"reflect"
	"strings"
)

// GetInsertSQL model insert query
func GetInsertSQL(model IModel, fields ...any) gosql.ISQL {
	isql := IndexCache.Get(IndexOperationCreate, model, fields...)
	if isql != nil {
		return isql
	}
	meta := PrepareMetaModel(model)
	idx := InitIndex(meta.Fields.Len())
	var values []any
	var insert = gosql.NewInsert()
	var conflict = gosql.NewConflict()
	var hasPrimaryKey bool
	var isConflict bool
	var updateSetPos = make([]int16, 0, meta.Fields.Len())
	var conflictColumns = strings.Builder{}

	// use when fields directly passed in args
	var direct = len(fields) > 0
	// define values for iteration
	if direct {
		values = fields
	} else {
		values = model.Values()
	}
	for j := 0; j < meta.Fields.Len(); j++ {
		tField := meta.Fields[j]
		for i := range values {
			fv := reflect.ValueOf(values[i])
			if fv.Kind() != reflect.Ptr {
				continue
			}
			if reflect.ValueOf(tField.Value).Elem().Addr().Pointer() == fv.Elem().Addr().Pointer() {
				if tField.IsPrimaryKey {
					hasPrimaryKey = true
					if !tField.IsNil {
						if fv.Elem().Kind() == reflect.Slice || fv.Elem().Kind() == reflect.Array {
							insert.Columns().Append(tField.Column, pq.Array(tField.Value))
						} else {
							insert.Columns().Append(tField.Column, tField.Value)
						}
						idx.AppendParamPos(int16(j))
						if !tField.IsSequence {
							isConflict = true
							if conflictColumns.Len() > 0 {
								conflictColumns.WriteString(", ")
							}
							conflictColumns.WriteString(tField.Column)
						}
					} else {
						if !tField.IsSequence {
							if fv.Elem().Kind() == reflect.Slice || fv.Elem().Kind() == reflect.Array {
								insert.Columns().Append(tField.Column, pq.Array(tField.Value))
							} else {
								insert.Columns().Append(tField.Column, tField.Value)
							}
							idx.AppendParamPos(int16(j))
						} else {
							insert.Returning().Append(tField.Column, tField.Value)
							idx.AppendReturningPos(int16(j))
						}
					}
				} else if tField.IsUnique && !hasPrimaryKey {
					if !tField.IsNil {
						if fv.Elem().Kind() == reflect.Slice || fv.Elem().Kind() == reflect.Array {
							insert.Columns().Append(tField.Column, pq.Array(tField.Value))
						} else {
							insert.Columns().Append(tField.Column, tField.Value)
						}
						idx.AppendParamPos(int16(j))
						if !tField.IsSequence {
							isConflict = true
							if conflictColumns.Len() > 0 {
								conflictColumns.WriteString(", ")
							}
							conflictColumns.WriteString(tField.Column)
						}
					} else {
						if !tField.IsSequence {
							if fv.Elem().Kind() == reflect.Slice || fv.Elem().Kind() == reflect.Array {
								insert.Columns().Append(tField.Column, pq.Array(tField.Value))
							} else {
								insert.Columns().Append(tField.Column, tField.Value)
							}
							idx.AppendParamPos(int16(j))
						} else {
							insert.Returning().Append(tField.Column, tField.Value)
							idx.AppendReturningPos(int16(j))
						}
					}
				} else if !tField.IsIgnored {
					if tField.IsCreatedAt {
						insert.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(j))
					} else if tField.IsUpdatedAt {
						insert.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(j))
					} else if tField.IsDeletedAt {
						insert.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(j))
					} else if tField.IsSequence {
						insert.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(j))
					} else {
						if fv.Elem().Kind() == reflect.Slice || fv.Elem().Kind() == reflect.Array {
							insert.Columns().Append(tField.Column, pq.Array(tField.Value))
							conflict.Set().Append(tField.Column+" = ?", pq.Array(tField.Value))
						} else {
							conflict.Set().Append(tField.Column+" = ?", tField.Value)
							insert.Columns().Append(tField.Column, tField.Value)
						}
						idx.AppendParamPos(int16(j))
						updateSetPos = append(updateSetPos, int16(j))
					}
				}
			}
		}
	}
	if !insert.IsEmpty() {
		insert.Into(model.Table())
		if isConflict {
			insert.Conflict().Object(conflictColumns.String())
			insert.Conflict().Action(gosql.ConflictActionUpdate)
			insert.Conflict().Set().Add(conflict.Set().Split()...)
			insert.Conflict().Set().Arg(conflict.Set().GetArguments()...)
			if len(updateSetPos) > 0 {
				idx.AppendParamPos(updateSetPos...)
			}
		}
	}
	idx.SetQuery(insert.String())
	IndexCache.Store(IndexCache.Key(IndexOperationCreate, model.Table(), model.Columns(), model.Values(), fields...), idx)
	return insert
}
