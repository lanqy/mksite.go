# mksite.go

Build html file from markdown folder

## Usage

-   create markdown files folder
-   create html template
-   create config.json

for example:

```text
{
    "siteName": "site name here",
    "staticDir": "static",
    "baseUrl": "https://lanqy.xyz",
    "sourceDir": "source/_posts/*",
    "targetDir": "website",
    "postTemplateFile": "template/post.html",
    "navTemplateFile": "template/nav.html",
    "indexTemplateFile": "template/index.html",
    "itemTemplateFile": "template/item.html"
}
```

## run

```text
go run ./mksite.go
```

## or build .exe

```text
go build -o mksite.exe ./mksite.go
```

then you got a website on website folder :)

## Use node.js ?

----> https://github.com/lanqy/mksite
