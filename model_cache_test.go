package gomodel

import (
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("get store", func(t *testing.T) {
		db, err := initDb()
		if err != nil {
			t.Fatal(err)
		}
		collection := NewBigTestTableCollection()
		collection.SetPagination(100, 0)
		e := collection.Load(db)
		if e != nil {
			t.Fatal(e)
		}
		model := NewBigTestTable()
		model.Id = collection.Last().Id
		o := IndexCache.Get(IndexOperationLoad, model)
		if o != nil {
			t.Fatal("must be nil")
		}
		isql := GetLoadSQL(model)
		_ = isql

		e = Load(db, model)
		if e != nil {
			t.Fatal(e)
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkCache_Key
// BenchmarkCache_Key-12    	37482244	        31.87 ns/op	      24 B/op	       1 allocs/op
func BenchmarkCache_Key(b *testing.B) {
	m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
	for i := 0; i < b.N; i++ {
		IndexCache.Get(IndexOperationCreate, &m)
	}
	b.ReportAllocs()
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gomodel
// cpu: Apple M2 Max
// BenchmarkCache_KeyValue
// BenchmarkCache_KeyValue-12    	16706208	        63.15 ns/op	     112 B/op	       1 allocs/op
func BenchmarkCache_KeyValue(b *testing.B) {
	m := InsertModel1{Name: &ACMName, SomeInt: &ACMSomeInt}
	for i := 0; i < b.N; i++ {
		IndexCache.Get(IndexOperationCreate, &m, &m.Name)
	}
	b.ReportAllocs()
}
