package gomodel

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gosql"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//go:embed dictionary_mapping.tmpl
var DefaultDictionaryTemplate string

// DictionaryItems dictionary collection items
type DictionaryItems []*DictionaryModel

// DictionaryModel model
type DictionaryModel struct {
	// Dictionary row identifier
	Id *int32 `db:"col~id;prk;req;unq;" json:"id" valid:"required"`
	// Dictionary row type
	Type *string `db:"col~type;req;" json:"type" valid:"required"`
	// Dictionary row code
	Code *string `db:"col~code;req;" json:"code" valid:"required"`
	// Dictionary row value label
	Label *string `db:"col~label;" json:"label"`
	// Dictionary row created time
	CreatedAt *time.Time `db:"col~created_at;req;cat;" json:"createdAt"`
	// Dictionary row updated time
	UpdatedAt *time.Time `db:"col~updated_at;uat;" json:"updatedAt"`
	// Dictionary row deleted time
	DeletedAt *time.Time `db:"col~deleted_at;dat;" json:"deletedAt"`
}

// Table Model columns
func (m *DictionaryModel) Table() string {
	return "dictionary"
}

// Columns Model columns
func (m *DictionaryModel) Columns() []string {
	return []string{"id", "type", "code", "label", "created_at", "updated_at", "deleted_at"}
}

// Values Model values
func (m *DictionaryModel) Values() (values []interface{}) {
	return []interface{}{&m.Id, &m.Type, &m.Code, &m.Label, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

// HasType check if type in collection
func (i DictionaryItems) HasType(dictionaryType string) bool {
	for _, model := range i {
		if *model.Type == dictionaryType {
			return true
		}
	}
	return false
}

// IsDictionaryColumn check if column name has dictionary reference
func (i DictionaryItems) IsDictionaryColumn(name string) (string, bool) {
	if strings.Index(name, "_id") != len(name)-3 {
		return "", false
	}
	clear := strings.Replace(name, "_id", "", -1)
	for _, model := range i {
		if *model.Type == clear {
			return clear, true
		}
	}
	return "", false
}

// GetTypeEnum get enum type for validation
func (i DictionaryItems) GetTypeEnum(dictionaryType string) string {
	var result string
	for _, dictionary := range i {
		if *dictionary.Type == dictionaryType {
			if result == "" {
				result += strconv.Itoa(int(*dictionary.Id))
			} else {
				result += "," + strconv.Itoa(int(*dictionary.Id))
			}
		}
	}
	return result
}

// Get all sql
func getDictionarySQList() gosql.SQList {
	var sqList gosql.SQList

	dict := gosql.CreateTable("dictionary").IfNotExists()
	dict.AddColumn("id").Type("INT").Constraint().NotNull().PrimaryKey()
	dict.AddColumn("type").Type("TEXT").Constraint().NotNull()
	dict.AddColumn("code").Type("TEXT").Constraint().NotNull()
	dict.AddColumn("label").Type("TEXT")

	gosql.TableModeler{TimestampModifier, SoftModifier}.Prepare(dict)

	sqList = append(sqList, dict,
		gosql.Comment().Column("dictionary.id", "Dictionary row identifier"),
		gosql.Comment().Column("dictionary.type", "Dictionary row type"),
		gosql.Comment().Column("dictionary.code", "Dictionary row code"),
		gosql.Comment().Column("dictionary.label", "Dictionary row value label"),
		gosql.Comment().Column("dictionary.created_at", "Dictionary row created time"),
		gosql.Comment().Column("dictionary.updated_at", "Dictionary row updated time"),
		gosql.Comment().Column("dictionary.deleted_at", "Dictionary row deleted time"),
		gosql.CreateIndex("dictionary", "type").IfNotExists().AutoName())

	return sqList
}

// Create Table
func CreateDictionaryTable(q godb.Queryer) error {
	query, _, _ := getDictionarySQList().Join()
	_, err := q.Exec(query)
	return err
}

// Create or update dictionary mapping
func GenerateDictionaryMapping(path string, q godb.Queryer) error {
	dictionaries := getDictionaryItems(q)
	if len(dictionaries) == 0 {
		return errors.New("no dictionary in database")
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	paths := strings.Split(path, fmt.Sprintf("%c", os.PathSeparator))
	packageName := paths[len(paths)-2]
	tml := getDictionaryTemplate()
	err = tml.Execute(f, struct {
		Dictionaries []*DictionaryModel
		Package      string
	}{
		Dictionaries: dictionaries,
		Package:      packageName,
	})

	if err != nil {
		_ = os.RemoveAll(path)
		return err
	}

	cmd := exec.Command("go", "fmt", path)

	return cmd.Run()
}

func getDictionaryTemplate() *template.Template {
	funcMap := template.FuncMap{
		"camelCase": func(str string) string {
			return gohelp.ToCamelCase(str, true)
		},
		"deref": func(str *string) string {
			if str != nil {
				return *str
			}
			return ""
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(DefaultDictionaryTemplate))
}

// Get all dictionary items sorted by type and created_at
func getDictionaryItems(q godb.Queryer) DictionaryItems {
	collection := NewCollection[DictionaryModel]()
	collection.AddOrder("type", "created_at", "id")
	e := collection.Load(q)
	if e != nil {
		return nil
	}
	return collection.Items()
}
