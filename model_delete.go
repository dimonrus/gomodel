package gomodel

import (
	"github.com/dimonrus/gosql"
	"github.com/lib/pq"
	"time"
)

// GetDeleteSQL model delete query
// model - target model
func GetDeleteSQL(model IModel) (iSQL gosql.ISQL) {
	if model == nil {
		return
	}
	var hasPrimaryKey bool
	meta := PrepareMetaModel(model)

	var now = time.Now()
	if meta.Fields.IsSoft() {
		upd := gosql.NewUpdate()
		for i := range meta.Fields {
			if meta.Fields[i].IsPrimaryKey {
				hasPrimaryKey = true
				if meta.Fields[i].Value != nil {
					if meta.Fields[i].IsArray {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
				}
			} else if meta.Fields[i].IsUnique && !hasPrimaryKey {
				if meta.Fields[i].Value != nil {
					if meta.Fields[i].IsArray {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						upd.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
				}
			} else if meta.Fields[i].IsDeletedAt {
				*meta.Fields[i].Value.(**time.Time) = &now
				upd.Set().Append(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
			} else if meta.Fields[i].IsUpdatedAt {
				upd.Set().Append(meta.Fields[i].Column + " = NOW()")
				upd.Returning().Append(meta.Fields[i].Column, meta.Fields[i].Value)
			}
		}
		if !upd.Where().IsEmpty() {
			upd.Table(model.Table())
			iSQL = upd
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
				}
			} else if meta.Fields[i].IsUnique && !hasPrimaryKey {
				if meta.Fields[i].Value != nil {
					if meta.Fields[i].IsArray {
						del.Where().AddExpression(meta.Fields[i].Column+" = ?", pq.Array(meta.Fields[i].Value))
					} else {
						del.Where().AddExpression(meta.Fields[i].Column+" = ?", meta.Fields[i].Value)
					}
				}
			}
		}
		if !del.Where().IsEmpty() {
			del.From(model.Table())
			iSQL = del
		}
	}
	return iSQL
}
