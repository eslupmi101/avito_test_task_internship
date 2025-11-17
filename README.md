# Avito Test Task Internship Service

## Описание

Сервис предназначен для управления командами, пользователями и pull request-ами.

Проект построен с упрощённым подходом **DDD (Domain-Driven Design)**, где:

* **Domain** — анемичные структуры (`PullRequest`, `User`) представляют агрегаты и сущности базы данных.
* **PullRequest** — агрегат таблиц `pull_requests` и `pull_request_authors`.
* **User** — повторяет таблицу `users`.
* Для простоты **контрактов/интерфейсов нет**, предполагается только одна база данных.

В слое **presentation** реализован HTTP API, слой **application** управляет сервисами и инициализацией зависимостей, а слой **domain/service** содержит бизнес-логику

---

## Используемые технологии

* Go 1.25
* [Chi](https://github.com/go-chi/chi) — легковесный HTTP router
* PostgreSQL (через `pgxpool`)
* DDD архитектура (анемичные структуры, агрегаты)
* Миграции с жесткими ограничениями через триггеры и индексы

---

## Архитектура

### Domain

Содержит структуры, отражающие таблицы базы данных:

### Application

Инициализирует сервисы и управляет их зависимостями:

```go
var (
	TeamServiceInstance        *domain_service.TeamService
	UserServiceInstance        *domain_service.UserService
	PullRequestServiceInstance *domain_service.PullRequestService
)

func InitRegistry(database *infrasctucture.PostgresDb) {
	TeamServiceInstance = domain_service.NewTeamService(database)
	UserServiceInstance = domain_service.NewUserService(database)
	PullRequestServiceInstance = domain_service.NewPullRequestService(database)
}
```

### Presentation

HTTP API реализован через Chi и контроллеры:

```go
func RegisterRoutes(r chi.Router) {
	r.Post("/team/add", controller.AddTeam)
	r.Get("/team/{team_name}", controller.GetTeam)

	r.Post("/users/setIsActive", controller.SetIsActive)
	r.Get("/users/getReview/{user_id}", controller.GetReview)

	r.Post("/pullRequest/create", controller.Create)
	r.Post("/pullRequest/merge", controller.Merge)
	r.Post("/pullRequest/reassigne", controller.Reassign)
}
```

---

## Миграции и ограничения базы данных

База данных имеет строгие бизнес-ограничения, реализованные через триггеры:

```sql
CREATE TABLE pull_request_reviewers (
    id SERIAL PRIMARY KEY,
    pull_request VARCHAR(50) NOT NULL REFERENCES pull_requests(id),
    reviewer VARCHAR(50) NOT NULL REFERENCES users(user_id)
);

-- Ограничение: максимум 2 ревьюеров на PR
-- Ограничение: ревьюер не может быть автором PR
-- Уникальность: один и тот же reviewer не может быть добавлен дважды
-- Мультиколоночный индекс для ускорения поиска
```

Триггеры проверяют:

* Максимальное количество ревьюеров на PR (`check_max_reviewers`)
* Ревьюер не совпадает с автором PR (`check_reviewer_not_author`)
* Уникальность reviewer для одного PR
* Индекс для быстрого поиска по сочетанию pull_request + reviewer

---

## Запуск

1. Настроить `.env` из .env_example с данными для подключения к PostgreSQL.


```bash
make docker-compose-infra-up
make run
```

4. HTTP API доступен на `http://localhost:8080`.
