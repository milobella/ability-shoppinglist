# Shopping list ability
The ability to manage the shopping list.

## Features
- Add/delete/list items into default or named shopping list;
- Create/delete/list named shopping list;

## Prerequisites

- Having access to [gitlab.milobella.com](https://gitlab.milobella.com/milobella)
- Having ``golang`` installed [instructions](https://golang.org/doc/install)
- Having ``go dep`` installed [instructions](https://golang.github.io/dep/docs/installation.html)

## Build

```bash
$ dep ensure
$ go build -o bin/ability cmd/ability/main.go
```

## Run

```bash
$ bin/ability -c config/ability.toml
```

## Requests example

#### Trigger the shopping list ability

```bash
$ curl -i -X POST http://localhost:4444/resolve/TRIGGER_SHOPPING_LIST -d '{}'
HTTP/1.1 200 OK
Date: Wed, 01 May 2019 20:55:22 GMT
Content-Length: 166
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"You have {{number}} items in your main shopping list, what do you want to do ?","params":[{"name":"number","value":2,"type":"enumerated_list"}]}}
```

#### Add item to the shopping list

```bash
$ curl -i -X POST http://localhost:4444/resolve/ADD_TO_LIST -d '{"nlu": {"entities": [{"label": "SHOPITEM", "text": "du pain"}]}}'
HTTP/1.1 200 OK
Date: Wed, 01 May 2019 21:03:11 GMT
Content-Length: 136
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"I added {{items}} to your shopping list","params":[{"name":"items","value":["du pain"],"type":"enumerated_list"}]}}
```

## CHANGELOGS
- [Application changelog](./CHANGELOG.md)
- [Helm chart changelog](./helm/oratio/CHANGELOG.md)