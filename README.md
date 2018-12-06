# mksite.go
Build html file from markdown folder

### Usage


- create markdown files folder
- create html template
- create config.json

for example:

```text
{
    "siteName": "site name here",
    "staticDir": "static",
    "sourceDir": "source/_posts/*",
    "targetDir": "website",
    "pageSize": 30,
    "templateFile": "template/tpl.html",
    "indexTemplateFile": "template/index.html",
    "itemTemplateFile": "template/item.html"
}
```

### run 
```text
go run ./mksite.go
```
