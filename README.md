# Pinstack API Gateway 🌐

**Pinstack API Gateway** — это единая точка входа для всех запросов к системе **Pinstack**, обеспечивающая маршрутизацию и взаимодействие между пользователями и микросервисами.

## Основные функции:
- Обработка HTTP-запросов и их маршрутизация к соответствующим gRPC-сервисам.
- Валидация запросов на уровне API Gateway.
- Предоставление единого интерфейса для взаимодействия с микросервисами.
- JWT-аутентификация и авторизация запросов.
- Обработка ошибок и возврат унифицированных HTTP-ответов.

## Технологии:
- **Go** — основной язык разработки.
- **gRPC** — для связи с микросервисами.
- **Docker** — для контейнеризации.

## CI/CD Pipeline 🚀

### GitHub Actions
Проект использует GitHub Actions для автоматического тестирования при каждом push/PR.

**Этапы CI:**
1. **Unit Tests** — юнит-тесты с покрытием кода
2. **Integration Tests** — интеграционные тесты с полной инфраструктурой всех микросервисов
3. **Auto Cleanup** — автоматическая очистка Docker ресурсов

### Makefile команды 📋

#### Основные команды разработки:
```bash
# Проверка кода и тесты
make fmt                    # Форматирование кода (gofmt)
make lint                   # Проверка кода (go vet)
make test-unit              # Юнит-тесты с покрытием
make test-integration       # Интеграционные тесты (с полной Docker инфраструктурой)
make test-all               # Все тесты: форматирование + линтер + юнит + интеграционные

# CI локально
make ci-local               # Полный CI процесс локально (имитация GitHub Actions)
```

#### Управление инфраструктурой:
```bash
# Настройка репозиториев
make setup-system-tests         # Клонирует/обновляет pinstack-system-tests репозиторий
make setup-monitoring           # Клонирует/обновляет pinstack-monitoring-service репозиторий

# Запуск инфраструктуры
make start-gateway-infrastructure  # Поднимает ВСЕ Docker контейнеры для тестов
make check-services               # Проверяет готовность всех сервисов

# Интеграционные тесты
make test-gateway-integration     # Запускает интеграционные тесты для всех endpoints
make quick-test                  # Быстрый запуск тестов без пересборки контейнеров

# Остановка и очистка
make stop-gateway-infrastructure  # Останавливает все тестовые контейнеры
make clean-gateway-infrastructure # Полная очистка (контейнеры + volumes + образы)
make clean                       # Полная очистка проекта + Docker
```

#### Мониторинг:
```bash
# Полный мониторинг стек
make start-monitoring           # Запускает Prometheus + Grafana + Loki + ELK + PgAdmin + Kafka UI
make start-prometheus-stack     # Только Prometheus + Grafana + Loki
make start-elk-stack           # Только ELK (Elasticsearch + Logstash + Kibana)
make check-monitoring-health   # Проверяет состояние всех monitoring сервисов
make stop-monitoring           # Останавливает monitoring стек
make clean-monitoring          # Очищает monitoring стек
```

#### Комбинированные команды разработки:
```bash
# Полная среда разработки
make start-dev-full            # Запускает ВСЕ сервисы + полный мониторинг
make start-dev-light           # Запускает все сервисы + только Prometheus стек (без ELK)
make stop-dev-full             # Останавливает всю dev среду
make clean-dev-full            # Полная очистка dev среды
```

#### Логи и отладка:
```bash
# Просмотр логов сервисов
make logs-user              # Логи User Service
make logs-auth              # Логи Auth Service
make logs-post              # Логи Post Service
make logs-notification      # Логи Notification Service
make logs-relation          # Логи Relation Service
make logs-gateway           # Логи API Gateway

# Просмотр логов баз данных
make logs-user-db           # Логи User Database
make logs-auth-db           # Логи Auth Database
make logs-post-db           # Логи Post Database
make logs-notification-db   # Логи Notification Database
make logs-relation-db       # Логи Relation Database

# Redis
make logs-redis             # Логи Redis
make redis-cli              # Подключение к Redis CLI
make redis-info             # Информация о Redis
make redis-keys             # Все ключи в Redis
make redis-flush            # Очистка данных Redis

# Monitoring логи
make logs-prometheus        # Логи Prometheus
make logs-grafana           # Логи Grafana
make logs-loki              # Логи Loki
make logs-elasticsearch     # Логи Elasticsearch
make logs-kibana            # Логи Kibana

# Экстренная очистка
make clean-docker-force     # Удаляет ВСЕ Docker ресурсы (с подтверждением)
```

### Зависимости для интеграционных тестов 🐳

Для интеграционных тестов автоматически поднимается полная инфраструктура:

**Базы данных:**
- **user-db-test** — PostgreSQL для User Service
- **auth-db-test** — PostgreSQL для Auth Service
- **post-db-test** — PostgreSQL для Post Service
- **notification-db-test** — PostgreSQL для Notification Service
- **relation-db-test** — PostgreSQL для Relation Service

**Миграции:**
- **user-migrator-test** — миграции User Service
- **auth-migrator-test** — миграции Auth Service
- **post-migrator-test** — миграции Post Service
- **notification-migrator-test** — миграции Notification Service
- **relation-migrator-test** — миграции Relation Service

**Микросервисы:**
- **user-service-test** — User Service
- **auth-service-test** — Auth Service
- **post-service-test** — Post Service
- **notification-service-test** — Notification Service
- **relation-service-test** — Relation Service
- **api-gateway-test** — API Gateway

**Кеш:**
- **redis** — Redis для кеширования и сессий

> 📍 **Требования:** Docker, docker-compose  
> 🚀 **Все сервисы собираются автоматически из Git репозиториев**  
> 🔄 **Репозитории `pinstack-system-tests` и `pinstack-monitoring-service` клонируются автоматически при запуске**

### Быстрый старт разработки ⚡

```bash
# 1. Проверить код
make fmt lint

# 2. Запустить юнит-тесты
make test-unit

# 3. Запустить интеграционные тесты (полная инфраструктура)
make test-integration

# 4. Или всё сразу
make ci-local

# 5. Полная среда разработки с мониторингом
make start-dev-full

# 6. Остановка и очистка после работы
make stop-dev-full
make clean
```

### Доступные URL после запуска полной среды 🌐

**Мониторинг:**
- 📊 **Prometheus**: http://localhost:9090
- 📈 **Grafana**: http://localhost:3000 (admin/admin)
- 🔍 **Loki**: http://localhost:3100
- 📋 **Kibana**: http://localhost:5601
- 💾 **Elasticsearch**: http://localhost:9200

**Базы данных и утилиты:**
- 🐧 **PgAdmin**: http://localhost:5050 (admin@admin.com/admin)
- 🔴 **Redis**:  http://localhost:6379
- 🐛 **Kafka UI**: http://localhost:9091

### Особенности 🔧

- **Полная инфраструктура:** интеграционные тесты поднимают все 5 микросервисов + базы данных + Redis
- **Мониторинг:** интегрированный стек Prometheus + Grafana + Loki + ELK для полного мониторинга
- **Docker сети:** автоматическое создание `pinstack` и `pinstack-test` сетей
- **Отключение кеша тестов:** все тесты запускаются с флагом `-count=1`
- **Тестирование всех endpoints:** интеграционные тесты покрывают все API Gateway endpoints
- **Автоочистка:** CI автоматически удаляет все Docker ресурсы после себя
- **Увеличенные таймауты:** 15 минут для интеграционных тестов из-за полной инфраструктуры
- **Логирование:** подробные логи всех сервисов для отладки
- **Redis утилиты:** встроенные команды для работы с Redis (CLI, мониторинг, очистка)
- **Комбинированные команды:** возможность запустить полную dev среду одной командой

> ✅ API Gateway готов к использованию как единая точка входа в систему Pinstack.
