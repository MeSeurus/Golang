### Инструкция по запуску

1. Установить зависимости:
```go mod download```

2. Запустить сервер:
```go run cmd/server/main.go```

### Эндпоинты

| Метод | Рут              | Описание                 |
|-------|------------------|-------------------------|
| GET   | /tasks           | Получить список задач    |
| POST  | /tasks           | Создать новую задачу    |
| GET   | /tasks/{id}      | Получить задачу по ID    |
| PUT   | /tasks/{id}      | Обновить задачу         |
| DELETE| /tasks/{id}      | Удалить задачу          |

### Примеры запросов через CURL

Получить список задач:
```curl localhost:8080/tasks```

Создать новую задачу:
```curl -X POST -H 'Content-Type: application/json' -d '{"title": "Buy milk", "done": false}' localhost:8080/tasks```

Обновить существующую задачу:
```curl -X PUT -H 'Content-Type: application/json' -d '{"title": "Finish report", "done": true}' localhost:8080/tasks/1```

Удалить задачу:
```curl -X DELETE localhost:8080/tasks/1```
