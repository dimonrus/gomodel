package gomodel

import (
	"database/sql"
	"github.com/dimonrus/gohelp"
	"github.com/lib/pq"
	"testing"
)

func TestGetInsertSQL(t *testing.T) {
	t.Run("insert_1", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_1 (name, pages, some_int) VALUES (?, ?, ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_1 sql")
		}
		if len(param) != 3 {
			t.Fatal("insert_1 wrong param len")
		}
		if **(param[0].(**string)) != ACMName {
			t.Fatal("insert_1 wrong param[0] value")
		}
		if **(param[2].(**int)) != ACMSomeInt {
			t.Fatal("insert_1 wrong param[2] value")
		}
		if len(returning) != 4 {
			t.Fatal("insert_1 wrong returning len")
		}
		if returning[0] != &m.Id {
			t.Fatal("insert_1 wrong returning[0] param")
		}
		if returning[3] != &m.DeletedAt {
			t.Fatal("insert_1 wrong returning[3] param")
		}
	})
	t.Run("insert_1", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_1 (name, pages, some_int) VALUES (?, ?, ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_1 sql")
		}
		if len(param) != 3 {
			t.Fatal("insert_1 wrong param len")
		}
		if **(param[0].(**string)) != ACMName {
			t.Fatal("insert_1 wrong param[0] value")
		}
		if **(param[2].(**int)) != ACMSomeInt {
			t.Fatal("insert_1 wrong param[2] value")
		}
		if len(returning) != 4 {
			t.Fatal("insert_1 wrong returning len")
		}
		if returning[0] != &m.Id {
			t.Fatal("insert_1 wrong returning[0] param")
		}
		if returning[3] != &m.DeletedAt {
			t.Fatal("insert_1 wrong returning[3] param")
		}
	})
	t.Run("insert_1_manual_vs_insert_1", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query1, _, _ := iSql.SQL()
		iSql = GetInsertSQL(&m, &m.Name)
		iSql = GetInsertSQL(&m, &m.Name)
		query2, _, _ := iSql.SQL()
		if query1 == query2 {
			t.Fatal("must be other insert sql")
		}
		t.Log(query1)
		t.Log(query2)
	})
	t.Run("insert_1_manual", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m, &m.Name)
		iSql = GetInsertSQL(&m, &m.Name)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_1 (name) VALUES (?);" {
			t.Fatal("wrong insert_1_manual sql")
		}
		if len(param) != 1 {
			t.Fatal("insert_1_manual wrong param len")
		}
		if **param[0].(**string) != ACMName {
			t.Fatal("insert_1_manual wrong param[0] value")
		}
	})
	t.Run("insert_2", func(t *testing.T) {
		m := InsertModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_2 (name, pages, some_int) VALUES (?, ?, ?) ON CONFLICT (some_int) DO UPDATE SET name = ?, pages = ? RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_2 sql")
		}
		if len(param) != 5 {
			t.Fatal("insert_2 wrong param len")
		}
		if **param[0].(**string) != ACMName {
			t.Fatal("insert_2 wrong param[0] value")
		}
		if **param[2].(**int) != ACMSomeInt {
			t.Fatal("insert_2 wrong param[2] value")
		}
		if len(returning) != 4 {
			t.Fatal("insert_2 wrong returning len")
		}
		if returning[0] != &m.Id {
			t.Fatal("insert_2 wrong returning[0] param")
		}
		if returning[3] != &m.DeletedAt {
			t.Fatal("insert_2 wrong returning[3] param")
		}
	})
	t.Run("insert_3", func(t *testing.T) {
		id := "some-12312"
		m := InsertModel3{Id: &id, Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_3 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET name = ?, pages = ?, some_int = ?;" {
			t.Fatal("wrong insert_3 sql")
		}
		if len(param) != 7 {
			t.Fatal("insert_3 wrong param len")
		}
		if **param[0].(**string) != id {
			t.Fatal("insert_3 wrong param[0] value")
		}
		if **param[3].(**int) != ACMSomeInt {
			t.Fatal("insert_3 wrong param[2] value")
		}
		if len(returning) != 0 {
			t.Fatal("insert_3 wrong returning len")
		}
	})
	t.Run("upsert_2", func(t *testing.T) {
		m := UpsertModel2{Id: &ACMId, Name: &ACMName}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_2 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET name = ?, pages = ?, some_int = ? RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_2 sql")
		}
		if len(param) != 7 {
			t.Fatal("upsert_2 wrong param len")
		}
		if **param[0].(**int) != ACMId {
			t.Fatal("upsert_2 wrong param[0] value")
		}
		if **param[1].(**string) != ACMName {
			t.Fatal("upsert_2 wrong param[1] value")
		}
		if len(returning) != 3 {
			t.Fatal("upsert_2 wrong returning len")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("upsert_2 wrong returning[0] param")
		}
		if returning[2] != &m.DeletedAt {
			t.Fatal("upsert_2 wrong returning[2] param")
		}
	})
	t.Run("upsert_4", func(t *testing.T) {
		m := UpsertModel4{Id: &ACMId, Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_4 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (some_int) DO UPDATE SET id = ?, name = ?, pages = ? RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_2 sql")
		}
		if len(param) != 7 {
			t.Fatal("upsert_4 wrong param len")
		}
		if **param[0].(**int) != ACMId {
			t.Fatal("upsert_4 wrong param[0] value")
		}
		if **param[1].(**string) != ACMName {
			t.Fatal("upsert_4 wrong param[1] value")
		}
		if **param[4].(**int) != ACMId {
			t.Fatal("upsert_4 wrong param[5] value")
		}
		if len(returning) != 3 {
			t.Fatal("upsert_4 wrong returning len")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("upsert_4 wrong returning[0] param")
		}
		if returning[2] != &m.DeletedAt {
			t.Fatal("upsert_4 wrong returning[2] param")
		}
	})
	t.Run("complex_model", func(t *testing.T) {
		m := ComplexTestModel{ComplexId: 1000, CategoryId: gohelp.Ptr(100), Name: sql.NullString{String: ACMName}, SomeInt: gohelp.Ptr(12121.22)}
		iSql := GetInsertSQL(&m)
		iSql = GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO complex_test_update_model (pages, name, complex_id, category_id, some_int, custom) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (complex_id, category_id) DO UPDATE SET pages = ?, name = ?, some_int = ?, custom = ? RETURNING created_at;" {
			t.Fatal("wrong complex_model sql")
		}
		if len(param) != 10 {
			t.Fatal("complex_model wrong param len")
		}
		if *param[0].(*pq.StringArray) != nil {
			t.Fatal("complex_model wrong param[0] value")
		}
		if *param[2].(*int) != 1000 {
			t.Fatal("complex_model wrong param[1] value")
		}
		if **param[8].(**float64) != 12121.22 {
			t.Fatal("complex_model wrong param[1] value")
		}
		if len(returning) != 1 {
			t.Fatal("complex_model wrong returning len")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("complex_model wrong returning[0] param")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkGetInsertSQL
// BenchmarkGetInsertSQL/insert_sql
// BenchmarkGetInsertSQL/insert_sql-12         	 7210249	       145.6 ns/op	     288 B/op	       4 allocs/op
func BenchmarkGetInsertSQL(b *testing.B) {
	b.Run("insert_sql", func(b *testing.B) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		for i := 0; i < b.N; i++ {
			GetInsertSQL(&m)
		}
		b.ReportAllocs()
	})
}
