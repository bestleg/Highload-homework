# README

Частично сгенерировано с помощью [Autostrada](https://autostrada.dev/).

## Getting started

```
Для запуска сервиса docker-compose up -d
```

```
Postman - collection : Losyakov-test.postman_collection.json
```

### Добавлены ручки:
```/user/get``` - без авторизации

```/user/register``` - без авторизации

```/login``` - получение авторизационного токена

## Managing SQL migrations

The `Makefile` in the project root contains commands to easily create and work with database migrations:

|     |     |
| --- | --- |
| `$ make migrations/new name=add_example_table` | Create a new database migration in the `assets/migrations` folder. |
| `$ make migrations/up` | Apply all up migrations. |
| `$ make migrations/down` | Apply all down migrations. |
| `$ make migrations/goto version=N` | Migrate up or down to a specific migration (where N is the migration version number). |
| `$ make migrations/force version=N` | Force the database to be specific version without running any migrations. |
| `$ make migrations/version` | Display the currently in-use migration version. |

Hint: You can run `$ make help` at any time for a reminder of these commands.

These `Makefile` tasks are simply wrappers around calls to the `github.com/golang-migrate/migrate/v4/cmd/migrate` tool. For more information, please see the [official documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

By default all 'up' migrations are automatically run on application startup using embeded files from the `assets/migrations` directory. You can disable this by setting the `DB_AUTOMIGRATE` environment variable to `false`.
