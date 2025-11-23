# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюеров на Pull Request'ы, управления командами и пользователями.

## Запуск

Для запуска требуется только **Docker** и **Docker Compose**.

1. Клонируйте репозиторий.
```
git clone https://github.com/Mandarkhanov/pr-reviewer-service.git
```

2. Запустите проект

```
docker-compose up --build
```

Сервис будет доступен по адресу: http://localhost:8080

## Технический стек
Язык: Go
Web Framework: Gin
База данных: PostgreSQL 15 + pgx/v5 
Миграции: golang-migrate