# Blog Platform

Расширенная система управления блогом на Go.

## Стек
- Go 1.21+
- PostgreSQL
- Docker Compose
- JWT аутентификация
- chi router
- sqlx
- bcrypt

## Запуск

1. Клонируйте репозиторий.
2. Создайте `.env` на основе `.env.example`.
3. Запустите базу данных: docker-compose up -d
4. Запустите приложение: go run cmd/api/main.go или через Make: make run

## API Эндпоинты

### Auth
- `POST /api/register` — регистрация (email, password)
- `POST /api/login` — вход, возвращает JWT

### Posts
- `GET /api/posts?limit=10&offset=0` — список опубликованных постов
- `GET /api/posts/{id}` — детальный пост
- `POST /api/posts` — создать пост (требуется авторизация)
- `PUT /api/posts/{id}` — обновить пост (только автор)
- `DELETE /api/posts/{id}` — удалить пост (только автор)

### Comments
- `POST /api/posts/{postId}/comments` — добавить комментарий (авторизация)
- `GET /api/posts/{postId}/comments` — список комментариев поста

### Health
- `GET /api/health`

## Тестирование
make test

## Примечания
- Планировщик публикует посты с `publish_at <= now()`.
- Используется bcrypt для паролей.
- Graceful shutdown обрабатывает завершение фоновых воркеров.