package gomodel

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed model.tmpl
var DefaultModelTemplate string

var DefaultSystemColumnsSoft = SystemColumns{Created: "created_at", Updated: "updated_at", Deleted: "deleted_at"}
var DefaultSystemColumns = SystemColumns{Created: "created_at", Updated: "updated_at"}

type SystemColumns struct {
	Created string
	Updated string
	Deleted string
}

// Column information
type Column struct {
	Name              string  // DB column name
	ModelName         string  // Model name
	Default           *string // DB default value
	IsNullable        bool    // DB is nullable
	IsByteArray       bool    // Do not need type pointer for []byte
	DataType          string  // DB column type
	ModelType         string  // Model type
	Schema            string  // DB Schema
	Table             string  // DB table
	Sequence          *string // DB sequence
	ForeignSchema     *string // DB foreign schema name
	ForeignTable      *string // DB foreign table name
	ForeignColumnName *string // DB foreign column name
	ForeignIsSoft     bool    // DB foreign table is soft
	Description       *string // DB column description
	IsPrimaryKey      bool    // DB is primary key
	Tags              string  // Model Tags name
	Import            string  // Model Import custom lib
	IsArray           bool    // Array column
	IsCreated         bool    // Is created at column
	IsUpdated         bool    // Is updated at column
	IsDeleted         bool    // Is deleted at column
	HasUniqueIndex    bool    // If column is a part of unique index
	UniqueIndexName   *string // Unique index name
	TableDescription  *string // Table description
	DefaultTypeValue  *string // Default value for type
}

// GetModelFieldTag Prepare ModelFiledTag by Column
func (c Column) GetModelFieldTag() (field ModelFiledTag) {
	field.Column = c.Name
	if c.ForeignColumnName != nil {
		field.ForeignKey = *c.ForeignSchema + "." + *c.ForeignTable + "." + *c.ForeignColumnName
	}
	field.IsSequence = c.Sequence != nil
	field.IsRequired = !c.IsNullable
	field.IsUnique = c.HasUniqueIndex
	field.IsPrimaryKey = c.IsPrimaryKey
	field.IsCreatedAt = c.IsCreated
	field.IsUpdatedAt = c.IsUpdated
	field.IsDeletedAt = c.IsDeleted
	field.IsArray = c.IsArray
	return
}

// PrepareValidTag if dictionary item
func (c Column) PrepareValidTag(dictionary DictionaryItems) string {
	var valid []string
	var validTag string
	if !c.IsNullable && c.Sequence == nil && !(c.IsCreated || c.IsUpdated || c.IsDeleted) {
		valid = []string{"required"}
	}
	if c.ForeignTable != nil && *c.ForeignTable == "dictionary" {
		if dType, ok := dictionary.IsDictionaryColumn(c.Name); ok {
			enum := dictionary.GetTypeEnum(dType)
			if len(enum) > 0 {
				valid = append(valid, "enum~"+enum)
			}
		}
	}
	if len(valid) > 0 {
		validTag = ` valid:"` + strings.Join(valid, ";") + `"`
	}
	return validTag
}

// Array of columns
type Columns []Column

// Get imports
func (c Columns) GetImports() []string {
	// imports in model file
	var imports []string

	for i := range c {
		if c[i].Import == "" {
			continue
		}
		imports = gohelp.AppendUnique(imports, c[i].Import)
	}

	return imports
}

// GetTableDescription get table comment
func (c Columns) GetTableDescription() string {
	if len(c) > 0 && c[0].TableDescription != nil {
		return strings.Join(strings.Split(*c[0].TableDescription, "\n"), "\n// ")
	}
	return ""
}

// Parse Row
func parseColumnRow(rows *sql.Rows) (*Column, error) {
	column := Column{}
	err := rows.Scan(
		&column.Name,
		&column.DataType,
		&column.IsNullable,
		&column.Schema,
		&column.Table,
		&column.IsPrimaryKey,
		&column.Default,
		&column.Sequence,
		&column.ForeignSchema,
		&column.ForeignTable,
		&column.ForeignColumnName,
		&column.ForeignIsSoft,
		&column.Description,
		&column.HasUniqueIndex,
		&column.UniqueIndexName,
		&column.TableDescription,
	)

	if err != nil {
		return nil, err
	}

	return &column, nil
}

// Get table columns from db
func GetTableColumns(dbo godb.Queryer, schema string, table string, sysCols SystemColumns, dictionary DictionaryItems) (*Columns, error) {
	query := fmt.Sprintf(`
SELECT a.attname                                                                       AS column_name,
       format_type(a.atttypid, a.atttypmod)                                            AS data_type,
       CASE WHEN a.attnotnull THEN FALSE ELSE TRUE END                                 AS is_nullable,
       s.nspname                                                                       AS schema,
       t.relname                                                                       AS table,
       (SELECT EXISTS(SELECT i.indisprimary
                      FROM pg_index i
                      WHERE i.indrelid = a.attrelid
                        AND a.attnum = ANY (i.indkey)
                        AND i.indisprimary IS TRUE))                                   AS is_primary,
       ic.column_default,
       pg_get_serial_sequence(ic.table_schema || '.' || ic.table_name, ic.column_name) AS sequence,
       max(ccu.table_schema)                                                           AS foreign_schema,
       max(ccu.table_name)                                                             AS foreign_table,
       max(ccu.column_name)                                                            AS foreign_column_name,
       (select EXISTS(SELECT 1
                      from information_schema.columns
                      where column_name = 'deleted_at'
                        and table_name = max(ccu.table_name)))                         AS is_foreign_soft,
       col_description(t.oid, a.attnum)                                                AS description,
       (SELECT EXISTS(SELECT i.indisunique
                      FROM pg_index i
                      WHERE i.indrelid = a.attrelid
                        AND i.indisunique IS TRUE
                        AND a.attnum = ANY (i.indkey)))                                AS has_unique_index,
       (SELECT ins.indexname
        FROM pg_indexes ins
                 JOIN pg_index i ON ins.indexdef = pg_get_indexdef(i.indexrelid)
        WHERE i.indisunique IS TRUE
          AND i.indrelid = a.attrelid
          AND a.attnum = ANY (i.indkey))                                               AS unique_index_name,
        obj_description(t.oid)                                                         AS table_description
FROM pg_attribute a
         JOIN pg_class t ON a.attrelid = t.oid
         JOIN pg_namespace s ON t.relnamespace = s.oid
         LEFT JOIN information_schema.columns AS ic
                   ON ic.column_name = a.attname AND ic.table_name = t.relname AND ic.table_schema = s.nspname
         LEFT JOIN information_schema.key_column_usage AS kcu
                   ON kcu.table_name = t.relname AND kcu.column_name = a.attname AND kcu.table_schema = s.nspname
         LEFT JOIN information_schema.table_constraints AS tc
                   ON tc.constraint_name = kcu.constraint_name AND tc.constraint_type = 'FOREIGN KEY' AND tc.table_schema = kcu.constraint_schema
         LEFT JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name AND tc.table_schema = ccu.table_schema
WHERE a.attnum > 0
  AND NOT a.attisdropped
  AND s.nspname = '%s'
  AND t.relname = '%s'
GROUP BY a.attname, a.atttypid, a.attrelid, a.atttypmod, a.attnotnull, s.nspname, t.relname, ic.column_default,
         ic.table_schema, ic.table_name, ic.column_name, a.attnum, t.oid, ic.ordinal_position
ORDER BY a.attnum;`, schema, table)

	rows, err := dbo.Query(query)
	if err != nil {
		return nil, err
	}

	var columns Columns
	var hasPrimary bool

	for rows.Next() {
		column, err := parseColumnRow(rows)
		if err != nil {
			return nil, err
		}
		name := gohelp.ToCamelCase(column.Name, true)
		json := gohelp.ToCamelCase(column.Name, false)
		column.ModelName = name
		if column.Sequence == nil && column.Default != nil {
			if strings.Contains(*column.Default, "seq") {
				column.Sequence = new(string)
				*column.Sequence = *column.Default
			}
		}
		if column.Name == sysCols.Created {
			column.IsCreated = true
		}
		if column.Name == sysCols.Updated {
			column.IsUpdated = true
		}
		if column.Name == sysCols.Deleted {
			column.IsDeleted = true
		}
		if column.Description != nil {
			*column.Description = strings.Join(strings.Split(*column.Description, "\n"), "\n// ")
		}
		fTag := column.GetModelFieldTag()

		validTag := column.PrepareValidTag(dictionary)

		column.Tags = fmt.Sprintf(`%cdb:"%s" json:"%s"%s%c`, '`', fTag.String(), json, validTag, '`')

		switch {
		case column.DataType == "bigint":
			column.ModelType = "int64"
		case column.DataType == "integer":
			column.ModelType = "int32"
		case column.DataType == "text":
			column.ModelType = "string"
		case column.DataType == "double precision":
			column.ModelType = "float64"
		case column.DataType == "boolean":
			column.ModelType = "bool"
		case column.DataType == "ARRAY":
			column.ModelType = "[]interface{}"
		case column.DataType == "json":
			column.ModelType = "json.RawMessage"
			column.Import = `"encoding/json"`
			column.IsByteArray = true
		case column.DataType == "smallint":
			column.ModelType = "int16"
		case column.DataType == "date":
			column.ModelType = "time.Time"
			column.Import = `"time"`
		case strings.Contains(column.DataType, "character varying"):
			column.ModelType = "string"
		case strings.Contains(column.DataType, "numeric"):
			column.ModelType = "float32"
		case column.DataType == "uuid":
			column.ModelType = "string"
		case column.DataType == "jsonb":
			column.ModelType = "json.RawMessage"
			column.Import = `"encoding/json"`
			column.IsByteArray = true
		case column.DataType == "uuid[]":
			column.ModelType = "[]string"
			column.IsArray = true
			column.Import = `"github.com/lib/pq"`
		case column.DataType == "integer[]":
			column.ModelType = "[]int64"
			column.IsArray = true
			column.Import = `"github.com/lib/pq"`
		case column.DataType == "bigint[]":
			column.ModelType = "[]int64"
			column.IsArray = true
			column.Import = `"github.com/lib/pq"`
		case column.DataType == "text[]":
			column.ModelType = "[]string"
			column.IsArray = true
			column.Import = `"github.com/lib/pq"`
		case strings.Contains(column.DataType, "timestamp"):
			column.ModelType = "time.Time"
			column.Import = `"time"`
		default:
			return nil, errors.New(fmt.Sprintf("unknown column type: %s", column.DataType))
		}

		if column.IsNullable && !column.IsArray {
			column.ModelType = "*" + column.ModelType
		}

		if !column.IsNullable {
			if strings.Contains(column.ModelType, "int") || strings.Contains(column.ModelType, "float") {
				column.DefaultTypeValue = new(string)
				*column.DefaultTypeValue = "0"
			} else {
				column.DefaultTypeValue = new(string)
				*column.DefaultTypeValue = `""`
			}
		} else {
			column.DefaultTypeValue = new(string)
			*column.DefaultTypeValue = "nil"
		}

		if column.IsPrimaryKey == true {
			hasPrimary = true
		}

		columns = append(columns, *column)
	}

	// column named id will be primary if no primary key
	if !hasPrimary {
		for key, column := range columns {
			if column.Name == "id" {
				columns[key].IsPrimaryKey = true
				if columns[key].ModelType[0] == '*' {
					columns[key].ModelType = columns[key].ModelType[1:]
					hasPrimary = true
				}
				break
			}
		}
		// if still no primary key
		if !hasPrimary {
			// Collect primary kye by unique index
			var uniqueIndexName *string
			for key, column := range columns {
				if column.HasUniqueIndex {
					if uniqueIndexName == nil {
						uniqueIndexName = column.UniqueIndexName
					}
					if *uniqueIndexName == *column.UniqueIndexName {
						columns[key].IsPrimaryKey = true
						if columns[key].ModelType[0] == '*' {
							columns[key].ModelType = columns[key].ModelType[1:]
						}
					}
				}
			}
		}
	}

	return &columns, nil
}

// Template helper functions
func getHelperFunc(systemColumns SystemColumns) template.FuncMap {
	return template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"system": func(column Column) bool {
			return gohelp.ExistsInArray(column.Name, []string{systemColumns.Created, systemColumns.Updated, systemColumns.Deleted}) ||
				(column.IsPrimaryKey && column.Sequence != nil)
		},
		"isSearchColumn": func(column Column) bool {
			return !(column.IsDeleted || column.IsCreated || column.IsUpdated || column.Import == `"time"`)
		},
		"cameled": func(name string) string {
			return gohelp.ToCamelCase(name, true)
		},
		"unref": func(modelType string) string {
			if modelType[0] == '*' {
				return modelType[1:]
			}
			return modelType
		},
		"icameled": func(name string) string {
			return gohelp.ToCamelCase(name, false)
		},
		"foreign": func(name string) string {
			if name[len(name)-3:] == "_id" {
				name = name[:len(name)-3]
			}
			return gohelp.ToCamelCase(name, true)
		},
		"model": func(schema string, table string) string {
			return getModelName(schema, table)
		},
		"pointerType": func(modelType string) string {
			if modelType[0] != '*' {
				return "*" + modelType
			}
			return modelType
		},
		"parseIdentifier": func(name string, modelType string) string {
			var result string
			switch modelType {
			case "string", "[]byte":
				result = fmt.Sprintf(`%s := mux.Vars(r)["%s"]`, name, gohelp.ToUnderscore(name))
			case "int", "int64", "int32":
				result = fmt.Sprintf(`%s, _ := strconv.ParseInt(mux.Vars(r)["%s"], 10, 64)`, name, gohelp.ToUnderscore(name))
			case "uint", "uint64", "uint32":
				result = fmt.Sprintf(`%s, _ := strconv.ParseUint(mux.Vars(r)["%s"], 10, 64)`, name, gohelp.ToUnderscore(name))
			}
			return result
		},
		"prepareIdentifierParam": func(name string, modelType string) string {
			var result string
			switch modelType {
			case "string", "[]byte":
				result = fmt.Sprintf(`{%s}`, gohelp.ToUnderscore(name))
			case "int", "int64", "int32", "uint", "uint64", "uint32":
				result = fmt.Sprintf(`{%s:^[0-9]+$}`, gohelp.ToUnderscore(name))
			}
			return result
		},
		"castToIdentifier": func(name string, modelType string) string {
			var result string
			switch modelType {
			case "string", "[]byte", "uint64", "int64":
				result = name
			case "int":
				result = fmt.Sprintf("int(%s)", name)
			case "int32":
				result = fmt.Sprintf("int32(%s)", name)
			case "uint":
				result = fmt.Sprintf("uint(%s)", name)
			case "uint32":
				result = fmt.Sprintf("uint32(%s)", name)
			}
			return result
		},
		"searchField": multipleValueName,
		"jsonTag": func(name string) string {
			return fmt.Sprintf(`%cjson:"%s"%c`, '`', gohelp.ToCamelCase(multipleValueName(name), false), '`')
		},
	}
}

func multipleValueName(name string) string {
	if name[len(name)-1] != 's' {
		return name + "s"
	}
	return name
}

// CreateFile Create File
func CreateFile(name, path string) (*os.File, string, error) {
	var filePath string
	if path != "" {
		folderPath := fmt.Sprintf(path)
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return nil, "", err
		}
		filePath = fmt.Sprintf("%s/%s.go", folderPath, name)
	} else {
		filePath = fmt.Sprintf("%s.go", name)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return nil, "", err
	}
	return f, filePath, nil
}

// Create file in os
func CreateModelFile(schema string, table string, path string) (*os.File, string, error) {
	fileName := fmt.Sprintf("%s", table)
	if schema != "public" {
		fileName = fmt.Sprintf("%s_%s", schema, table)
	}
	return CreateFile(fileName, path)
}

// Prepare model name
func getModelName(schema string, table string) string {
	var name string
	if schema == "public" || schema == "" {
		name = gohelp.ToCamelCase(table, true)
	} else {
		name = gohelp.ToCamelCase(schema+"_"+table, true)
	}
	return name
}

// Create model
func MakeModel(db godb.Queryer, dir string, schema string, table string, templatePath string, systemColumns SystemColumns) error {
	// Imports in model file
	var imports = []string{
		`"github.com/dimonrus/gomodel"`,
	}

	if table == "" {
		return errors.New("table name is empty")
	}

	// New Template
	tmpl := template.New("model").Funcs(getHelperFunc(systemColumns))

	var tmlString = DefaultModelTemplate

	templateFile, err := os.Open(templatePath)
	if err == nil {
		// Read template
		data, err := ioutil.ReadAll(templateFile)
		if err != nil {
			return err
		}
		tmlString = string(data)
	} else if tmlString == "" {
		return err
	}

	dictionaryItems := getDictionaryItems(db)

	// Open model template
	tmpl = template.Must(tmpl.Parse(tmlString))

	// Columns
	columns, err := GetTableColumns(db, schema, table, systemColumns, dictionaryItems)
	if err != nil {
		return err
	}

	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	// Table comment
	tableDescription := columns.GetTableDescription()

	// Collect imports
	for _, column := range *columns {
		imports = gohelp.AppendUnique(imports, column.Import)
	}

	// To camel case
	modelName := getModelName(schema, table)

	var hasSequence bool
	// Check for sequence and primary key
	for _, column := range *columns {
		if column.IsPrimaryKey && column.Sequence != nil {
			hasSequence = true
			break
		}
	}
	packageNames := strings.Split(dir, string(os.PathSeparator))
	var packageName string
	if len(packageNames) > 0 {
		packageName = packageNames[len(packageNames)-1]
	} else {
		packageName = dir
	}

	// Create file
	file, path, err := CreateModelFile(schema, table, dir)
	if err != nil {
		return err
	}

	// Parse template to file
	err = tmpl.Execute(file, struct {
		Package          string
		Model            string
		Table            string
		TableDescription string
		Columns          Columns
		HasSequence      bool
		Imports          []string
	}{
		Package:          packageName,
		Model:            modelName,
		Table:            table,
		TableDescription: tableDescription,
		Columns:          *columns,
		HasSequence:      hasSequence,
		Imports:          imports,
	})

	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	err = cmd.Run()
	if err != nil {
		return err
	}

	if dbo, ok := db.(*godb.DBO); ok {
		dbo.Logger.Printf("Model file created: %s", path)
	}

	// Create all foreign models if not exists
	for i := range *columns {
		c := (*columns)[i]
		if c.ForeignTable != nil {
			var found bool
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if info == nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				data, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}
				modelName := getModelName(schema, *c.ForeignTable)
				if strings.Contains(string(data), fmt.Sprintf("type %s struct {", modelName)) {
					found = true
				}
				return nil
			})
			if err != nil {
				return err
			}
			if !found {
				err = MakeModel(db, dir, *c.ForeignSchema, *c.ForeignTable, templatePath, systemColumns)
				if err != nil {
					db.(*godb.DBO).Logger.Printf("Model file generator error: %s", err.Error())
				}
			}
		}
	}
	return nil
}

// TODO generate crud in core
// TODO generate list and search
// TODO generate client with methods
// TODO generate apis
