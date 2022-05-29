package gomodel

import (
	"testing"
	"time"
)

func TestGetDeleteSQL(t *testing.T) {
	t.Run("soft", func(t *testing.T) {
		model := &InsertModel1{}
		model.Id = &ACMId
		iSql := GetDeleteSQL(model)
		if iSql == nil {
			t.Fatal("soft must be not nil")
		}
		query, params, returning := iSql.SQL()
		t.Log(query)
		if query != `UPDATE test_model SET updated_at = NOW(), deleted_at = ? WHERE (id = ?) RETURNING updated_at;` {
			t.Fatal("soft wrong query")
		}
		if len(params) != 2 {
			t.Fatal("soft must have 3 param")
		}
		if *(params[1].(**int)) != model.Id {
			t.Fatal("soft wrong param ref")
		}
		if **(params[1].(**int)) != *model.Id {
			t.Fatal("soft wrong param value")
		}
		if **(params[0].(**time.Time)) != *model.DeletedAt {
			t.Fatal("soft wrong param value")
		}
		if len(returning) != 1 {
			t.Fatal("soft must have 1 returning")
		}
	})
	t.Run("soft_unique", func(t *testing.T) {
		model := &InsertModel2{}
		model.SomeInt = &ACMSomeInt
		iSql := GetDeleteSQL(model)
		if iSql == nil {
			t.Fatal("soft must be not nil")
		}
		query, params, returning := iSql.SQL()
		t.Log(query)
		if query != `UPDATE test_model SET updated_at = NOW(), deleted_at = ? WHERE (id = ? AND some_int = ?) RETURNING updated_at;` {
			t.Fatal("soft wrong query")
		}
		if len(params) != 2 {
			t.Fatal("soft must have 3 param")
		}
		if *(params[1].(**int)) != model.Id {
			t.Fatal("soft wrong param ref")
		}
		if **(params[1].(**int)) != *model.Id {
			t.Fatal("soft wrong param value")
		}
		if **(params[0].(**time.Time)) != *model.DeletedAt {
			t.Fatal("soft wrong param value")
		}
		if len(returning) != 1 {
			t.Fatal("soft must have 1 returning")
		}
	})
	t.Run("classic", func(t *testing.T) {
		model := &DeleteModel1{}
		model.Id = &ACMId
		iSql := GetDeleteSQL(model)
		if iSql == nil {
			t.Fatal("soft must be not nil")
		}
		query, params, returning := iSql.SQL()
		t.Log(query)
		if len(params) != 1 {
			t.Fatal("classic must have 1 param")
		}
		if *(params[0].(**int)) != model.Id {
			t.Fatal("soft wrong param ref id")
		}
		if **(params[0].(**int)) != *model.Id {
			t.Fatal("soft wrong param value id")
		}
		if len(returning) != 0 {
			t.Fatal("soft must have 0 returning")
		}
	})
}

func BenchmarkName(b *testing.B) {
	b.Run("soft", func(b *testing.B) {
		model := NewTestModel()
		model.Id = &ACMId
		for i := 0; i < b.N; i++ {
			GetDeleteSQL(model)
		}
		b.ReportAllocs()
	})
	b.Run("classic", func(b *testing.B) {
		model := &InsertModel1{}
		model.Id = &ACMId
		for i := 0; i < b.N; i++ {
			GetDeleteSQL(model)
		}
		b.ReportAllocs()
	})
}
