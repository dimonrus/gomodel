package gomodel

import "testing"

func TestGetLoadSQL(t *testing.T) {
	t.Run("classic_pk", func(t *testing.T) {
		m := InsertModel1{Id: &ACMId}
		q := GetLoadSQL(&m)
		q = GetLoadSQL(&m)
		query, param, _ := q.SQL()
		t.Log(query)
		if query != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM test_model_1 WHERE (id = ? AND deleted_at IS NULL)" {
			t.Fatal("wrong sql classic_pk")
		}
		if len(param) != 1 {
			t.Fatal("classic_pk wrong param len")
		}
		if **(param[0].(**int)) != ACMId {
			t.Fatal("classic_pk wrong param[0] value")
		}
	})
	t.Run("classic_unique", func(t *testing.T) {
		m := InsertModel2{Id: &ACMId}
		q := GetLoadSQL(&m)
		q = GetLoadSQL(&m)
		query, param, _ := q.SQL()
		t.Log(query)
		if query != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM test_model_2 WHERE (id = ? AND deleted_at IS NULL)" {
			t.Fatal("wrong sql classic_unique")
		}
		if len(param) != 1 {
			t.Fatal("classic_unique wrong param len")
		}
		if **(param[0].(**int)) != ACMId {
			t.Fatal("classic_unique wrong param[0] value")
		}
	})
	t.Run("classic_string_prk", func(t *testing.T) {
		someStringId := "sdsdf-12312"
		m := InsertModel3{Id: &someStringId}
		q := GetLoadSQL(&m)
		q = GetLoadSQL(&m)
		query, param, _ := q.SQL()
		t.Log(query)
		if query != "SELECT id, name, pages, some_int FROM test_model_3 WHERE (id = ?)" {
			t.Fatal("wrong sql classic_string_prk")
		}
		if len(param) != 1 {
			t.Fatal("classic_string_prk wrong param len")
		}
		if **param[0].(**string) != someStringId {
			t.Fatal("classic_string_prk wrong param[0] value")
		}
	})
	t.Run("classic_unique_2", func(t *testing.T) {
		m := UpdateModel2{SomeInt: &ACMSomeInt}
		q := GetLoadSQL(&m)
		q = GetLoadSQL(&m)
		query, param, _ := q.SQL()
		t.Log(query)
		if query != "SELECT id, name, pages, some_int, created_at, updated_at, deleted_at FROM test_model_upd_2 WHERE (some_int = ? AND deleted_at IS NULL)" {
			t.Fatal("wrong sql classic_unique_2")
		}
		if len(param) != 1 {
			t.Fatal("classic_unique_2 wrong param len")
		}
		if **param[0].(**int) != ACMSomeInt {
			t.Fatal("classic_unique_2 wrong param[0] value")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkGetLoadSQL
// BenchmarkGetLoadSQL-12    	 7622672	       150.4 ns/op	     304 B/op	       4 allocs/op
func BenchmarkGetLoadSQL(b *testing.B) {
	m := &InsertModel1{Id: &ACMId}
	for i := 0; i < b.N; i++ {
		GetLoadSQL(m)
	}
	b.ReportAllocs()
}
