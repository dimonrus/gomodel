package gomodel

import (
	_ "embed"
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

// GenerateCrud generate model in client pack and generate crud in core
func GenerateCrud(crudPath, clientPath, project string, schema, table, tmpl string, q godb.Queryer) error {
	err := MakeModel(q, clientPath, schema, table, tmpl, DefaultSystemColumnsSoft)
	if err != nil {
		return err
	}

	// New Template
	tmp := template.New("crud")
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
	}{
		Package: packageName,
		Model:   getModelName(schema, table),
		Imports: imports,
	})

	if err != nil {
		return err
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}
