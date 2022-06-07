package gomodel

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed crud.tmpl
var CrudTemplate string

//go:embed api_read.tmpl
var ApiReadTemplate string

// MakeCoreCrud Make core cud file
func MakeCoreCrud(q godb.Queryer, crudPath, clientPath, project, schema, table, tmpl string) error {
	// New Template
	tmp := template.New("crud").Funcs(getHelperFunc(DefaultSystemColumns))
	var tmlString = CrudTemplate

	templateFile, err := os.Open(tmpl)
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

	packageNames := strings.Split(crudPath, string(os.PathSeparator))
	var packageName string
	if len(packageNames) > 0 {
		packageName = packageNames[len(packageNames)-1]
	} else {
		packageName = crudPath
	}

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}

	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateModelFile(schema, table, crudPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"github.com/dimonrus/gomodel"`,
		`"github.com/dimonrus/porterr"`,
		`"github.com/dimonrus/godb/v2"`,
		fmt.Sprintf(`"%s/%s"`, project, clientPath),
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

// GenerateCrud generate model in client pack and generate crud in core
func GenerateCrud(crudPath, clientPath, project string, schema, table, tmpl string, q godb.Queryer) error {
	err := MakeModel(q, clientPath, schema, table, tmpl, DefaultSystemColumnsSoft)
	if err != nil {
		return err
	}
	err = MakeCoreCrud(q, crudPath, clientPath, project, schema, table, tmpl)
	if err != nil {
		return err
	}
	//err = MakeAPIs(q, project, "io/web/api/user/v1", "", schema, table)
	return err
}

func MakeAPIs(q godb.Queryer, project, versionPath, tmpl, schema, table string) error {
	// New Template
	tmp := template.New("crud").Funcs(getHelperFunc(DefaultSystemColumns))
	var tmlString = ApiReadTemplate

	templateFile, err := os.Open(tmpl)
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

	packageNames := strings.Split(versionPath, string(os.PathSeparator))
	var packageName string
	if len(packageNames) > 0 {
		packageName = packageNames[len(packageNames)-1]
	} else {
		packageName = versionPath
	}

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}

	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateModelFile(schema, table, versionPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"strconv"`,
		`"github.com/gorilla/mux"`,
		`"github.com/dimonrus/gorest"`,
		fmt.Sprintf(`"%s/app/base"`, project),
		fmt.Sprintf(`"%s/%s"`, project, "core"),
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
