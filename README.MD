## Model generator
*Will create model.go in app/models directory*
```
err := godb.MakeModel(dbo, "app/models", "schema", "table", "vendor/github.com/dimonrus/godb/model.tmpl", godb.DefaultSystemColumnsSoft)
if err != nil {
   panic(err)
}
```