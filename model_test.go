package gomodel

import (
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
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

func TestModelDeleteQuery(t *testing.T) {
	t.Run("with_cond", func(t *testing.T) {
		m := NewTestModel()
		c := gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
		c.AddExpression("created_at >= NOW()")
		query, e := GetDelete(m, c)
		if e != nil {
			t.Fatal(e)
		}
		if query.String() != "DELETE FROM test_model WHERE (created_at >= NOW());" {
			t.Fatal("Wrong sql prepared with_cond")
		}
	})
	t.Run("without_cond", func(t *testing.T) {
		m := NewTestModel()
		query, e := GetDelete(m, nil)
		if e != nil {
			t.Fatal(e)
		}
		if query.String() != "DELETE FROM test_model;" {
			t.Fatal("Wrong sql prepared")
		}
	})
}

func TestModelInsertQuery(t *testing.T) {
	m := NewTestModel()
	q := GetInsertSQL(m)
	q.Columns().Arg("one", []string{"1", "2", "3"}, 10, time.Now())
	q.Columns().Arg("two", []string{"2", "3", "4"}, 10, time.Now())
	t.Log(q.String())
	if q.String() != "INSERT INTO test_model (name, pages, some_int, created_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?);" {
		t.Fatal("wrong insert")
	}
	q = GetInsertSQL(m, &m.Name, &m.SomeInt)
	q.Columns().Arg("one", 10)
	q.Columns().Arg("two", 20)
	if q.String() != "INSERT INTO test_model (name, some_int) VALUES (?, ?), (?, ?), (?, ?);" {
		t.Fatal("wrong insert by columns")
	}

	q = GetInsertSQL(m, &m.Id, &m.Name, &m.Pages, &m.SomeInt)
	if q.String() != "INSERT INTO test_model (id, name, some_int) VALUES (?, ?, ?, ?);" {
		t.Fatal("wrong insert by keys")
	}
}

//classic BenchmarkModelInsertQuery-8   	  689485	      1741 ns/op	     224 B/op	      11 allocs/op
func BenchmarkModelInsertQuery(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		q := GetInsertSQL(m, &m.Id, &m.Name, &m.Pages)
		_ = q.String()
	}
	b.ReportAllocs()
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

func TestModelColumn(t *testing.T) {
	m := NewTestModel()
	cond := gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	cond.AddExpression("id = ?", 1)
	q, e := GetUpdate(m, cond, &m.Name, &m.SomeInt)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(q)
	if q.String() != "UPDATE test_model SET name = ?, some_int = ? WHERE (id = ?);" {
		t.Fatal("wrong update model")
	}
	if len(q.GetArguments()) != 3 {
		t.Fatal("wrong args count")
	}
}

func BenchmarkGetUpdate(b *testing.B) {
	m := NewTestModel()
	cond := gosql.NewSqlCondition(gosql.ConditionOperatorAnd)
	cond.AddExpression("id = ?", 1)
	for i := 0; i < b.N; i++ {
		_, _ = GetUpdate(m, cond, &m.Name, &m.SomeInt)
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
		GetColumns(m, &m.Name, &m.SomeInt)
	}
	b.ReportAllocs()
}
