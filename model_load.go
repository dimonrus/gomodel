package gomodel

import (
	"github.com/dimonrus/gosql"
)

// GetLoadSQL return sql query fot load model
func GetLoadSQL(model IModel) gosql.ISQL {
	isql := IndexCache.Get(IndexOperationLoad, model)
	if isql != nil {
		return isql
	}
	meta := PrepareMetaModel(model)
	if meta == nil {
		return nil
	}
	selectSql := gosql.NewSelect()
	selectSql.From(model.Table())
	cond := gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	idx := InitIndex(meta.Fields.Len())
	for i := 0; i < meta.Fields.Len(); i++ {
		tField := meta.Fields[i]
		if tField.IsIgnored || tField.Column == "" {
			continue
		}
		if tField.IsPrimaryKey && !tField.IsNil {
			cond.AddExpression(tField.Column+" = ?", tField.Value)
			idx.AppendParamPos(int16(i))
		} else if tField.IsUnique && !tField.IsNil {
			if cond.IsEmpty() {
				cond.AddExpression(tField.Column+" = ?", tField.Value)
				idx.AppendParamPos(int16(i))
			}
		} else if tField.IsDeletedAt {
			cond.AddExpression(tField.Column + " IS NULL")
		}
		selectSql.Columns().Append(tField.Column, tField.Value)
		idx.AppendReturningPos(int16(i))
	}
	if !cond.IsEmpty() {
		selectSql.Where().Replace(cond)
	}
	idx.SetQuery(selectSql.String())
	IndexCache.Store(IndexCache.Key(IndexOperationLoad, model.Table(), model.Columns(), model.Values()), idx)
	return selectSql
}
