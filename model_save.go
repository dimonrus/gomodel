package gomodel

import (
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
	"reflect"
	"strings"
)

// GetSaveSQL prepare a save query
// it can be insert or update or upsert query
// some popular scenario was implemented. not all
func GetSaveSQL(model IModel) gosql.ISQL {
	var result gosql.ISQL
	insert, update, upsert := getSaveScenario(model)
	if insert {
		result = IndexCache.Get(IndexOperationCreate, model)
	} else if update {
		result = IndexCache.Get(IndexOperationUpdate, model)
	} else if upsert {
		result = IndexCache.Get(IndexOperationSave, model)
	}
	if result != nil {
		return result
	}
	meta := PrepareMetaModel(model)
	if meta == nil {
		return nil
	}
	idx := InitIndex(meta.Fields.Len())
	var columnsInsert = gosql.NewExpression()
	var columnsUpdate = gosql.NewExpression()
	var returning = gosql.NewExpression()
	var condition = gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	var conflict = gosql.NewConflict().Action(gosql.ConflictActionUpdate)
	var hasPrimaryKey bool

	var conditionPos = make([]int16, 0, meta.Fields.Len())
	var columnPos = make([]int16, 0, meta.Fields.Len())
	var upsertPos = make([]int16, 0, meta.Fields.Len())
	var conflictColumns = strings.Builder{}

	var tField ModelFiledTag
	for i := 0; i < meta.Fields.Len(); i++ {
		tField = meta.Fields[i]
		if tField.IsPrimaryKey {
			hasPrimaryKey = true
			if !tField.IsNil {
				if tField.IsSequence {
					update = true
					condition.AddExpression(tField.Column+" = ?", tField.Value)
					conditionPos = append(conditionPos, int16(i))
				} else {
					upsert = true
					if tField.IsArray {
						columnsInsert.Append(tField.Column, pq.Array(tField.Value))
					} else {
						columnsInsert.Append(tField.Column, tField.Value)
					}
					columnPos = append(columnPos, int16(i))
					if conflictColumns.Len() > 0 {
						conflictColumns.WriteString(", ")
					}
					conflictColumns.WriteString(tField.Column)
				}
			} else {
				if tField.IsSequence {
					insert = true
					returning.Append(tField.Column, tField.Value)
					idx.AppendReturningPos(int16(i))
				} else {
					upsert = true
					if tField.IsArray {
						columnsInsert.Append(tField.Column, pq.Array(tField.Value))
					} else {
						columnsInsert.Append(tField.Column, tField.Value)
					}
					columnPos = append(columnPos, int16(i))
					if conflictColumns.Len() > 0 {
						conflictColumns.WriteString(", ")
					}
					conflictColumns.WriteString(tField.Column)
				}
			}
		} else if tField.IsUnique && !hasPrimaryKey {
			if !tField.IsNil {
				if tField.IsSequence {
					update = true
					condition.AddExpression(tField.Column+" = ?", tField.Value)
					conditionPos = append(conditionPos, int16(i))
				} else if !insert && !update {
					upsert = true
					if tField.IsArray {
						columnsInsert.Append(tField.Column, pq.Array(tField.Value))
					} else {
						columnsInsert.Append(tField.Column, tField.Value)
					}
					columnPos = append(columnPos, int16(i))
					if conflictColumns.Len() > 0 {
						conflictColumns.WriteString(", ")
					}
					conflictColumns.WriteString(tField.Column)
				}
			} else {
				if tField.IsSequence {
					insert = true
					returning.Append(tField.Column, tField.Value)
					idx.AppendReturningPos(int16(i))
				} else if !insert && !update {
					upsert = true
					if tField.IsArray {
						columnsInsert.Append(tField.Column, pq.Array(tField.Value))
					} else {
						columnsInsert.Append(tField.Column, tField.Value)
					}
					columnPos = append(columnPos, int16(i))
					if conflictColumns.Len() > 0 {
						conflictColumns.WriteString(", ")
					}
					conflictColumns.WriteString(tField.Column)
				}
			}
		} else if !tField.IsIgnored {
			if tField.IsArray {
				columnsInsert.Append(tField.Column, pq.Array(tField.Value))
				columnsUpdate.Append(tField.Column+" = ?", pq.Array(tField.Value))
				columnPos = append(columnPos, int16(i))
				upsertPos = append(upsertPos, int16(i))
			} else {
				if tField.IsCreatedAt {
					returning.Append(tField.Column, tField.Value)
					idx.AppendReturningPos(int16(i))
				} else if tField.IsUpdatedAt {
					columnsUpdate.Append(tField.Column + " = NOW()")
					returning.Append(tField.Column, tField.Value)
					idx.AppendReturningPos(int16(i))
				} else if tField.IsDeletedAt {
					returning.Append(tField.Column, tField.Value)
					idx.AppendReturningPos(int16(i))
				} else if tField.IsSequence {
					returning.Append(tField.Column, tField.Value)
					idx.AppendReturningPos(int16(i))
				} else {
					columnsInsert.Append(tField.Column, tField.Value)
					columnsUpdate.Append(tField.Column+" = ?", tField.Value)
					columnPos = append(columnPos, int16(i))
					upsertPos = append(upsertPos, int16(i))
				}
			}
		}
	}
	if update {
		uQuery := gosql.NewUpdate()
		uQuery.Table(model.Table())
		if columnsUpdate.Len() > 0 {
			uQuery.Set().Append(columnsUpdate.String(", "), columnsUpdate.GetArguments()...)
			idx.AppendParamPos(columnPos...)
		}
		if condition != nil {
			uQuery.Where().Replace(condition)
			idx.AppendParamPos(conditionPos...)
		}
		if returning.Len() > 0 {
			uQuery.Returning().Append(returning.String(", "), returning.GetArguments()...)
		}
		idx.SetQuery(uQuery.String())
		result = uQuery
	} else if insert {
		insertQuery := gosql.NewInsert()
		insertQuery.Into(model.Table())
		insertQuery.Columns().Add(columnsInsert.Split()...)
		insertQuery.Columns().Arg(columnsInsert.GetArguments()...)
		idx.AppendParamPos(columnPos...)
		if returning.Len() > 0 {
			insertQuery.Returning().Append(returning.String(", "), returning.GetArguments()...)
		}
		idx.SetQuery(insertQuery.String())
		result = insertQuery
	} else if upsert {
		upsertQuery := gosql.NewInsert()
		upsertQuery.Into(model.Table())
		upsertQuery.Columns().Add(columnsInsert.Split()...)
		upsertQuery.Columns().Arg(columnsInsert.GetArguments()...)
		idx.AppendParamPos(columnPos...)
		if returning.Len() > 0 {
			upsertQuery.Returning().Append(returning.String(", "), returning.GetArguments()...)
		}
		if columnsUpdate.Len() > 0 {
			conflict.Set().Add(columnsUpdate.Split()...)
			conflict.Set().Arg(columnsUpdate.GetArguments()...)
			idx.AppendParamPos(upsertPos...)
		}
		if conflictColumns.Len() > 0 {
			conflict.Object(conflictColumns.String())
			upsertQuery.SetConflict(*conflict)
		}
		idx.SetQuery(upsertQuery.String())
		result = upsertQuery
	}
	var key ModelOperation
	if insert {
		key = IndexCache.Key(IndexOperationCreate, model.Table(), model.Columns(), model.Values())
	} else if update {
		key = IndexCache.Key(IndexOperationUpdate, model.Table(), model.Columns(), model.Values())
	} else if upsert {
		key = IndexCache.Key(IndexOperationSave, model.Table(), model.Columns(), model.Values())
	}
	IndexCache.Store(key, idx)
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
	var tField ModelFiledTag
	for i := 0; i < ve.NumField(); i++ {
		field := ve.Field(i)
		ParseModelFiledTag(te.Field(i).Tag.Get("db"), &tField)
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
		tField.Clear()
	}
	return
}
