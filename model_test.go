package gomodel

import (
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
	var now = time.Now()
	m.CreatedAt = &now
	if *(meta.Fields[4].Value.(**time.Time)) != m.CreatedAt {
		t.Fatal("wrong field pointer CreatedAt")
	}
	if **(meta.Fields[4].Value.(**time.Time)) != *m.CreatedAt {
		t.Fatal("wrong field value CreatedAt")
	}
}

func BenchmarkPrepareMetaModel(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		PrepareMetaModel(m)
	}
	b.ReportAllocs()
}
