# Shopping list ability
The ability to manage the shopping list.

## Features
- Add/delete/list items into default or named shopping list;
- Create/delete/list named shopping list;

## Prerequisites

- Having access to [milobella.com/gitlab](https://milobella.com/gitlab/milobella)
- Having ``golang`` installed [instructions](https://golang.org/doc/install)

## Build

```bash
$ go build -o bin/ability cmd/ability/main.go
```

## Run

```bash
$ bin/ability -c config/ability.toml
```

## Requests example

#### Trigger the shopping list ability

```bash
$ curl -i -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "TRIGGER_SHOPPING_LIST"}}'
HTTP/1.1 200 OK
Date: Sun, 19 May 2019 14:54:18 GMT
Content-Length: 206
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"You have {{count}} items in your main shopping list, what do you want to do ?","params":[{"name":"count","value":0,"type":"string"}]},"auto_reprompt":true,"context":{"slot_filling":{}}}
```

#### Add item to the shopping list

```bash
$ curl -i -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "ADD_TO_LIST", "entities": [{"label": "SHOPITEM", "text": "haricots"}]}}'                                                                                                     130 ↵
HTTP/1.1 200 OK
Date: Sun, 19 May 2019 14:56:58 GMT
Content-Length: 167
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"I added {{items}} to your shopping list","params":[{"name":"items","value":["haricots"],"type":"enumerated_list"}]},"context":{"slot_filling":{}}}
```

#### Remove item from the shopping list

```bash
$ curl -i -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "REMOVE_FROM_LIST", "entities": [{"label": "SHOPITEM", "text": "haricots"}]}}'
HTTP/1.1 200 OK
Date: Sun, 19 May 2019 14:58:33 GMT
Content-Length: 170
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"I removed {{items}} from your shopping list.","params":[{"name":"items","value":["haricots"],"type":"enumerated_list"}]},"context":{"slot_filling":{}}}
```

#### Remove all items from the shopping list

> TODO: Provide example request

#### Count items on the shopping list

> TODO: Provide example request

#### List items on the shopping list

> TODO: Provide example request

## CHANGELOGS
- [Application changelog](./CHANGELOG.md)
- [Helm chart changelog](./helm/oratio/CHANGELOG.md)