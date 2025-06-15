# 📝 Мини-блог / Сервис заметок с пользователями

REST API для создания, хранения и получения заметок, привязанных к конкретным пользователям. Авторизация осуществляется через JWT.

---

## 📦 Технологии

- Go
- PostgreSQL
- Docker + Docker Compose
- Chi (роутер)
- JWT (аутентификация)
- Goose (миграции)
- Postman (тестирование API)

---

## 🚀 Бысткий старт

### 1. Клонируем проект и переходим в папку

```bash
git clone https://github.com/AlTrubinov/mini-blog.git
cd mini-blog
```

---

### 2. Конфигурация

#### `.env` для запуска контейнеров
Скопируйте пример .env-example в .env:

```bash
cp docker/.env-example docker/.env
```
#### `config/<your_config_name>.yaml` для конфигурации приложения
Файл `config/default-example.yaml` для локального использования можно использовать как есть — он уже указан в `.env-example`.
В случае использования конфига с другим названием, требуется указать актуальный путь до конфига в `docker/.env`.

---

### 3. Запуск приложения

```bash
docker compose -f docker/docker-compose.dev.yml up --build
```

---

### 4. Применение миграций

Миграции применяются вручную с использованием `goose`.

```bash
goose --dir ./internal/storage/migrations/ postgres "postgres://<user>:<password>@<host>:<port>/<name>?sslmode=disable" up
```

Значения `user`, `password`, `host`, `port`, `name` должны соответствовать вашему конфигу в папке `config`. Например, для `config/default-example.yaml`:

```bash
goose --dir ./internal/storage/migrations/ postgres "postgres://postgres:postgres@db:5432/mini-blog?sslmode=disable" up
```

В случае проблем с хостом, можно попробовать указать localhost:

```bash
goose --dir ./internal/storage/migrations/ postgres "postgres://postgres:postgres@localhost:5432/mini-blog?sslmode=disable" up
```

---

## 🔐 Аутентификация

| Метод | Путь       | Описание                 |
| ----- | ---------- |--------------------------|
| POST  | `/users`   | Регистрация и выдача JWT |
| POST  | `/login`   | Авторизация и выдача JWT |
| POST  | `/refresh` | Обновление JWT токена    |

Все защищённые маршруты требуют JWT access токен.

---

## 🧪 Тестирование

В папке `tests/` лежит коллекция и окружение Postman:

- `tests/mini-blog.postman_collection.json` — коллекция запросов и тестовые запросы
- `tests/mini-blog.postman_environment.json` — окружение

Вы можете использовать их в Postman для тестирования API вручную (пока только вручную)

---