## Model generator
*Will create model.go in app/models directory*
```
err := gomodel.MakeModel(dbo, "app/models", "schema", "table", "vendor/github.com/dimonrus/gomodel/model.tmpl", gomodel.DefaultSystemColumnsSoft)
if err != nil {
   panic(err)
}
```