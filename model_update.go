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
	isql := IndexCache.Get(IndexOperationUpdate, model, fields...)
	if isql != nil {
		return isql
	}
	meta := PrepareMetaModel(model)
	if meta == nil {
		return nil
	}
	idx := InitIndex(meta.Fields.Len())
	if fields == nil {
		fields = model.Values()
	}
	var conditionParams = make([]int16, 0, meta.Fields.Len())
	var hasPrimaryKey bool
	var condition = gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	var update = gosql.NewUpdate()
	for i := 0; i < meta.Fields.Len(); i++ {
		tField := meta.Fields[i]
		for _, v := range fields {
			cte := reflect.ValueOf(v)
			if cte.Kind() != reflect.Ptr {
				return nil
			}
			if reflect.ValueOf(tField.Value).Elem().Addr().Pointer() == cte.Elem().Addr().Pointer() {
				if tField.IsPrimaryKey {
					hasPrimaryKey = true
					if !tField.IsNil {
						if cte.Elem().Kind() == reflect.Slice || cte.Elem().Kind() == reflect.Array {
							condition.AddExpression(tField.Column+" = ?", pq.Array(tField.Value))
						} else {
							condition.AddExpression(tField.Column+" = ?", tField.Value)
						}
						conditionParams = append(conditionParams, int16(i))
					}
				} else if tField.IsUnique && !hasPrimaryKey {
					if !tField.IsNil {
						if cte.Elem().Kind() == reflect.Slice || cte.Elem().Kind() == reflect.Array {
							condition.AddExpression(tField.Column+" = ?", pq.Array(tField.Value))
						} else {
							condition.AddExpression(tField.Column+" = ?", tField.Value)
						}
						conditionParams = append(conditionParams, int16(i))
					}
				} else if !tField.IsIgnored {
					if tField.IsCreatedAt {
						update.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(i))
					} else if tField.IsUpdatedAt {
						if !tField.IsNil {
							update.Set().Append(tField.Column+" = ?", tField.Value)
							idx.AppendParamPos(int16(i))
						} else {
							update.Set().Append(tField.Column + " = NOW()")
						}
						update.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(i))
					} else if tField.IsDeletedAt {
						update.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(i))
					} else if tField.IsSequence {
						if !tField.IsNil {
							if cte.Elem().Kind() == reflect.Slice || cte.Elem().Kind() == reflect.Array {
								condition.AddExpression(tField.Column+" = ?", pq.Array(tField.Value))
							} else {
								condition.AddExpression(tField.Column+" = ?", tField.Value)
							}
							conditionParams = append(conditionParams, int16(i))
						}
						update.Returning().Append(tField.Column, tField.Value)
						idx.AppendReturningPos(int16(i))
					} else {
						if cte.Elem().Kind() == reflect.Slice || cte.Elem().Kind() == reflect.Array {
							update.Set().Append(tField.Column+" = ?", pq.Array(tField.Value))
						} else {
							update.Set().Append(tField.Column+" = ?", tField.Value)
						}
						idx.AppendParamPos(int16(i))
					}
				}
			}
		}
	}
	if update.IsEmpty() && condition.IsEmpty() {
		return nil
	}
	update.Table(model.Table())
	if !condition.IsEmpty() {
		update.Where().Replace(condition)
	}
	idx.SetQuery(update.String())
	idx.AppendParamPos(conditionParams...)
	IndexCache.Store(IndexCache.Key(IndexOperationUpdate, model.Table(), model.Columns(), model.Values(), fields...), idx)
	return update
}
