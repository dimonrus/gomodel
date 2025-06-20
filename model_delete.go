package gomodel

import (
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
)

// GetDeleteSQL model delete query
// model - target model
func GetDeleteSQL(model IModel) (iSQL gosql.ISQL) {
	isql := IndexCache.Get(IndexOperationDelete, model)
	if isql != nil {
		return isql
	}
	meta := PrepareMetaModel(model)
	if meta == nil {
		return
	}
	var hasPrimaryKey bool
	idx := InitIndex(meta.Fields.Len())
	if meta.Fields.IsSoft() {
		upd := gosql.NewUpdate()
		for i := range meta.Fields {
			if meta.Fields[i].IsPrimaryKey {
				hasPrimaryKey = true
				if !meta.Fields[i].IsNil {
					if meta.Fields[i].IsArray {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
					idx.AppendParamPos(int16(meta.Fields[i].Index))
				}
			} else if meta.Fields[i].IsUnique && !hasPrimaryKey {
				if !meta.Fields[i].IsNil {
					if meta.Fields[i].IsArray {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
					idx.AppendParamPos(int16(meta.Fields[i].Index))
				}
			} else if meta.Fields[i].IsDeletedAt {
				upd.Set().Append(meta.Fields[i].Column + " = NOW()")
				upd.Returning().Append(meta.Fields[i].Column, meta.Fields[i].Value)
				idx.AppendReturningPos(int16(meta.Fields[i].Index))
			} else if meta.Fields[i].IsUpdatedAt {
				upd.Set().Append(meta.Fields[i].Column + " = NOW()")
				upd.Returning().Append(meta.Fields[i].Column, meta.Fields[i].Value)
				idx.AppendReturningPos(int16(meta.Fields[i].Index))
			}
		}
		if !upd.Where().IsEmpty() {
			upd.Table(model.Table())
			iSQL = upd
			idx.SetQuery(upd.String())
		}
	} else {
		del := gosql.NewDelete()
		for i := range meta.Fields {
			if meta.Fields[i].IsPrimaryKey {
				hasPrimaryKey = true
				if meta.Fields[i].Value != nil {
					if meta.Fields[i].IsArray {
						del.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						del.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
					idx.AppendParamPos(int16(meta.Fields[i].Index))
				}
			} else if meta.Fields[i].IsUnique && !hasPrimaryKey {
				if meta.Fields[i].Value != nil {
					if meta.Fields[i].IsArray {
						del.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						del.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
					idx.AppendParamPos(int16(meta.Fields[i].Index))
				}
			}
		}
		if !del.Where().IsEmpty() {
			del.From(model.Table())
			iSQL = del
			idx.SetQuery(del.String())
		}
	}
	IndexCache.Store(IndexCache.Key(IndexOperationDelete, model.Table(), model.Columns(), model.Values()), idx)
	return iSQL
}
