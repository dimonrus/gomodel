package generator

import (
	_ "embed"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gomodel"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

//go:embed dictionary_mapping.tmpl
var DefaultDictionaryTemplate string

type DictionaryModel struct {
	// Dictionary row identifier
	Id *int32 `json:"id"`
	// Dictionary row type
	Type *string `json:"type"`
	// Dictionary row code
	Code *string `json:"code"`
	// Dictionary row value label
	Label *string `json:"label"`
	// Dictionary row created time
	CreatedAt *time.Time `json:"createdAt"`
	// Dictionary row updated time
	UpdatedAt *time.Time `json:"updatedAt"`
	// Dictionary row deleted time
	DeletedAt *time.Time `json:"deletedAt"`
}

// Model columns
func (m *DictionaryModel) Table() string {
	return "dictionary"
}

// Model columns
func (m *DictionaryModel) Columns() []string {
	return []string{"id", "type", "code", "label", "created_at", "updated_at", "deleted_at"}
}

// Model values
func (m *DictionaryModel) Values() (values []interface{}) {
	return []interface{}{&m.Id, &m.Type, &m.Code, &m.Label, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

// Create Table
func CreateDictionaryTable(q godb.Queryer) error {
	query := `
CREATE TABLE IF NOT EXISTS dictionary
(
  id         INT PRIMARY KEY                                 NOT NULL,
  type       TEXT                                            NOT NULL,
  code       TEXT                                            NOT NULL,
  label      TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT localtimestamp NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

COMMENT ON COLUMN dictionary.id IS 'Dictionary row identifier';
COMMENT ON COLUMN dictionary.type IS 'Dictionary row type';
COMMENT ON COLUMN dictionary.code IS 'Dictionary row code';
COMMENT ON COLUMN dictionary.label IS 'ОDictionary row value label';
COMMENT ON COLUMN dictionary.created_at IS 'Dictionary row created time';
COMMENT ON COLUMN dictionary.updated_at IS 'Dictionary row updated time';
COMMENT ON COLUMN dictionary.deleted_at IS 'Dictionary row deleted time';

CREATE INDEX IF NOT EXISTS dictionary_type_idx ON dictionary (type);`

	_, err := q.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Create or update dictionary mapping
func GenerateDictionaryMapping(path string, q godb.Queryer) error {
	collection := gomodel.NewCollection[DictionaryModel]()
	collection.AddOrder("type", "created_at", "id")
	e := collection.Load(q)
	if e != nil {
		return e
	}
	dictionaries := collection.Items()

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
		err = os.RemoveAll(path)
	}
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "fmt", path)
	err = cmd.Run()

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
