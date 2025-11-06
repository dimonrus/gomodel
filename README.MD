## Model and crud generator

### This package contains logic to work with db models

*Create model in specific directory*
```
_, _, err := gomodel.MakeModel(base.App.GetDB(), "app/io/db/models", schema, table, "", gomodel.DefaultSystemColumnsSoft)
if err != nil {
    base.App.FatalError(err)
}
```

*Create dictionary table*
```
err := CreateDictionaryTable(base.App.GetDB())
if err != nil {
    base.App.FatalError(err)
}
```

*Example of generated model*

```
// Dictionary model
type Dictionary struct {
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

// Table get Dictionary model table name
func (m *Dictionary) Table() string {
	return "dictionary"
}

// Columns get all Dictionary model columns
func (m *Dictionary) Columns() []string {
	return []string{"id", "type", "code", "label", "created_at", "updated_at", "deleted_at"}
}

// Values get all Dictionary model values
func (m *Dictionary) Values() []any {
	return []any{&m.Id, &m.Type, &m.Code, &m.Label, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt}
}

// NewDictionary init Dictionary model method
func NewDictionary() *Dictionary {
	return &Dictionary{}
}

// NewDictionaryCollection init Dictionary collection
func NewDictionaryCollection() *gomodel.Collection[Dictionary] {
	return gomodel.NewCollection[Dictionary]()
}
```
You can use collection for load, save, and delete collection items in mass action scenarios. Example
```
// Update items
c := NewDictionaryCollection()
c.Map(func(item *Dictionary) {
    if item.Label != nil {
        item.Label = nil
    }
})
c.Save(db)

// Load dictionary collection where label is not null and updated_at >= "2022-01-01"
c := NewDictionaryCollection()
c.Where().AddExpression("label IS NOT NULL")
c.Where().AddExpression("updated_at >= ?", time.Date(2022, 1, 1, 0,0,0,1, time.UTC))
e := c.Load(db)
if e != nil {
    return
}

// Delete if found items with such condition
if c.Count() > 0 {
    e = c.Delete(db)
    if e != nil {
        return
    }
}
```


*Create dictionary mapping*
```
err := GenerateDictionaryMapping("app/io/db/models/dictionary_mapping.go", base.App.GetDB())
if err != nil {
    base.App.FatalError(err)
}
```

*Create CRUD operations and http client with swagger docs*
```
crud := gomodel.NewCRUDGenerator("app/core", "app/client", "app/io/web/api", "gost")
err := crud.Generate(base.App.GetDB(), schema, table, "v1", gomodel.CrudNumber(args["num"].GetInt()))
if err != nil {
    base.App.FatalError(err)
}
```

*Model wizard allows to set model fields step by step*
```
wiz := NewWizard[DefaultCustomModel]("name", "pages", "some_int")
wiz.Set(&sql.NullString{String: "foo"})
wiz.Set(&pq.StringArray{"foo", "bar"})
val := gohelp.Ptr[float64](123)
wiz.Set(&val)

model := wiz.Get()

// model.Name == "foo" {
// model.Pages == []string{"foo", "bar"}
// model.SomeInt == 123
// wiz.IsComplete() == true
```

#### If you find this project useful or want to support the author, you can send tokens to any of these wallets
- Bitcoin: bc1qgx5c3n7q26qv0tngculjz0g78u6mzavy2vg3tf
- Ethereum: 0x62812cb089E0df31347ca32A1610019537bbFe0D
- Dogecoin: DET7fbNzZftp4sGRrBehfVRoi97RiPKajV