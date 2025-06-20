package gomodel

import (
	"database/sql"
	"github.com/dimonrus/gohelp"
	"github.com/lib/pq"
	"testing"
	"time"
)

type TestModel struct {
	Id        *int       `json:"id" db:"col~id;req;seq;"`
	Name      *string    `json:"name" db:"col~name;req;"`
	Pages     []string   `json:"pages" db:"col~pages;"`
	SomeInt   *int       `json:"someInt" db:"col~some_int;"`
	CreatedAt *time.Time `json:"createdAt" db:"col~created_at;cat;"`
	Custom    *int32     `json:"custom" db:"ign"`
}

// Model table name
func (m *TestModel) Table() string {
	return "test_model"
}

// Model columns
func (m *TestModel) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at"}
}

// Model values
func (m *TestModel) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt}
}

func NewTestModel() *TestModel {
	id := gohelp.GetRndId()
	someInt := gohelp.GetRndNumber(10, 3000)
	name := gohelp.RandString(10)
	pages := []string{"one", "two"}
	return &TestModel{
		Id:      &id,
		Name:    &name,
		Pages:   pages,
		SomeInt: &someInt,
	}
}

func TestModelValues(t *testing.T) {
	m := NewTestModel()
	values := GetValues(m, "id", "pages", "some_int")
	if len(values) != 3 {
		t.Fatal("wrong values len")
	}
	**values[0].(**int) = 10
	if *m.Id != 10 {
		t.Fatal("wrong conversation")
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkModelValues
// BenchmarkModelValues-12    	 2906820	       404.2 ns/op	     128 B/op	       2 allocs/op
func BenchmarkModelValues(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		GetValues(m, "id", "pages", "some_int")
	}
	b.ReportAllocs()
}

func TestModelColumns(t *testing.T) {
	m := NewTestModel()
	columns := GetColumns(m, &m.Name, &m.SomeInt, &m.Pages)
	if len(columns) != 3 {
		t.Fatal("wrong")
	}
	if columns[0] != "name" {
		t.Fatal("wrong")
	}
	if columns[1] != "some_int" {
		t.Fatal("wrong")
	}
	if columns[2] != "pages" {
		t.Fatal("wrong")
	}
}

func BenchmarkGetColumns(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		GetColumns(m)
	}
	b.ReportAllocs()
}

func BenchmarkExtract(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		extract(m)
	}
	b.ReportAllocs()
}

func TestPrepareMetaModel(t *testing.T) {
	m := NewTestModel()
	meta := PrepareMetaModel(m)
	if meta.TableName != "test_model" {
		t.Fatal("wrong table name")
	}
	if meta.Fields.Len() != 5 {
		t.Fatal("wrong field count")
	}
	if *(meta.Fields[0].Value.(**int)) != m.Id {
		t.Fatal("wrong field pointer id")
	}
	if **(meta.Fields[0].Value.(**int)) != *m.Id {
		t.Fatal("wrong field value id")
	}
	if !meta.Fields[4].IsNil {
		t.Fatal("must be nil")
	}
	var now = time.Now()
	m.CreatedAt = &now
	if *(meta.Fields[4].Value.(**time.Time)) != m.CreatedAt {
		t.Fatal("wrong field pointer CreatedAt")
	}
	if **(meta.Fields[4].Value.(**time.Time)) != *m.CreatedAt {
		t.Fatal("wrong field value CreatedAt")
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkPrepareMetaModel
// BenchmarkPrepareMetaModel-12    	 2036199	       583.7 ns/op	     496 B/op	       2 allocs/op
func BenchmarkPrepareMetaModel(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		PrepareMetaModel(m)
	}
	b.ReportAllocs()
}

type DefaultCustomModel struct {
	Id        int                `json:"id" db:"col~id;req;seq;"`
	Name      sql.NullString     `json:"name" db:"col~name;req;"`
	Pages     []string           `json:"pages" db:"col~pages;"`
	SomeInt   *float64           `json:"someInt" db:"col~some_int;"`
	CreatedAt time.Time          `json:"createdAt" db:"col~created_at;cat;"`
	Custom    *struct{ Foo int } `json:"custom" db:"col~custom;"`
}

// Model table name
func (m *DefaultCustomModel) Table() string {
	return "default_custom_model"
}

// Model columns
func (m *DefaultCustomModel) Columns() []string {
	return []string{"id", "name", "pages", "some_int", "created_at", "custom"}
}

// Model values
func (m *DefaultCustomModel) Values() []any {
	return []any{&m.Id, &m.Name, pq.Array(&m.Pages), &m.SomeInt, &m.CreatedAt, &m.Custom}
}

func TestPrepareMetaModelCustom(t *testing.T) {
	m := &DefaultCustomModel{}
	meta := PrepareMetaModel(m)
	if meta.Fields.Len() != 6 {
		t.Fatal("wrong fields count")
	}
	if *(meta.Fields[0].Value.(*int)) != m.Id {
		t.Fatal("wrong field pointer id")
	}
	if !meta.Fields[0].IsZero {
		t.Fatal("id must be zero")
	}
	if meta.Fields[0].IsNil {
		t.Fatal("id must be not nil")
	}
	if *(meta.Fields[1].Value.(*sql.NullString)) != m.Name {
		t.Fatal("wrong field pointer name")
	}
	if !meta.Fields[1].IsZero {
		t.Fatal("name must be zero")
	}
	if meta.Fields[0].IsNil {
		t.Fatal("name must be not nil")
	}
	if len(*(meta.Fields[2].Value.(*[]string))) != len(m.Pages) {
		t.Fatal("wrong field pointer pages")
	}
	if !meta.Fields[2].IsZero {
		t.Fatal("pages must be zero")
	}
	if meta.Fields[2].IsNil {
		t.Fatal("pages must be not nil")
	}
	if *(meta.Fields[3].Value.(**float64)) != m.SomeInt {
		t.Fatal("wrong field pointer some int")
	}
	if !meta.Fields[3].IsZero {
		t.Fatal("someint must be zero")
	}
	if !meta.Fields[3].IsNil {
		t.Fatal("someint must be nil")
	}
	if *(meta.Fields[4].Value.(*time.Time)) != m.CreatedAt {
		t.Fatal("wrong field pointer created at")
	}
	if !meta.Fields[4].IsZero {
		t.Fatal("created must be zero")
	}
	if meta.Fields[4].IsNil {
		t.Fatal("created at must be not nil")
	}
	if *(meta.Fields[5].Value.(**struct{ Foo int })) != m.Custom {
		t.Fatal("wrong field pointer custom")
	}
	if !meta.Fields[5].IsZero {
		t.Fatal("custom must be zero")
	}
	if !meta.Fields[5].IsNil {
		t.Fatal("custom must be nil")
	}
}

type ComplexTestModel struct {
	Pages      []string           `json:"pages" db:"col~pages;"`
	Name       sql.NullString     `json:"name" db:"col~name;req;"`
	ComplexId  int                `json:"complexId" db:"col~complex_id;req;prk;"`
	CategoryId *int               `json:"categoryId" db:"col~category_id;req;prk;"`
	SomeInt    *float64           `json:"someInt" db:"col~some_int;"`
	CreatedAt  time.Time          `json:"createdAt" db:"col~created_at;cat;"`
	Custom     *struct{ Foo int } `json:"custom" db:"col~custom;"`
}

// Model table name
func (m *ComplexTestModel) Table() string {
	return "complex_test_update_model"
}

// Model columns
func (m *ComplexTestModel) Columns() []string {
	return []string{"pages", "name", "complex_id", "category_id", "some_int", "created_at", "custom"}
}

// Model values
func (m *ComplexTestModel) Values() []any {
	return []any{pq.Array(&m.Pages), &m.Name, &m.ComplexId, &m.CategoryId, &m.SomeInt, &m.CreatedAt, &m.Custom}
}
