.PHONY: test clean wait-for-db

# Domyślne cele
all: test

# Czyszczenie środowiska: zatrzymuje kontenery i usuwa wolumeny
clean:
	@echo "🧹 Cleaning up environment..."
	docker-compose down -v --remove-orphans
	@# Opcjonalnie: zabij procesy blokujące porty 5432 i 8080 jeśli docker-compose down nie wystarczyło
	@lsof -ti :5432 | xargs kill -9 2>/dev/null || true
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true

# Uruchomienie środowiska w tle
up: clean
	@echo "🚀 Starting environment..."
	docker-compose up -d --build

# Czekanie aż baza danych będzie gotowa
wait-for-db:
	@echo "⏳ Waiting for database to be ready..."
	@timeout=30; \
	while ! docker-compose exec -T db-service pg_isready -U postgres > /dev/null 2>&1; do \
		if [ $$timeout -le 0 ]; then \
			echo "❌ Database failed to start within timeout"; \
			exit 1; \
		fi; \
		printf "."; \
		sleep 1; \
		timeout=$$((timeout-1)); \
	done; \
	echo "\n✅ Database is ready!"

# Uruchomienie wszystkich testów (jednostkowe + integracyjne)
test: up wait-for-db
	@echo "🧪 Running tests..."
	go test ./... -v || (echo "❌ Tests failed"; $(MAKE) clean; exit 1)
	@echo "✅ All tests passed!"
	$(MAKE) clean
