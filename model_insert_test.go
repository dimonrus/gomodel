package gomodel

import (
	"testing"
)

func TestGetInsertSQL(t *testing.T) {
	t.Run("insert_1", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO insert_model1 (name, pages, some_int) VALUES (?, ?, ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_1 sql")
		}
		if len(param) != 3 {
			t.Fatal("insert_1 wrong param len")
		}
		if *param[0].(*string) != ACMName {
			t.Fatal("insert_1 wrong param[0] value")
		}
		if *param[2].(*int) != ACMSomeInt {
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
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO insert_model1 (name, pages, some_int) VALUES (?, ?, ?) RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_1 sql")
		}
		if len(param) != 3 {
			t.Fatal("insert_1 wrong param len")
		}
		if *param[0].(*string) != ACMName {
			t.Fatal("insert_1 wrong param[0] value")
		}
		if *param[2].(*int) != ACMSomeInt {
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
	t.Run("insert_1_manual", func(t *testing.T) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m, &m.Name)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO insert_model1 (name) VALUES (?);" {
			t.Fatal("wrong insert_1_manual sql")
		}
		if len(param) != 1 {
			t.Fatal("insert_1_manual wrong param len")
		}
		if *param[0].(*string) != ACMName {
			t.Fatal("insert_1_manual wrong param[0] value")
		}
	})
	t.Run("insert_2", func(t *testing.T) {
		m := InsertModel2{Name: &ACMName, SomeInt: &ACMSomeInt}
		iSql := GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO insert_model2 (name, pages, some_int) VALUES (?, ?, ?) ON CONFLICT (some_int) DO UPDATE SET name = ?, pages = ? RETURNING id, created_at, updated_at, deleted_at;" {
			t.Fatal("wrong insert_2 sql")
		}
		if len(param) != 5 {
			t.Fatal("insert_2 wrong param len")
		}
		if *param[0].(*string) != ACMName {
			t.Fatal("insert_2 wrong param[0] value")
		}
		if *param[2].(*int) != ACMSomeInt {
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
	t.Run("upsert_2", func(t *testing.T) {
		m := UpsertModel2{Id: &ACMId, Name: &ACMName}
		iSql := GetInsertSQL(&m)
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO upsert_model2 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET name = ?, pages = ?, some_int = ? RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_2 sql")
		}
		if len(param) != 7 {
			t.Fatal("upsert_2 wrong param len")
		}
		if *param[0].(*int) != ACMId {
			t.Fatal("upsert_2 wrong param[0] value")
		}
		if *param[1].(*string) != ACMName {
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
		query, param, returning := iSql.SQL()
		t.Log(query, param, returning)
		if query != "INSERT INTO upsert_model4 (id, name, pages, some_int) VALUES (?, ?, ?, ?) ON CONFLICT (some_int) DO UPDATE SET id = ?, name = ?, pages = ? RETURNING created_at, updated_at, deleted_at;" {
			t.Fatal("wrong upsert_2 sql")
		}
		if len(param) != 7 {
			t.Fatal("upsert_4 wrong param len")
		}
		if *param[0].(*int) != ACMId {
			t.Fatal("upsert_4 wrong param[0] value")
		}
		if *param[1].(*string) != ACMName {
			t.Fatal("upsert_4 wrong param[1] value")
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
}

func BenchmarkGetInsertSQL(b *testing.B) {
	b.Run("insert_sql", func(b *testing.B) {
		m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
		for i := 0; i < b.N; i++ {
			GetInsertSQL(&m)
		}
		b.ReportAllocs()
	})
}
