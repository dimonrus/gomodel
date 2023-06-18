package gomodel

import (
	"testing"
)

func TestGetUpdateSQL(t *testing.T) {
	t.Run("classic_update", func(t *testing.T) {
		m := NewTestModel()
		m.Id = new(int)
		*m.Id = 1000
		q := GetUpdateSQL(m)
		query, params, returning := q.SQL()
		t.Log(query)
		if query != "UPDATE test_model SET name = ?, pages = ?, some_int = ? WHERE (id = ?) RETURNING id, created_at;" {
			t.Fatal("wrong classic_update")
		}
		if len(params) != 4 {
			t.Fatal("wrong classic_update args len must be 4")
		}
		if params[0] != m.Name {
			t.Fatal("wrong classic_update args 0 must be name")
		}
		if params[3] != m.Id {
			t.Fatal("wrong classic_update args 0 must be id")
		}
		if len(returning) != 2 {
			t.Fatal("wrong classic_update returning arg len must be 2")
		}
		if returning[0] != &m.Id {
			t.Fatal("wrong classic_update returning id")
		}
	})

	t.Run("classic_update_1", func(t *testing.T) {
		m := &UpdateModel1{}
		m.Id = new(int)
		*m.Id = 1000
		q := GetUpdateSQL(m)
		query, params, returning := q.SQL()
		t.Log(query)
		if query != "UPDATE test_model SET name = ?, pages = ?, some_int = ?, updated_at = NOW() WHERE (id = ?) RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong classic_update_1")
		}
		if len(params) != 4 {
			t.Fatal("wrong classic_update_1 args len must be 4")
		}
		if params[0] != m.Name {
			t.Fatal("wrong classic_update_1 args 0 must be name")
		}
		if params[3] != m.Id {
			t.Fatal("wrong classic_update_1 args 0 must be id")
		}
		if len(returning) != 3 {
			t.Fatal("wrong classic_update_1 returning arg len must be 3")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong classic_update_1 returning CreatedAt")
		}
	})
	t.Run("classic_update_2", func(t *testing.T) {
		m := &UpdateModel2{}
		m.SomeInt = new(int)
		*m.SomeInt = 1000
		q := GetUpdateSQL(m)
		query, params, returning := q.SQL()
		t.Log(query)
		if query != "UPDATE test_model SET name = ?, pages = ?, updated_at = NOW() WHERE (some_int = ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong classic_update_2")
		}
		if len(params) != 3 {
			t.Fatal("wrong classic_update_2 args len must be 4")
		}
		if params[0] != m.Name {
			t.Fatal("wrong classic_update_2 args 0 must be name")
		}
		if params[2] != m.SomeInt {
			t.Fatal("wrong classic_update_2 args 2 must be SomeInt")
		}
		if len(returning) != 4 {
			t.Fatal("wrong classic_update_2 returning arg len must be 4")
		}
		if returning[0] != &m.Id {
			t.Fatal("wrong classic_update_2 returning id")
		}
	})

	t.Run("classic_update_3", func(t *testing.T) {
		m := &InsertModel2{}
		m.Id = new(int)
		*m.Id = 1000

		m.SomeInt = new(int)
		*m.SomeInt = 1000

		m.Pages = []string{"aas"}

		m.Name = new(string)
		*m.Name = "clear"

		q := GetUpdateSQL(m)
		query, params, returning := q.SQL()
		t.Log(query)
		if query != "UPDATE test_model SET name = ?, pages = ?, updated_at = NOW() WHERE (id = ? AND some_int = ?) RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong classic_update_3")
		}
		if len(params) != 4 {
			t.Fatal("wrong classic_update_3 args len must be 4")
		}
		if params[0] != m.Name {
			t.Fatal("wrong classic_update_3 args 0 must be name")
		}
		if params[2] != m.Id {
			t.Fatal("wrong classic_update_3 args 2 must be Id")
		}
		if len(returning) != 3 {
			t.Fatal("wrong classic_update_3 returning arg len must be 3")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong classic_update_3 returning CreatedAt")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// BenchmarkGetUpdate
// BenchmarkGetUpdate-12    	 1488288	       796.1 ns/op	     688 B/op	      10 allocs/op
func BenchmarkGetUpdate(b *testing.B) {
	m := NewTestModel()
	for i := 0; i < b.N; i++ {
		GetUpdateSQL(m, &m.Name, &m.SomeInt)
	}
	b.ReportAllocs()
}
