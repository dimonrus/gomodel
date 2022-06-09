package gomodel

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed crud.tmpl
var CrudTemplate string

//go:embed api_read.tmpl
var ApiReadTemplate string

//go:embed api_update.tmpl
var ApiUpdateTemplate string

//go:embed api_create.tmpl
var ApiCreateTemplate string

//go:embed api_delete.tmpl
var ApiDeleteTemplate string

// CRUDGenerator struct for crud generation
type CRUDGenerator struct {
	// Path for crud folder
	CRUDPath string
	// Path for client folder
	ClientPath string
	// Path for api folder
	APIPath string
	// Path for project
	ProjectPath string
}

// NewCRUDGenerator init crud generator
func NewCRUDGenerator(CRUDPath, ClientPath, APIPath, ProjectPath string) *CRUDGenerator {
	return &CRUDGenerator{
		CRUDPath:    CRUDPath,
		ClientPath:  ClientPath,
		APIPath:     APIPath,
		ProjectPath: ProjectPath,
	}
}

// MakeCoreCrud generate core file
func (c CRUDGenerator) MakeCoreCrud(q godb.Queryer, schema, table string) error {
	// New Template
	tmp := template.New("crud").Funcs(getHelperFunc(DefaultSystemColumns))
	// init template
	tmlString := CrudTemplate

	// Get package name
	packageNames := strings.Split(c.CRUDPath, string(os.PathSeparator))
	var packageName string
	if len(packageNames) > 0 {
		packageName = packageNames[len(packageNames)-1]
	} else {
		packageName = c.CRUDPath
	}

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp, err = tmp.Parse(tmlString)
	if err != nil {
		return err
	}

	file, path, err := CreateModelFile(schema, table, c.CRUDPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// crud imports
	var imports = []string{
		`"github.com/dimonrus/gomodel"`,
		`"github.com/dimonrus/porterr"`,
		`"github.com/dimonrus/godb/v2"`,
		fmt.Sprintf(`"%s/%s"`, c.ProjectPath, c.ClientPath),
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Model   string
		Imports []string
		Columns Columns
	}{
		Package: packageName,
		Model:   getModelName(schema, table),
		Imports: imports,
		Columns: *columns,
	})

	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPIRead generate qpi read
func (c CRUDGenerator) MakeAPIRead(q godb.Queryer, schema, table, version string) error {
	// New Template
	tmp := template.New("api_read").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiReadTemplate

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("read", c.APIPath+string(os.PathSeparator)+table+string(os.PathSeparator)+version)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"strconv"`,
		`"github.com/gorilla/mux"`,
		`"github.com/dimonrus/gorest"`,
		fmt.Sprintf(`"%s/app/base"`, c.ProjectPath),
		fmt.Sprintf(`"%s/%s"`, c.ProjectPath, c.CRUDPath),
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Model   string
		Imports []string
		Columns Columns
	}{
		Package: version,
		Model:   getModelName(schema, table),
		Imports: imports,
		Columns: *columns,
	})

	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPIDelete generate qpi delete
func (c CRUDGenerator) MakeAPIDelete(q godb.Queryer, schema, table, version string) error {
	// New Template
	tmp := template.New("api_delete").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiDeleteTemplate

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("delete", c.APIPath+string(os.PathSeparator)+table+string(os.PathSeparator)+version)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"strconv"`,
		`"github.com/gorilla/mux"`,
		`"github.com/dimonrus/gorest"`,
		fmt.Sprintf(`"%s/app/base"`, c.ProjectPath),
		fmt.Sprintf(`"%s/%s"`, c.ProjectPath, c.CRUDPath),
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Model   string
		Imports []string
		Columns Columns
	}{
		Package: version,
		Model:   getModelName(schema, table),
		Imports: imports,
		Columns: *columns,
	})

	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPIRead generate qpi update
func (c CRUDGenerator) MakeAPIUpdate(q godb.Queryer, schema, table, version string) error {
	// New Template
	tmp := template.New("api_update").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiUpdateTemplate

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("update", c.APIPath+string(os.PathSeparator)+table+string(os.PathSeparator)+version)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"strconv"`,
		`"github.com/gorilla/mux"`,
		`"github.com/dimonrus/gorest"`,
		fmt.Sprintf(`"%s/app/base"`, c.ProjectPath),
		fmt.Sprintf(`"%s/%s"`, c.ProjectPath, c.CRUDPath),
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Model   string
		Imports []string
		Columns Columns
	}{
		Package: version,
		Model:   getModelName(schema, table),
		Imports: imports,
		Columns: *columns,
	})

	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPICreate generate qpi create
func (c CRUDGenerator) MakeAPICreate(q godb.Queryer, schema, table, version string) error {
	// New Template
	tmp := template.New("api_create").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiCreateTemplate

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("create", c.APIPath+string(os.PathSeparator)+table+string(os.PathSeparator)+version)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"github.com/dimonrus/gorest"`,
		fmt.Sprintf(`"%s/app/base"`, c.ProjectPath),
		fmt.Sprintf(`"%s/%s"`, c.ProjectPath, c.CRUDPath),
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Model   string
		Imports []string
		Columns Columns
	}{
		Package: version,
		Model:   getModelName(schema, table),
		Imports: imports,
		Columns: *columns,
	})

	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// Generate generate crud, client, api
func (c CRUDGenerator) Generate(q godb.Queryer, schema, table, version string) error {
	modelTemplate := DefaultModelTemplate
	err := MakeModel(q, c.ClientPath, schema, table, modelTemplate, DefaultSystemColumnsSoft)
	if err != nil {
		return err
	}
	err = c.MakeCoreCrud(q, schema, table)
	if err != nil {
		return err
	}
	err = c.MakeAPIRead(q, schema, table, version)
	if err != nil {
		return err
	}
	err = c.MakeAPIUpdate(q, schema, table, version)
	if err != nil {
		return err
	}
	err = c.MakeAPIDelete(q, schema, table, version)
	if err != nil {
		return err
	}
	err = c.MakeAPICreate(q, schema, table, version)
	if err != nil {
		return err
	}
	return err
}
