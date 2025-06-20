package gomodel

import (
	"github.com/dimonrus/gosql"
	"sync"
)

// IndexOperation type of operation for specific model
type IndexOperation string

// ModelOperation concat of model table(1) and index operation(2)
// example "table_create" or "dictionary_load"
type ModelOperation string

const (
	// IndexOperationLoad load operation
	IndexOperationLoad IndexOperation = "load"
	// IndexOperationCreate create operation
	IndexOperationCreate IndexOperation = "create"
	// IndexOperationUpdate update operation
	IndexOperationUpdate IndexOperation = "update"
	// IndexOperationDelete delete operation
	IndexOperationDelete IndexOperation = "delete"
	// IndexOperationSave save operation
	IndexOperationSave IndexOperation = "save"

	// IndexCacheDefaultLength Init size for a cache map
	IndexCacheDefaultLength = 16
)

// IndexCache index cache object
var IndexCache = &cache{
	models:  make(map[ModelOperation]Index, IndexCacheDefaultLength),
	columns: make(map[string][]string, 16),
}

// cache type
type cache struct {
	// list of cached items
	models map[ModelOperation]Index
	// list of columns
	columns map[string][]string
	// rw mutex
	m sync.RWMutex
}

// set model columns to cache
func (c *cache) setModelColumns(table string, columns []string) {
	c.m.Lock()
	defer c.m.Unlock()
	if _, ok := c.columns[table]; ok {
		return
	}
	c.columns[table] = columns
	return
}

// set model columns to cache
func (c *cache) getModelColumns(table string) (columns []string) {
	c.m.RLock()
	defer c.m.RUnlock()
	if v, ok := c.columns[table]; ok {
		return v
	}
	return
}

// Get a model index object
func (c *cache) Get(io IndexOperation, model IModel, field ...any) gosql.ISQL {
	table := model.Table()
	values := model.Values()
	columns := c.getModelColumns(table)
	if columns == nil {
		columns = model.Columns()
		c.setModelColumns(table, columns)
	}
	key := c.Key(io, table, columns, values, field...)
	if v, ok := c.models[key]; ok {
		return v.ToISQL(values)
	}
	return nil
}

// Reset map
func (c *cache) Reset() {
	c.m.Lock()
	defer c.m.Unlock()
	c.models = make(map[ModelOperation]Index, IndexCacheDefaultLength)
}

// Key GET cache key
func (c *cache) Key(io IndexOperation, table string, columns []string, values []any, field ...any) ModelOperation {
	var suffix = []byte("_full_")
	if len(field) > 0 && len(field) < len(columns) {
		suffix = suffix[:0]
		for i := range field {
			for j := range values {
				if field[i] == values[j] {
					suffix = append(suffix, byte(j))
					break
				}
			}
		}
	}
	return ModelOperation(table + string(suffix) + string(io))
}

// Store model index object
func (c *cache) Store(mo ModelOperation, i Index) {
	c.m.Lock()
	defer c.m.Unlock()
	c.models[mo] = i
}

// Index model index internal struct
type Index struct {
	// sql query
	query string
	// position for parameters
	paramsPos []int16
	// positions for returning values
	returningPos []int16
}

// SetQuery set query
func (c *Index) SetQuery(query string) {
	c.query = query
}

// AppendParamPos add parameter position
func (c *Index) AppendParamPos(pos ...int16) {
	c.paramsPos = append(c.paramsPos, pos...)
}

// PrependParamPos starts with parameter position
func (c *Index) PrependParamPos(pos ...int16) {
	c.paramsPos = append(pos, c.paramsPos...)
}

// AppendReturningPos add returning position
func (c *Index) AppendReturningPos(pos ...int16) {
	c.returningPos = append(c.returningPos, pos...)
}

// PrependReturningPos add returning position
func (c *Index) PrependReturningPos(pos ...int16) {
	c.returningPos = append(pos, c.returningPos...)
}

// ToISQL prepare ISQL struct
func (c *Index) ToISQL(values []any) gosql.ISQL {
	if c == nil {
		return nil
	}
	index := indexISQL{
		query:     c.query,
		params:    make([]any, len(c.paramsPos)),
		returning: make([]any, len(c.returningPos)),
	}
	for i := range index.params {
		index.params[i] = values[c.paramsPos[i]]
	}
	for i := range index.returning {
		index.returning[i] = values[c.returningPos[i]]
	}
	return index
}

// result gosql.ISQL struct with a query, params and returning
type indexISQL struct {
	// string query
	query string
	// list of query params
	params []any
	// list of returning params
	returning []any
}

// SQL Implementation for gosql.ISQL
func (c indexISQL) SQL() (query string, params []any, returning []any) {
	return c.query, c.params, c.returning
}

// InitIndex Init index object
func InitIndex(size int) Index {
	return Index{paramsPos: make([]int16, 0, size), returningPos: make([]int16, 0, size)}
}
