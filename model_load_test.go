package gomodel

import "testing"

func TestGetLoadSQL(t *testing.T) {
	t.Run("classic_pk", func(t *testing.T) {
		m := InsertModel1{Id: &ACMId}
		q := GetLoadSQL(&m)
		query, _, _ := q.SQL()
		t.Log(query)

		if query != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM test_model WHERE (id = ? AND deleted_at IS NOT NULL)" {
			t.Fatal("wrong sql classic_pk")
		}
	})
	t.Run("classic_unique", func(t *testing.T) {
		m := InsertModel2{Id: &ACMId}
		q := GetLoadSQL(&m)
		query, _, _ := q.SQL()
		t.Log(query)

		if query != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM test_model WHERE (id = ? AND deleted_at IS NOT NULL)" {
			t.Fatal("wrong sql classic_unique")
		}
	})
	t.Run("classic_unique_2", func(t *testing.T) {
		m := UpdateModel2{SomeInt: &ACMSomeInt}
		q := GetLoadSQL(&m)
		query, _, _ := q.SQL()
		t.Log(query)

		if query != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM test_model WHERE (some_int = ? AND deleted_at IS NOT NULL)" {
			t.Fatal("wrong sql classic_unique_2")
		}
	})
}

func BenchmarkGetLoadSQL(b *testing.B) {
	m := &InsertModel1{Id: &ACMId}
	for i := 0; i < b.N; i++ {
		GetLoadSQL(m)
	}
	b.ReportAllocs()
}
