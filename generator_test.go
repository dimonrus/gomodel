package gomodel

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"testing"
)

type connection struct{}

func (c *connection) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "gomodel", "gomodel", "gomodel")
}
func (c *connection) GetDbType() string       { return "postgres" }
func (c *connection) GetMaxConnection() int   { return 200 }
func (c *connection) GetMaxIdleConns() int    { return 15 }
func (c *connection) GetConnMaxLifetime() int { return 50 }

func TestDBO_InitError(t *testing.T) {
	_, err := godb.DBO{
		Options: godb.Options{
			Debug:  true,
			Logger: gocli.NewLogger(gocli.LoggerConfig{}),
		},
		Connection: &connection{},
	}.Init()
	if err == nil {
		t.Fatal("must be an error case")
	}
}

func initDb() (*godb.DBO, error) {
	return godb.DBO{
		Options: godb.Options{
			Debug:  true,
			Logger: gocli.NewLogger(gocli.LoggerConfig{}),
		},
		Connection: &connection{},
	}.Init()
}

func TestMakeModel(t *testing.T) {
	db, err := initDb()
	db.Debug = false
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = MakeModel(db, "models", "public", "dictionary", "", DefaultSystemColumnsSoft)
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
	crud := NewCRUDGenerator("app/core", "app/client", "app/io/web/api", "github.com/dimonrus/gomodel")
	err = crud.Generate(db, "public", "dictionary", "v2", 31)
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

func TestCreateDictionaryTable(t *testing.T) {
	list := getDictionarySQList()
	query, _, _ := list.Join()
	t.Log(query)
	if query != `CREATE TABLE IF NOT EXISTS dictionary (id INT NOT NULL PRIMARY KEY, type TEXT NOT NULL, code TEXT NOT NULL, label TEXT, created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT localtimestamp, updated_at TIMESTAMP WITH TIME ZONE, deleted_at TIMESTAMP WITH TIME ZONE);COMMENT ON COLUMN dictionary.id IS 'Dictionary row identifier';COMMENT ON COLUMN dictionary.type IS 'Dictionary row type';COMMENT ON COLUMN dictionary.code IS 'Dictionary row code';COMMENT ON COLUMN dictionary.label IS 'Dictionary row value label';COMMENT ON COLUMN dictionary.created_at IS 'Dictionary row created time';COMMENT ON COLUMN dictionary.updated_at IS 'Dictionary row updated time';COMMENT ON COLUMN dictionary.deleted_at IS 'Dictionary row deleted time';CREATE INDEX IF NOT EXISTS dictionary_type_idx ON dictionary (type);` {
		t.Fatal("wrong dictionary query")
	}
	db, err := initDb()
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
