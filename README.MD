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

#### If you find this project useful or want to support the author, you can send tokens to any of these wallets
- Bitcoin: bc1qgx5c3n7q26qv0tngculjz0g78u6mzavy2vg3tf
- Ethereum: 0x62812cb089E0df31347ca32A1610019537bbFe0D
- Dogecoin: DET7fbNzZftp4sGRrBehfVRoi97RiPKajV