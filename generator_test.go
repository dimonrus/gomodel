package gomodel

import (
	"testing"
)

func TestMakeModel(t *testing.T) {
	db, err := initDb()
	db.Debug = false
	if err != nil {
		t.Fatal(err)
	}
	err = MakeModel(db, "models", "public", "reset_password", "", DefaultSystemColumnsSoft)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCrud(t *testing.T) {
	db, err := initDb()
	db.Debug = false
	if err != nil {
		t.Fatal(err)
	}
	err = GenerateCrud("core", "client", "public", "login", "", db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateDictionaryMapping(t *testing.T) {
	db, _ := initDb()
	e := GenerateDictionaryMapping("models/dictionary_mapping.go", db)
	if e != nil {
		t.Fatal(e)
	}
}

func TestDictionaryUtils(t *testing.T) {
	db, _ := initDb()
	items := getDictionaryItems(db)
	if _, ok := items.IsDictionaryColumn("login_type_id"); !ok {
		t.Fatal("must be a dictionary item")
	}
	if _, ok := items.IsDictionaryColumn("some_new_column"); ok {
		t.Fatal("must not be a dictionary item")
	}
}
