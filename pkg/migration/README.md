# Database Migrations with Goose

## Instalation

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Create Migration

```bash
cd migrations # Change direectory to migrations
goose -s create [migration_name] sql
```