# mksite.go
Build html file from markdown folder

### Usage


- create markdown files folder
- create html template
- create config.json

for example:

```text
{
    "sourceDir": "source/_posts/*",
    "targetDir": "website",
    "templateFile": "template/tpl.html"
}
```

### run 
```text
go run ./mksite.go
```
