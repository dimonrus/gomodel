package gomodel

import (
	"testing"
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
		if query != `UPDATE test_model SET updated_at = NOW(), deleted_at = NOW() WHERE (id = ?) RETURNING updated_at, deleted_at;` {
			t.Fatal("soft wrong query")
		}
		if len(params) != 1 {
			t.Fatal("soft must have 1 param")
		}
		if *(params[0].(**int)) != model.Id {
			t.Fatal("soft wrong param ref")
		}
		if **(params[0].(**int)) != *model.Id {
			t.Fatal("soft wrong param value")
		}
		if len(returning) != 2 {
			t.Fatal("soft must have 2 returning")
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
		if query != `UPDATE test_model SET updated_at = NOW(), deleted_at = NOW() WHERE (some_int = ?) RETURNING updated_at, deleted_at;` {
			t.Fatal("soft wrong query")
		}
		if len(params) != 1 {
			t.Fatal("soft must have 2 param")
		}
		if *(params[0].(**int)) != model.SomeInt {
			t.Fatal("soft wrong param ref")
		}
		if **(params[0].(**int)) != *model.SomeInt {
			t.Fatal("soft wrong param value")
		}
		if len(returning) != 2 {
			t.Fatal("soft must have 2 returning")
		}
	})
	t.Run("classic", func(t *testing.T) {
		model := &DeleteModel1{}
		model.Id = &ACMId
		iSql := GetDeleteSQL(model)
		if iSql == nil {
			t.Fatal("classic must be not nil")
		}
		query, params, returning := iSql.SQL()
		t.Log(query)
		if query != "DELETE FROM test_model WHERE (id = ?);" {
			t.Fatal("classic wrong query")
		}
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
	t.Run("classic_unique", func(t *testing.T) {
		model := &DeleteModel2{}
		model.SomeInt = &ACMSomeInt
		iSql := GetDeleteSQL(model)
		if iSql == nil {
			t.Fatal("classic_unique must be not nil")
		}
		query, params, returning := iSql.SQL()
		t.Log(query)
		if query != "DELETE FROM test_model WHERE (some_int = ?);" {
			t.Fatal("classic_unique wrong query")
		}
		if len(params) != 1 {
			t.Fatal("classic must have 1 param")
		}
		if *(params[0].(**int)) != model.SomeInt {
			t.Fatal("soft wrong param ref id")
		}
		if **(params[0].(**int)) != *model.SomeInt {
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
