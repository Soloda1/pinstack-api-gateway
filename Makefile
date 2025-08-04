.PHONY: test test-unit test-integration test-gateway-integration clean build run docker-build setup-system-tests

BINARY_NAME=api-gateway
DOCKER_IMAGE=pinstack-api-gateway:latest
GO_VERSION=1.24.2
SYSTEM_TESTS_DIR=../pinstack-system-tests
SYSTEM_TESTS_REPO=https://github.com/Soloda1/pinstack-system-tests.git

# Проверка версии Go
check-go-version:
	@echo "🔍 Проверка версии Go..."
	@go version | grep -q "go$(GO_VERSION)" || (echo "❌ Требуется Go $(GO_VERSION)" && exit 1)
	@echo "✅ Go $(GO_VERSION) найден"

# Настройка system tests репозитория
setup-system-tests:
	@echo "🔄 Проверка system tests репозитория..."
	@if [ ! -d "$(SYSTEM_TESTS_DIR)" ]; then \
		echo "📥 Клонирование pinstack-system-tests..."; \
		git clone $(SYSTEM_TESTS_REPO) $(SYSTEM_TESTS_DIR); \
	else \
		echo "🔄 Обновление pinstack-system-tests..."; \
		cd $(SYSTEM_TESTS_DIR) && git pull origin main; \
	fi
	@echo "✅ System tests готовы"

# Форматирование и проверки
fmt: check-go-version
	gofmt -s -w .
	go fmt ./...

lint: check-go-version
	go vet ./...
	golangci-lint run

# Юнит тесты
test-unit: check-go-version
	go test -v -count=1 -race -coverprofile=coverage.txt ./...

# Запуск полной инфраструктуры для интеграционных тестов из существующего docker-compose
start-gateway-infrastructure: setup-system-tests
	@echo "🚀 Запуск полной инфраструктуры для интеграционных тестов API Gateway..."
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml up -d
	@echo "⏳ Ожидание готовности всех сервисов..."
	@sleep 90

# Проверка готовности сервисов
check-services:
	@echo "🔍 Проверка готовности всех сервисов..."
	@docker exec pinstack-user-db-test pg_isready -U postgres || (echo "❌ User база данных не готова" && exit 1)
	@docker exec pinstack-auth-db-test pg_isready -U postgres || (echo "❌ Auth база данных не готова" && exit 1)
	@docker exec pinstack-post-db-test pg_isready -U postgres || (echo "❌ Post база данных не готова" && exit 1)
	@docker exec pinstack-notification-db-test pg_isready -U postgres || (echo "❌ Notification база данных не готова" && exit 1)
	@docker exec pinstack-relation-db-test pg_isready -U postgres || (echo "❌ Relation база данных не готова" && exit 1)
	@echo "✅ Все базы данных готовы"
	@echo "=== User Service logs ==="
	@docker logs pinstack-user-service-test --tail=10
	@echo "=== Auth Service logs ==="
	@docker logs pinstack-auth-service-test --tail=10
	@echo "=== Post Service logs ==="
	@docker logs pinstack-post-service-test --tail=10
	@echo "=== Notification Service logs ==="
	@docker logs pinstack-notification-service-test --tail=10
	@echo "=== Relation Service logs ==="
	@docker logs pinstack-relation-service-test --tail=10
	@echo "=== API Gateway logs ==="
	@docker logs pinstack-api-gateway-test --tail=10

# Интеграционные тесты для всех endpoints API Gateway
test-gateway-integration: start-gateway-infrastructure check-services
	@echo "🧪 Запуск интеграционных тестов для API Gateway..."
	cd $(SYSTEM_TESTS_DIR) && \
	go test -v -count=1 -timeout=15m ./internal/scenarios/integration/...

# Остановка всех контейнеров
stop-gateway-infrastructure:
	@echo "🛑 Остановка всей инфраструктуры..."
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml stop
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml rm -f

# Полная очистка (включая volumes)
clean-gateway-infrastructure:
	@echo "🧹 Полная очистка всей инфраструктуры..."
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml down -v
	@echo "🧹 Очистка Docker контейнеров, образов и volumes..."
	docker container prune -f
	docker image prune -a -f
	docker volume prune -f
	docker network prune -f
	@echo "✅ Полная очистка завершена"

# Полные интеграционные тесты (с очисткой)
test-integration: test-gateway-integration stop-gateway-infrastructure

# Все тесты
test-all: fmt lint test-unit test-integration

# Логи сервисов
logs-user:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f user-service-test

logs-auth:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f auth-service-test

logs-post:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f post-service-test

logs-notification:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f notification-service-test

logs-relation:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f relation-service-test

logs-gateway:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f api-gateway-test

logs-user-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f user-db-test

logs-auth-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f auth-db-test

logs-post-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f post-db-test

logs-notification-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f notification-db-test

logs-relation-db:
	cd $(SYSTEM_TESTS_DIR) && \
	docker compose -f docker-compose.test.yml logs -f relation-db-test

# Очистка
clean: clean-gateway-infrastructure
	go clean
	rm -f $(BINARY_NAME)
	@echo "🧹 Финальная очистка Docker системы..."
	docker system prune -a -f --volumes
	@echo "✅ Вся очистка завершена"

# Экстренная полная очистка Docker (если что-то пошло не так)
clean-docker-force:
	@echo "🚨 ЭКСТРЕННАЯ ПОЛНАЯ ОЧИСТКА DOCKER..."
	@echo "⚠️  Это удалит ВСЕ Docker контейнеры, образы, volumes и сети!"
	@read -p "Продолжить? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	docker stop $$(docker ps -aq) 2>/dev/null || true
	docker rm $$(docker ps -aq) 2>/dev/null || true
	docker rmi $$(docker images -q) 2>/dev/null || true
	docker volume rm $$(docker volume ls -q) 2>/dev/null || true
	docker network rm $$(docker network ls -q) 2>/dev/null || true
	docker system prune -a -f --volumes
	@echo "💥 Экстренная очистка завершена"

# CI локально (имитация GitHub Actions)
ci-local: test-all
	@echo "🎉 Локальный CI завершен успешно!"

# Быстрый тест (только запуск без пересборки)
quick-test: start-gateway-infrastructure
	@echo "⚡ Быстрый запуск всех интеграционных тестов без пересборки..."
	cd $(SYSTEM_TESTS_DIR) && \
	go test -v -count=1 -timeout=10m ./internal/scenarios/integration/...
	$(MAKE) stop-gateway-infrastructure