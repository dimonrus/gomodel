package generator

import (
	"testing"
)

func TestMakeModel(t *testing.T) {
	db, err := initDb()
	db.Debug = false
	if err != nil {
		t.Fatal(err)
	}
	err = MakeModel(db, "models", "public", "login", "", DefaultSystemColumnsSoft)
	if err != nil {
		t.Fatal(err)
	}
}
