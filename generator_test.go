package gomodel

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
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

var testTable = func(schema string) gosql.SQList {
	table := SerialTable(schema + "test_table")
	table.AddColumn("name").Type("TEXT").Constraint().NotNull().Unique()
	table.AddColumn("sort_order").Type("INT")
	table.AddColumn("user_id").Type("BIGINT").Constraint().Default("1")
	table.AddColumn("description").Type("TEXT")
	table.AddColumn("params").Type("JSONB")
	table.AddColumn("external_ids").Type("INT[]")
	table.AddColumn("file_ids").Type("BIGINT[]").Constraint().NotNull()
	table.AddColumn("status_id").Type("INT").Constraint().NotNull().
		References().RefTable("dictionary").Column("id").OnDelete(gosql.ActionRestrict).OnUpdate(gosql.ActionCascade)
	table.AddColumn("is_failed").Type("bool")
	table.AddColumn("uuids").Type("UUID[]")
	table.AddColumn("number").Type("NUMERIC(6, 2)")
	table.AddColumn("prices").Type("NUMERIC(6, 2)[]")

	list := gosql.SQList{table,
		gosql.NewComment().Table(schema+"test_table", "Test table"),
		gosql.NewComment().Column(schema+"test_table.id", "Test table identifier"),
		gosql.NewComment().Column(schema+"test_table.created_at", "Test table created time"),
		gosql.NewComment().Column(schema+"test_table.updated_at", "Test table updated time"),
		gosql.NewComment().Column(schema+"test_table.name", "Test table name"),
		gosql.NewComment().Column(schema+"test_table.sort_order", "Test table sort order"),
		gosql.NewComment().Column(schema+"test_table.user_id", "Test table user_id"),
		gosql.NewComment().Column(schema+"test_table.description", "Test table description"),
		gosql.NewComment().Column(schema+"test_table.params", "Test table parameters"),
		gosql.NewComment().Column(schema+"test_table.external_ids", "Test table external ids"),
		gosql.NewComment().Column(schema+"test_table.file_ids", "Test table file ids"),
		gosql.NewComment().Column(schema+"test_table.is_failed", "Test table is failed flag"),
		gosql.NewComment().Column(schema+"test_table.uuids", "Test table list of uuids"),
		gosql.NewComment().Column(schema+"test_table.number", "Test table float number"),
		gosql.NewComment().Column(schema+"test_table.prices", "Test table list of prices"),
	}
	return list
}

func TestDBO_InitError(t *testing.T) {
	_, err := godb.DBO{
		Options: godb.Options{
			Debug:          true,
			Logger:         gocli.NewLogger(gocli.LoggerConfig{}),
			QueryProcessor: godb.PreparePositionalArgsQuery,
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
			Debug:          true,
			Logger:         gocli.NewLogger(gocli.LoggerConfig{}),
			QueryProcessor: godb.PreparePositionalArgsQuery,
		},
		Connection: &connection{},
	}.Init()
}

func TestCreateTestTable(t *testing.T) {
	db, err := initDb()
	if err != nil {
		t.Fatal(err)
	}
	err = CreateDictionaryTable(db)
	if err != nil {
		t.Fatal(err)
	}
	var schema = "kpi"
	var schemaQuery = "CREATE SCHEMA IF NOT EXISTS %s;"
	schemaQuery = fmt.Sprintf(schemaQuery, schema)
	_, err = db.Exec(schemaQuery)
	if err != nil {
		t.Fatal(err)
	}
	schema += "."
	list := testTable(schema)
	query, _, _ := list.Join()
	_, err = db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
	col := NewCollection[DictionaryModel]()
	col.AddItem(&DictionaryModel{
		Id:    gohelp.Ptr[int32](1000),
		Type:  gohelp.Ptr("test_table_status"),
		Code:  gohelp.Ptr("new"),
		Label: gohelp.Ptr("New status"),
	})
	e := col.Save(db)
	if e != nil {
		t.Fatal(e)
	}
}

//func NATestSaveTestTable(t *testing.T) {
//	db, err := initDb()
//	if err != nil {
//		t.Fatal(err)
//	}
//	model := NewTestTable()
//	model.Name = gohelp.Ptr(gohelp.RandString(10))
//	model.SortOrder = gohelp.Ptr[int32](10)
//	model.UserId = gohelp.Ptr[int64](1000)
//	model.FileIds = pq.Int64Array{1, 2, 100, 3000}
//	model.IsFailed = gohelp.Ptr(true)
//	model.Uuids = pq.StringArray{gohelp.NewUUID(), gohelp.NewUUID(), gohelp.NewUUID()}
//	model.Number = gohelp.Ptr[float32](1.2333)
//	model.Prices = pq.Float32Array{1.233, 2.3444, 4.555}
//	e := Save(db, model)
//	if e != nil {
//		t.Fatal(e)
//	}
//	model.Uuids = nil
//	model.Params = gohelp.Ptr(json.RawMessage(`{"some": 123}`))
//	model.ExternalIds = pq.Int32Array{1, 2, 3, 4}
//	model.Prices = nil
//	e = Save(db, model)
//	if e != nil {
//		t.Fatal(e)
//	}
//
//	tt := NewTestTable()
//	tt.Id = model.Id
//	e = Load(db, tt)
//	if e != nil {
//		t.Fatal(e)
//	}
//}

func TestMakeTestTableModel(t *testing.T) {
	db, err := initDb()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = MakeModel(db, "models", "public", "test_table", "", DefaultSystemColumnsSoft)
	if err != nil {
		t.Fatal(err)
	}
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
	err = crud.Generate(db, "public", "test_table", "v2", 31)
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
