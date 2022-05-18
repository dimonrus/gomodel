package gomodel

import "testing"

func TestGetLoadSQL(t *testing.T) {
	t.Run("classic_pk", func(t *testing.T) {
		m := InsertModel1{Id: &ACMId}
		q := GetLoadSQL(&m)
		t.Log(q.String())

		if q.String() != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM insert_model1 WHERE (id = ?)" {
			t.Fatal("wrong sql classic_pk")
		}
	})
	t.Run("classic_unique", func(t *testing.T) {
		m := InsertModel2{Id: &ACMId}
		q := GetLoadSQL(&m)
		t.Log(q.String())

		if q.String() != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM insert_model2 WHERE (id = ?)" {
			t.Fatal("wrong sql classic_unique")
		}
	})
	t.Run("classic_unique_2", func(t *testing.T) {
		m := UpdateModel2{SomeInt: &ACMSomeInt}
		q := GetLoadSQL(&m)
		t.Log(q.String())

		if q.String() != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM update_model2 WHERE (some_int = ?)" {
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
