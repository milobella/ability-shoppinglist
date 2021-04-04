# Shopping list Ability
Milobella Ability to manage the shopping list.

## Features
- Add/delete/list items into default or named shopping list;
- Create/delete/list named shopping list;

## Prerequisites

- Having ``golang`` installed [instructions](https://golang.org/doc/install)

## Build

```bash
$ go build -o bin/ability cmd/ability/main.go
```

## Run

```bash
$ bin/ability
```

## Requests example

#### Trigger the shopping list ability

```bash
$ curl -i -H "Content-Type":"application/json" -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "TRIGGER_SHOPPING_LIST"}}'
HTTP/1.1 200 OK
Date: Sun, 19 May 2019 14:54:18 GMT
Content-Length: 206
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"You have {{count}} items in your main shopping list, what do you want to do ?","params":[{"name":"count","value":0,"type":"string"}]},"auto_reprompt":true,"context":{"slot_filling":{}}}
```

#### Add item to the shopping list

```bash
$ curl -i -H "Content-Type":"application/json" -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "ADD_TO_LIST", "entities": [{"label": "SHOPITEM", "text": "haricots"}]}}'                                                                                                     130 ↵
HTTP/1.1 200 OK
Date: Sun, 19 May 2019 14:56:58 GMT
Content-Length: 167
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"I added {{items}} to your shopping list","params":[{"name":"items","value":["haricots"],"type":"enumerated_list"}]},"context":{"slot_filling":{}}}
```

#### Remove item from the shopping list

```bash
$ curl -i -H "Content-Type":"application/json" -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "REMOVE_FROM_LIST", "entities": [{"label": "SHOPITEM", "text": "haricots"}]}}'
HTTP/1.1 200 OK
Date: Sun, 19 May 2019 14:58:33 GMT
Content-Length: 170
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"I removed {{items}} from your shopping list.","params":[{"name":"items","value":["haricots"],"type":"enumerated_list"}]},"context":{"slot_filling":{}}}
```

#### Remove all items from the shopping list

```bash
$ curl -i -H "Content-Type":"application/json" -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "EMPTY_LIST_ITEMS"}}'
HTTP/1.1 200 OK
Date: Wed, 30 Oct 2019 18:43:16 GMT
Content-Length: 90
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"Your shopping list has been cleared."},"context":{"slot_filling":{}}}

```

#### Count items on the shopping list

```bash
$ curl -i -H "Content-Type":"application/json" -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "COUNT_LIST_ITEMS"}}'
HTTP/1.1 200 OK
Date: Wed, 30 Oct 2019 18:42:06 GMT
Content-Length: 155
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"You have {{count}} items in your shopping list.","params":[{"name":"count","value":2,"type":"string"}]},"context":{"slot_filling":{}}}
```

#### List items on the shopping list

```bash
$ curl -i -H "Content-Type":"application/json" -X POST http://localhost:4444/resolve -d '{"nlu":{"BestIntent": "LIST_LIST_ITEMS"}}'
HTTP/1.1 200 OK
Date: Wed, 30 Oct 2019 18:39:44 GMT
Content-Length: 250
Content-Type: text/plain; charset=utf-8

{"nlg":{"sentence":"You have {{count}} items in your shopping list. There are {{items}}.","params":[{"name":"count","value":2,"type":"string"},{"name":"items","value":["haricots","haricots"],"type":"enumerated_list"}]},"context":{"slot_filling":{}}}
```
