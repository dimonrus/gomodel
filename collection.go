package gomodel

import (
	"database/sql"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
	"github.com/dimonrus/porterr"
)

// Collection struct contain items and collection common methods
type Collection[T any] struct {
	// List of collection items
	items []*T
	// Iterator
	*gohelp.Iterator
	// Query builder
	*gosql.Select
	// Count
	CountOver int
}

// Items Get all items
func (c *Collection[T]) Items() []*T {
	return c.items
}

// AddItem Add item to collection
func (c *Collection[T]) AddItem(item ...*T) {
	c.items = append(c.items, item...)
	c.SetCount(len(c.items))
	return
}

// First get item of model collection
func (c *Collection[T]) First() *T {
	if len(c.items) > 0 {
		return c.items[0]
	}
	return nil
}

// Last get item of model collection
func (c *Collection[T]) Last() *T {
	if c.Count() > 0 {
		return c.items[c.Count()-1]
	}
	return nil
}

// Item get item of collection
func (c *Collection[T]) Item() *T {
	if c.Cursor() < c.Count() && c.Cursor() > -1 {
		return c.items[c.Cursor()]
	}
	return nil
}

// Clear collection
func (c *Collection[T]) clear() {
	c.Iterator.Reset()
	c.items = c.items[:0]
}

// fetch collection data private method
func (c *Collection[T]) preload(q godb.Queryer) (rows *sql.Rows, e porterr.IError) {
	var err error
	rows, err = q.Query(c.String(), c.GetArguments()...)
	if err != nil {
		e = porterr.New(porterr.PortErrorDatabaseQuery, "Collection search query error: "+err.Error())
	}
	return
}

// scan collection method
func (c *Collection[T]) scan(rows *sql.Rows) (e porterr.IError) {
	if rows == nil {
		return
	}
	// row values
	var values []interface{}
	for rows.Next() {
		var model interface{} = new(T)
		values = (model).(IModel).Values()
		if c.CountOver >= 0 {
			values = append(values, &c.CountOver)
		}
		err := rows.Scan(values...)
		if err != nil {
			e = porterr.New(porterr.PortErrorIO, (model).(IModel).Table()+" model scan error: "+err.Error())
			return
		}
		c.AddItem(model.(*T))
	}
	return
}

// AddCountOver add count column to SQL
func (c *Collection[T]) AddCountOver() {
	if c.CountOver < 0 {
		c.Columns().Add("COUNT(*) OVER()")
		c.CountOver = 0
	}
}

// RemoveCountOver remove count column from SQL
func (c *Collection[T]) RemoveCountOver() {
	var item interface{} = new(T)
	c.CountOver = -1
	c.Columns().Reset()
	c.Columns().Add(item.(IModel).Columns()...)
}

// Load collection
func (c *Collection[T]) Load(q godb.Queryer) porterr.IError {
	rows, e := c.preload(q)
	if e != nil {
		return e
	}
	defer func() { _ = rows.Close() }()
	c.clear()
	return c.scan(rows)
}

// Map collection
func (c *Collection[T]) Map(callback func(*T)) {
	c.Reset()
	for c.Next() {
		callback(c.Item())
	}
	return
}

// Filter collection
func (c *Collection[T]) Filter(callback func(*T) bool) {
	c.Reset()
	i := 0
	items := make([]*T, c.Count())
	for c.Next() {
		if callback(c.Item()) {
			items[i] = c.Item()
			i++
		}
	}
	c.items = items[:i]
	c.SetCount(len(c.items))
	return
}

// Save Create or Update collection items
func (c *Collection[T]) Save(q godb.Queryer) (e porterr.IError) {
	var m interface{} = new(T)
	if _, ok := m.(IModel); !ok {
		e = porterr.New(porterr.PortErrorArgument, "Type T is not implements IModel interface")
		return
	}
	var stmts = make(map[string]*godb.SqlStmt)
	defer func() {
		for s := range stmts {
			_ = stmts[s].Close()
		}
	}()
	var err error
	for c.Next() {
		query, params, returning := GetSaveSQL((interface{})(c.Item()).(IModel)).SQL()
		if _, ok := stmts[query]; !ok {
			stmts[query], err = q.Prepare(query)
			if err != nil {
				return porterr.New(porterr.PortErrorIO, err.Error())
			}
		}
		err = stmts[query].QueryRow(params...).Scan(returning...)
		if err != nil {
			return porterr.New(porterr.PortErrorIO, err.Error())
		}
	}
	return
}

// NewCollection Create new model collection
func NewCollection[T any]() *Collection[T] {
	var item interface{} = new(T)
	// *T must implements IModel interface
	if _, ok := item.(IModel); !ok {
		return nil
	}
	query := gosql.NewSelect()
	query.Columns().Add(item.(IModel).Columns()...)
	query.From(item.(IModel).Table())
	collection := &Collection[T]{
		items:     make([]*T, 0),
		Select:    query,
		Iterator:  gohelp.NewIterator(0),
		CountOver: -1,
	}
	return collection
}
