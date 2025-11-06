package gomodel

import (
	"database/sql"
	"testing"
	"time"

	"github.com/dimonrus/gohelp"
	"github.com/lib/pq"
)

func TestNewWizard(t *testing.T) {
	t.Run("new wizard", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
		if wiz.IsComplete() {
			t.Fatal("must be incomplete")
		}
	})
	t.Run("new wizard empty order", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]()
		if wiz.IsComplete() {
			t.Fatal("must be incomplete")
		}
	})
	t.Run("incorrect type", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name")
		if wiz.IsComplete() {
			t.Fatal("must be incomplete")
		}
		_, e := wiz.Set(123)
		if e == nil {
			t.Fatal("must be an error")
		}
		if wiz.IsComplete() {
			t.Fatal("must be incomplete")
		}
	})
	t.Run("new wizard panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				if err == nil {
					t.Fatal("wrong model")
				} else if err.Error() != "interface conversion: *gomodel.MetaModel is not gomodel.IModel: missing method Columns" {
					t.Fatal(err.Error())
				}
			}
		}()
		wiz := NewWizard[MetaModel]("name", "pages", "some_int")
		t.Log(wiz.IsComplete())
	})
}

func TestWizard_Set(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
		_, e := wiz.Set(&sql.NullString{String: "foo"})
		if e != nil {
			t.Fatal(e)
		}
		_, e = wiz.Set(&pq.StringArray{"foo", "bar"})
		if e != nil {
			t.Fatal(e)
		}
		val := gohelp.Ptr[float64](123)
		_, e = wiz.Set(&val)
		if e != nil {
			t.Fatal(e)
		}
		model := wiz.Get()
		if model.Name.String != "foo" {
			t.Fatal("name must be foo")
		}
		if model.Pages[0] != "foo" || model.Pages[1] != "bar" {
			t.Fatal("pages is incorrect")
		}
		if model.SomeInt == nil || *model.SomeInt != 123 {
			t.Fatal("some int must be 123")
		}
		if !wiz.IsComplete() {
			t.Fatal("must be complete")
		}
	})
	t.Run("set values more the 3", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
		_, e := wiz.Set(&sql.NullString{String: "foo"})
		if e != nil {
			t.Fatal(e)
		}
		_, e = wiz.Set(&pq.StringArray{"foo", "bar"})
		if e != nil {
			t.Fatal(e)
		}
		val := gohelp.Ptr[float64](123)
		_, e = wiz.Set(&val)
		if e != nil {
			t.Fatal(e)
		}
		tt := time.Now()
		_, e = wiz.Set(&tt)
		if e == nil {
			t.Fatal("must be overflow value error")
		}
	})
	t.Run("all model fields", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("id", "name", "pages", "some_int", "created_at", "custom")
		if wiz.GetCursor() != 0 {
			t.Fatal("must be 0")
		}
		_, e := wiz.Set(gohelp.Ptr(12))
		if e != nil {
			t.Fatal(e)
		}
		if wiz.GetCursor() != 1 {
			t.Fatal("must be 0")
		}
		_, e = wiz.Set(&sql.NullString{String: "foo"})
		if e != nil {
			t.Fatal(e)
		}
		if wiz.GetCursor() != 2 {
			t.Fatal("must be 2")
		}
		_, e = wiz.Set(&pq.StringArray{"foo", "bar"})
		if e != nil {
			t.Fatal(e)
		}
		if wiz.GetCursor() != 3 {
			t.Fatal("must be 3")
		}
		val := gohelp.Ptr[float64](123)
		_, e = wiz.Set(&val)
		if e != nil {
			t.Fatal(e)
		}
		if wiz.GetCursor() != 4 {
			t.Fatal("must be 4")
		}
		tt := time.Now()
		_, e = wiz.Set(&tt)
		if e != nil {
			t.Fatal(e)
		}
		if wiz.GetCursor() != 5 {
			t.Fatal("must be 5")
		}
		custom := &struct{ Foo int }{Foo: 1234}
		next, e := wiz.Set(&custom)
		if e != nil {
			t.Fatal(e)
		}
		if next {
			t.Fatal("no next expected")
		}
		if wiz.GetCursor() != 6 {
			t.Fatal("must be 6")
		}
		model := wiz.Get()
		if model.Custom == nil || model.Custom.Foo != 1234 {
			t.Fatal("wrong custom")
		}
		if !wiz.IsComplete() {
			t.Fatal("must be completed")
		}
	})
	t.Run("empty", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("id", "name", "pages", "some_int", "created_at", "custom")
		if wiz.IsComplete() {
			t.Fatal("must not be completed")
		}
		_, e := wiz.Set(gohelp.Ptr(12))
		if e != nil {
			t.Fatal(e)
		}
		_, e = wiz.Set(&sql.NullString{String: "foo"})
		if e != nil {
			t.Fatal(e)
		}
		wiz.ResetCursor()
		_, e = wiz.Set(gohelp.Ptr(341))
		if e != nil {
			t.Fatal(e)
		}
		model := wiz.Get()
		if model.Id != 341 {
			t.Fatal("must be 431")
		}
	})
}

func BenchmarkWizardGet(b *testing.B) {
	wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
	_, e := wiz.Set(&sql.NullString{String: "foo"})
	if e != nil {
		b.Fatal(e)
	}
	_, e = wiz.Set(&pq.StringArray{"foo", "bar"})
	if e != nil {
		b.Fatal(e)
	}
	val := gohelp.Ptr[float64](123)
	_, e = wiz.Set(&val)
	if e != nil {
		b.Fatal(e)
	}
	for i := 0; i < b.N; i++ {
		wiz.Get()
	}
	b.ReportAllocs()
}

func TestWizard_SetByName(t *testing.T) {
	t.Run("set by name", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
		num := gohelp.Ptr[float64](123)
		e := wiz.SetByName("some_int", &num)
		if e != nil {
			t.Fatal(e)
		}
		model := wiz.Get()
		if model.SomeInt == nil || *model.SomeInt != 123 {
			t.Fatal("wrong")
		}
		if wiz.IsComplete() {
			t.Fatal("must not be completed")
		}
	})

}

func TestWizard_Next(t *testing.T) {
	t.Run("next", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
		var cursor uint16
		if cursor != 0 {
			t.Fatal("must be 1")
		}
		hasNext := wiz.Next()
		cursor = wiz.GetCursor()
		if cursor != 1 {
			t.Fatal("must be 1")
		}
		if !hasNext {
			t.Fatal("next required")
		}
		hasNext = wiz.Next()
		cursor = wiz.GetCursor()
		if cursor != 2 {
			t.Fatal("must be 2")
		}
		if hasNext {
			t.Fatal("no more next")
		}
	})
	t.Run("by name", func(t *testing.T) {
		wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
		var cursor uint16
		if cursor != 0 {
			t.Fatal("must be 1")
		}
		num := gohelp.Ptr(123)
		e := wiz.SetByName("some_int", &num)
		if e != nil {
			t.Fatal("no error expected")
		}
		if wiz.GetCursor() != 2 {
			t.Fatal("must be 2")
		}
		if wiz.IsComplete() {
			t.Fatal("must not be completed")
		}
	})

}
