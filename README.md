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

## CHANGELOGS
- [Application changelog](./CHANGELOG.md)
- [Helm chart changelog](./helm/oratio/CHANGELOG.md)