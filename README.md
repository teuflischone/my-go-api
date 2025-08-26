# my-go-api

Небольшой сервис на Go (Echo) с базовой инфраструктурой.

## Что реализовано
- /healthz (проверка доступности БД/Redis через переменные окружения)
- /metrics (Prometheus)
- Request ID middleware
- Makefile (tidy/fmt/test/run)
- CI (Go vet + tests)

## Запуск локально
- переменные окружения — в .env.example
- команды разработки — см. Makefile

## Технологии
Go, Echo, PostgreSQL, Redis, Prometheus.
