package gomodel

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/dimonrus/godb/v2"
	"github.com/dimonrus/gohelp"
	"io"
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

//go:embed api_search.tmpl
var ApiSearchTemplate string

//go:embed search_form.tmpl
var SearchFormTemplate string

//go:embed api_route.tmpl
var ApiRouteTemplate string

//go:embed api_client.tmpl
var ApiClientTemplate string

// CrudNumber calculate for type of crud method
type CrudNumber uint8

// PossibleCrudMethods All possible crud methods
type PossibleCrudMethods struct {
	// Create method
	Create bool
	// Read method
	Read bool
	// Update method
	Update bool
	// Delete method
	Delete bool
	// Search method
	Search bool
}

// GetPossibleMethods Calculate possible methods
func (n CrudNumber) GetPossibleMethods() PossibleCrudMethods {
	return PossibleCrudMethods{
		Create: n&1 == 1,
		Read:   n&2 == 2,
		Update: n&4 == 4,
		Delete: n&8 == 8,
		Search: n&16 == 16,
	}
}

// GetPossibleMethodsArray return list of crud method based on num
func (n CrudNumber) GetPossibleMethodsArray() []string {
	var result = make([]string, 0)
	if n&1 == 1 {
		result = append(result, "create")
	}
	if n&2 == 2 {
		result = append(result, "read")
	}
	if n&4 == 4 {
		result = append(result, "update")
	}
	if n&8 == 8 {
		result = append(result, "delete")
	}
	if n&16 == 16 {
		result = append(result, "list")
	}
	return result
}

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
		`"github.com/dimonrus/gorest"`,
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

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("Core %s file created: %s", getModelName(schema, table), path)
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeSearchForm generate search form
func (c CRUDGenerator) MakeSearchForm(q godb.Queryer, schema, table string) error {
	// New Template
	tmp := template.New("search_form").Funcs(getHelperFunc(DefaultSystemColumns))
	// init template
	tmlString := SearchFormTemplate

	// Get package name
	packageNames := strings.Split(c.ClientPath, string(os.PathSeparator))
	var packageName string
	if len(packageNames) > 0 {
		packageName = packageNames[len(packageNames)-1]
	} else {
		packageName = c.ClientPath
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

	file, path, err := CreateModelFile("public", table+"_search", c.ClientPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// crud imports
	var imports = []string{
		`"github.com/dimonrus/gosql"`,
		`"github.com/lib/pq"`,
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

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("Client %s search file created: %s", getModelName(schema, table), path)
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

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API read %s file created: %s", getModelName(schema, table), path)
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

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API delete %s file created: %s", getModelName(schema, table), path)
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
	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API update %s file created: %s", getModelName(schema, table), path)
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

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API create %s file created: %s", getModelName(schema, table), path)
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPISearch generate qpi search
func (c CRUDGenerator) MakeAPISearch(q godb.Queryer, schema, table, version string) error {
	// New Template
	tmp := template.New("api_search").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiSearchTemplate

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("search", c.APIPath+string(os.PathSeparator)+table+string(os.PathSeparator)+version)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"github.com/dimonrus/gorest"`,
		fmt.Sprintf(`"%s/app/base"`, c.ProjectPath),
		fmt.Sprintf(`"%s/%s"`, c.ProjectPath, c.ClientPath),
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

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API search %s file created: %s", getModelName(schema, table), path)
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPIRoute generate qpi route
func (c CRUDGenerator) MakeAPIRoute(q godb.Queryer, schema, table, version string, num CrudNumber) error {
	// New Template
	tmp := template.New("api_route").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiRouteTemplate

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("route", c.APIPath+string(os.PathSeparator)+table+string(os.PathSeparator)+version)
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"net/http"`,
		`"github.com/gorilla/mux"`,
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Model   string
		Imports []string
		Columns Columns
		PossibleCrudMethods
	}{
		Package:             version,
		Model:               getModelName(schema, table),
		Imports:             imports,
		Columns:             *columns,
		PossibleCrudMethods: num.GetPossibleMethods(),
	})

	if err != nil {
		return err
	}

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API route %s file created: %s", getModelName(schema, table), path)
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// MakeAPIClient generate qpi search
func (c CRUDGenerator) MakeAPIClient(q godb.Queryer, schema, table, version string, num CrudNumber) error {
	// update if file already exists
	path := c.ClientPath + string(os.PathSeparator) + "api_client.go"
	f, _ := os.Open(path)
	if f != nil {
		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		_ = f.Close()
		return c.UpdateAPIClient(q, path, content, schema, table, num)
	}

	// New Template
	tmp := template.New("api_client").Funcs(getHelperFunc(DefaultSystemColumns))

	tmlString := ApiClientTemplate

	tmp = template.Must(tmp.Parse(tmlString))

	file, path, err := CreateFile("api_client", c.ClientPath+string(os.PathSeparator))
	if err != nil {
		return err
	}
	defer file.Close()

	var imports = []string{
		`"fmt"`,
		`"github.com/dimonrus/goreq"`,
		`"github.com/dimonrus/gorest"`,
		`"github.com/dimonrus/porterr"`,
		`"net/http"`,
	}

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	// Parse template to file
	err = tmp.Execute(file, struct {
		Package string
		Imports []string
		Model   string
		PossibleCrudMethods
		Columns Columns
	}{
		Package:             "client",
		Imports:             imports,
		Model:               getModelName(schema, table),
		PossibleCrudMethods: num.GetPossibleMethods(),
		Columns:             *columns,
	})

	if err != nil {
		return err
	}

	if dbo, ok := q.(*godb.DBO); ok {
		dbo.Logger.Printf("API client file created: %s", path)
	}

	// Format code
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

func (c CRUDGenerator) UpdateAPIClient(q godb.Queryer, path string, content []byte, schema, table string, num CrudNumber) error {
	pcm := num.GetPossibleMethodsArray()
	modelName := getModelName(schema, table)

	createApi := `
    // Create{{ .Model }} Create {{ icameled .Model }} http method
	func (s Service) Create{{ .Model }}(request {{ .Model }}) ({{ icameled .Model }} {{ .Model }}, e porterr.IError) {
		response := gorest.JsonResponse{Data: &{{ icameled .Model }}}
		_, err := s.EnsureJSON(http.MethodPost, "api/v1/{{ icameled .Model }}", nil, request, &response)
		if err != nil {
			e = err.(*porterr.PortError)
		}
    	return
	}`
	readApi := `
    // Read{{ .Model }} Read {{ icameled .Model }} http method
	func (s Service) Read{{ .Model }}({{ $index := 0 }}{{ range $key, $column := .Columns }}{{ if $column.IsPrimaryKey }}{{ if $index }}, {{ end }}{{ $index = inc $index }} {{ icameled $column.Name }} {{ $column.ModelType }} {{ end }}{{ end }}) ({{ icameled .Model }} {{ .Model }}, e porterr.IError) {
		response := gorest.JsonResponse{Data: &{{ icameled .Model }}}
		_, err := s.EnsureJSON(http.MethodGet, fmt.Sprintf("api/v1/{{ icameled .Model }}/%v", id), nil, nil, &response)
		if err != nil {
			e = err.(*porterr.PortError)
		}
		return
	}`
	updateApi := `
   	// Update{{ .Model }} Update user http method
	func (s Service) Update{{ .Model }}({{ $index := 0 }}{{ range $key, $column := .Columns }}{{ if $column.IsPrimaryKey }}{{ if $index }}, {{ end }}{{ $index = inc $index }} {{ icameled $column.Name }} {{ $column.ModelType }} {{ end }}{{ end }}, request {{ .Model }}) ({{ icameled .Model }} {{ .Model }}, e porterr.IError) {
		response := gorest.JsonResponse{Data: &{{ icameled .Model }}}
		_, err := s.EnsureJSON(http.MethodPatch, fmt.Sprintf("api/v1/{{ icameled .Model }}/%v", id), nil, request, &response)
		if err != nil {
			e = err.(*porterr.PortError)
		}
		return
	}`
	deleteApi := `
   	// Delete{{ .Model }} Delete {{ icameled .Model }} http method
	func (s Service) Delete{{ .Model }}({{ $index := 0 }}{{ range $key, $column := .Columns }}{{ if $column.IsPrimaryKey }}{{ if $index }}, {{ end }}{{ $index = inc $index }} {{ icameled $column.Name }} {{ $column.ModelType }} {{ end }}{{ end }}) (e porterr.IError) {
		_, err := s.EnsureJSON(http.MethodDelete, fmt.Sprintf("api/v1/{{ icameled .Model }}/%v", id), nil, nil, nil)
		if err != nil {
			e = err.(*porterr.PortError)
		}
		return
	}`
	listApi := `
	// List{{ .Model }} Get list of {{ searchField (icameled .Model) }} http method
	func (s Service) List{{ .Model }}(form {{ .Model }}SearchForm) (list {{ searchField .Model }}, meta gorest.Meta, e porterr.IError) {
		response := gorest.JsonResponse{Data: &list, Meta: &meta}
		_, err := s.EnsureJSON(http.MethodPost, "api/v1/{{ icameled .Model }}/list", nil, form, &response)
		if err != nil {
			e = err.(*porterr.PortError)
		}
		return
	}`
	crudMap := map[string]string{"create": createApi, "read": readApi, "update": updateApi, "delete": deleteApi, "list": listApi}

	// define all needed api
	reader := bufio.NewReader(bytes.NewBuffer(content))
	foundMethods := make([]string, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		for _, s := range pcm {
			methodName := gohelp.ToCamelCase(s+modelName, true) + "("
			if strings.Contains(string(line), methodName) {
				foundMethods = append(foundMethods, s)
				break
			}
		}
	}

	var err error
	var needleContent string
	tml := template.New("needed_api").Funcs(getHelperFunc(DefaultSystemColumnsSoft))

	for _, s := range pcm {
		if !gohelp.ExistsInArray(s, foundMethods) {
			needleContent += crudMap[s]
		}
	}

	// Columns
	columns, err := GetTableColumns(q, schema, table, DefaultSystemColumns, getDictionaryItems(q))
	if err != nil {
		return err
	}
	if columns == nil || len(*columns) == 0 {
		return errors.New("No table found or no columns in table ")
	}

	tml, err = tml.Parse(needleContent)
	if err != nil {
		return err
	}
	var data = make([]byte, 0)
	buf := bytes.NewBuffer(data)
	var str = bufio.NewWriter(buf)
	err = tml.Execute(str, struct {
		Model   string
		Columns Columns
	}{
		Model:   modelName,
		Columns: *columns,
	})
	if err != nil {
		return err
	}
	str.Flush()
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// AddToGlobalRoute add global route for target entity
func (c CRUDGenerator) AddToGlobalRoute(schema, table, version string) error {
	path := c.APIPath + string(os.PathSeparator) + "../route.go"
	var content []byte
	var err error
	content, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(bytes.NewReader(content))

	var newContent strings.Builder
	var alias = gohelp.ToCamelCase(getModelName(schema, table), false) + strings.ToUpper(version)
	var initString = fmt.Sprintf("%s.Init(ApiRoute%s)", alias, strings.ToUpper(version))
	var versionSubRoute bool
	var importVersion bool
	var subRouteAdded = strings.Contains(string(content), initString)
	for {
		var line []byte
		line, _, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if strings.Contains(string(line), "ApiRoute"+strings.ToUpper(version)) {
			versionSubRoute = true
		}
		if strings.Contains(string(line), alias) {
			importVersion = true
		}
		if strings.Contains(string(line), "net/http") && !importVersion {
			importString := alias + " \"" + c.ProjectPath + "/" + c.APIPath + "/" + gohelp.ToUnderscore(getModelName(schema, table)) + "/" + version + "\"\n\n"
			newContent.WriteString(importString)
		} else if strings.Contains(string(line), "Setup middleware") {
			if !versionSubRoute {
				subRoute := fmt.Sprintf("// Api %s routes \n ApiRoute%s := ApiRoute.PathPrefix(\"%s\").Subrouter() \n\n", version, strings.ToUpper(version), "/"+version)
				newContent.WriteString(subRoute)
			}
			if !subRouteAdded {
				newContent.WriteString(fmt.Sprintf("// %s sub route \n %s.Init(ApiRoute%s) \n\n", getModelName(schema, table), alias, strings.ToUpper(version)))
			}
		}
		newContent.WriteString(string(line) + "\n")
	}
	f, err := os.OpenFile(path, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(newContent.String())
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "fmt", path)
	return cmd.Run()
}

// Generate generate crud, client, api
// q - database connection
// schema - db schema (table namespace)
// table - name of table
// num - crud scenario (1 - create, 2 - read, 4 - update, 8 - delete, 16 - list)
func (c CRUDGenerator) Generate(q godb.Queryer, schema, table, version string, num CrudNumber) error {
	modelTemplate := DefaultModelTemplate
	err := MakeModel(q, c.ClientPath, schema, table, modelTemplate, DefaultSystemColumnsSoft)
	if err != nil {
		return err
	}
	err = c.MakeCoreCrud(q, schema, table)
	if err != nil {
		return err
	}
	if num&1 == 1 {
		err = c.MakeAPICreate(q, schema, table, version)
		if err != nil {
			return err
		}
	}
	if num&2 == 2 {
		err = c.MakeAPIRead(q, schema, table, version)
		if err != nil {
			return err
		}
	}
	if num&4 == 4 {
		err = c.MakeAPIUpdate(q, schema, table, version)
		if err != nil {
			return err
		}
	}
	if num&8 == 8 {
		err = c.MakeAPIDelete(q, schema, table, version)
		if err != nil {
			return err
		}
	}
	err = c.MakeSearchForm(q, schema, table)
	if err != nil {
		return err
	}
	if num&16 == 16 {
		err = c.MakeAPISearch(q, schema, table, version)
		if err != nil {
			return err
		}
	}
	if num > 0 {
		err = c.MakeAPIRoute(q, schema, table, version, num)
		if err != nil {
			return err
		}
		err = c.MakeAPIClient(q, schema, table, version, num)
		if err != nil {
			return err
		}
		err = c.AddToGlobalRoute(schema, table, version)
		if err != nil {
			return err
		}
	}
	return err
}
