# mksite.go
Build html file from markdown folder

### Usage

create config.json

for example:

```text
{
    "sourceDir": "source/_posts/*",
    "targetDir": "website",
    "templateFile": "template/tpl.html"
}
```

```text
go run ./mksite.go
```
