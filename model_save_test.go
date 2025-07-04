package gomodel

import (
	"github.com/lib/pq"
	"testing"
)

func TestGetSaveScenario(t *testing.T) {
	t.Run("nil_model", func(t *testing.T) {
		var m IModel
		insert, update, upsert := getSaveScenario(m)
		if insert || update || upsert {
			t.Fatal("wrong nil_model")
		}
	})
	t.Run("insert_1", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if !insert || update || upsert {
			t.Fatal("wrong insert_1")
		}
	})
	t.Run("insert_2", func(t *testing.T) {
		m := InsertModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if !insert || update || upsert {
			t.Fatal("wrong insert_2")
		}
	})
	t.Run("upsert_1", func(t *testing.T) {
		m := UpsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if insert || update || !upsert {
			t.Fatal("wrong upsert_1")
		}
	})
	t.Run("upsert_2", func(t *testing.T) {
		m := UpsertModel2{Id: &ACMId, Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if insert || update || !upsert {
			t.Fatal("wrong upsert_2")
		}
	})
	t.Run("upsert_3", func(t *testing.T) {
		m := UpsertModel3{Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if insert || update || !upsert {
			t.Fatal("wrong upsert_3")
		}
	})
	t.Run("upsert_4", func(t *testing.T) {
		m := UpsertModel4{Name: &ACMName}
		insert, update, upsert := getSaveScenario(&m)
		if insert || update || !upsert {
			t.Fatal("wrong upsert_4")
		}
	})
	t.Run("update_1", func(t *testing.T) {
		m := UpdateModel1{Id: &ACMId, Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if insert || !update || upsert {
			t.Fatal("wrong update_1")
		}
	})
	t.Run("update_2", func(t *testing.T) {
		m := UpdateModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
		insert, update, upsert := getSaveScenario(&m)
		if insert || !update || upsert {
			t.Fatal("wrong update_2")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkGetSaveScenario
// BenchmarkGetSaveScenario-12    	 2134198	       554.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGetSaveScenario(b *testing.B) {
	m := UpdateModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
	for i := 0; i < b.N; i++ {
		getSaveScenario(&m)
	}
	b.ReportAllocs()
}

func TestGetSaveSQL(t *testing.T) {
	t.Run("insert_1", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_1 (name, pages, some_int) VALUES (?, ?, ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_1 sql")
		}
		if len(param) != 3 {
			t.Fatal("wrong insert_1 param len")
		}
		if len(returning) != 4 {
			t.Fatal("wrong insert_1 returning len")
		}
		if *param[0].(**string) != m.Name {
			t.Fatal("wrong insert_1 name addr")
		}
		if *param[2].(**int) != m.SomeInt {
			t.Fatal("wrong insert_1 some int addr")
		}
		if returning[0] != &m.Id {
			t.Fatal("wrong insert_1 returning addr")
		}
	})
	t.Run("insert_2", func(t *testing.T) {
		m := InsertModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_2 (name, pages) VALUES (?, ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_2 sql")
		}
		if len(param) != 2 {
			t.Fatal("wrong insert_2 param len")
		}
		if len(returning) != 4 {
			t.Fatal("wrong insert_2 returning len")
		}
		if *param[0].(**string) != m.Name {
			t.Fatal("wrong insert_2 name addr")
		}
		if *param[1].(*pq.StringArray) != nil {
			t.Fatal("wrong insert_2 int addr")
		}
		if returning[0] != &m.Id {
			t.Fatal("wrong insert_2 returning addr")
		}
	})
	t.Run("upsert_1", func(t *testing.T) {
		m := UpsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_1 (name, pages, some_int) VALUES (?, ?, ?) ON CONFLICT (some_int) DO UPDATE SET name = ?, pages = ?, updated_at = NOW() RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_1 sql")
		}
		if len(param) != 5 {
			t.Fatal("wrong upsert_1 param len")
		}
		if len(returning) != 4 {
			t.Fatal("wrong upsert_1 returning len")
		}
		if *param[0].(**string) != m.Name {
			t.Fatal("wrong upsert_1 name addr")
		}
		if *param[3].(**string) != m.Name {
			t.Fatal("wrong upsert_1 name addr")
		}
		if returning[0] != &m.Id {
			t.Fatal("wrong upsert_1 returning addr")
		}
	})
	t.Run("upsert_2", func(t *testing.T) {
		m := UpsertModel2{Id: &ACMId, Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_2 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET name = ?, pages = ?, some_int = ?, updated_at = NOW() RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_2 sql")
		}
		if len(param) != 7 {
			t.Fatal("wrong upsert_2 param len")
		}
		if len(returning) != 3 {
			t.Fatal("wrong upsert_2 returning len")
		}
		if *param[0].(**int) != m.Id {
			t.Fatal("wrong upsert_2 name addr")
		}
		if *param[1].(**string) != m.Name {
			t.Fatal("wrong upsert_2 name addr")
		}
		if **param[1].(**string) != *m.Name {
			t.Fatal("wrong upsert_2 name val")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong upsert_2 returning addr")
		}
	})
	t.Run("upsert_3", func(t *testing.T) {
		m := UpsertModel3{Name: &ACMName, Pages: []string{"one", "two"}, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_3 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET name = ?, pages = ?, some_int = ?, updated_at = NOW() RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_3 sql")
		}
		if len(param) != 7 {
			t.Fatal("wrong upsert_3 param len")
		}
		if len(returning) != 3 {
			t.Fatal("wrong upsert_3 returning len")
		}
		if *param[0].(**int) != m.Id {
			t.Fatal("wrong upsert_3 id addr")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong upsert_3 returning addr")
		}
	})
	t.Run("upsert_4", func(t *testing.T) {
		m := UpsertModel4{Name: &ACMName}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_4 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (some_int) DO UPDATE SET id = ?, name = ?, pages = ?, updated_at = NOW() RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_3 sql")
		}
		if len(param) != 7 {
			t.Fatal("wrong upsert_3 param len")
		}
		if len(returning) != 3 {
			t.Fatal("wrong upsert_3 returning len")
		}
		if *param[0].(**int) != m.Id {
			t.Fatal("wrong upsert_3 name addr")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong upsert_3 returning addr")
		}
	})
	t.Run("upsert_5", func(t *testing.T) {
		complexId := 121
		m := UpsertModel5{Id: &ACMId, ComplexId: &complexId, Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO test_model_up_5 (id, complex_id, name, pages, some_int) VALUES (?, ?, ?, ?, ?) ON CONFLICT (id, complex_id) DO UPDATE SET name = ?, pages = ?, some_int = ?, updated_at = NOW() RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_5 sql")
		}
		if len(param) != 8 {
			t.Fatal("wrong upsert_5 param len")
		}
		if len(returning) != 3 {
			t.Fatal("wrong upsert_5 returning len")
		}
		if *param[0].(**int) != m.Id {
			t.Fatal("wrong upsert_5 Id addr")
		}
		if *param[1].(**int) != m.ComplexId {
			t.Fatal("wrong upsert_5 ComplexId addr")
		}
		if *param[7].(**int) != m.SomeInt {
			t.Fatal("wrong upsert_5 SomeInt addr")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong upsert_5 returning CreatedAt addr")
		}
	})
	t.Run("update_1", func(t *testing.T) {
		m := UpdateModel1{Id: &ACMId, Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "UPDATE test_model_upd_1 SET name = ?, pages = ?, some_int = ?, updated_at = NOW() WHERE (id = ?) RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong update_1 sql")
		}
		if len(param) != 4 {
			t.Fatal("wrong update_1 param len")
		}
		if len(returning) != 3 {
			t.Fatal("wrong update_1 returning len")
		}
		if *param[0].(**string) != m.Name {
			t.Fatal("wrong update_1 name addr")
		}
		if *param[3].(**int) != m.Id {
			t.Fatal("wrong update_1 name addr")
		}
		if returning[0] != &m.CreatedAt {
			t.Fatal("wrong update_1 returning addr")
		}
	})
	t.Run("update_2", func(t *testing.T) {
		m := UpdateModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetSaveSQL(&m)
		iSql = GetSaveSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "UPDATE test_model_upd_2 SET name = ?, pages = ?, updated_at = NOW() WHERE (some_int = ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong update_2 sql")
		}
		if len(param) != 3 {
			t.Fatal("wrong update_2 param len")
		}
		if len(returning) != 4 {
			t.Fatal("wrong update_2 returning len")
		}
		if *param[0].(**string) != m.Name {
			t.Fatal("wrong update_2 name addr")
		}
		if *param[2].(**int) != m.SomeInt {
			t.Fatal("wrong update_2 name addr")
		}
		if returning[0] != &m.Id {
			t.Fatal("wrong update_2 returning addr")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkGetSaveSQL
// BenchmarkGetSaveSQL-12    	 1777431	       666.5 ns/op	     288 B/op	       4 allocs/op
func BenchmarkGetSaveSQL(b *testing.B) {
	m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
	for i := 0; i < b.N; i++ {
		GetSaveSQL(&m)
	}
	b.ReportAllocs()
}
