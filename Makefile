.PHONY: help install build test test-unit test-e2e check lint clean

# Help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*## "} /^([a-zA-Z_-]+:.*)? ## / {printf "%-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Install dependencies
install: ## Install all dependencies
	cd oms && go mod tidy
	cd frontend && npm install

# Build
build: ## Build both backend and frontend
	cd oms && go build -o server ./cmd/server
	cd frontend && npm run build

# Test
test: test-unit ## Run all tests

test-unit: ## Run unit tests
	cd oms && go test ./...
	cd frontend && npm test

test-e2e: ## Run E2E tests (requires backend running)
	cd frontend && npx playwright test

# Check
check: build test ## Run build and test

# Lint
lint: ## Run all linters
	@echo "=== Go Lint ===" && \
	if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found. Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.61.0"; \
		echo "Running gofmt check instead..."; \
		cd oms && gofmt -l . && [ -z "$$(gofmt -l .)" ] || { echo "gofmt failed: run 'cd oms && gofmt -w .'"; exit 1; }; \
	else \
		cd oms && golangci-lint run ./...; \
	fi
	@echo "=== Frontend Lint ===" && \
	cd frontend && npm run lint

# Clean
clean: ## Clean build artifacts
	cd oms && rm -f server && go clean
	cd frontend && rm -rf dist playwright-report test-results

# Backend only
backend-run: ## Run backend server
	cd oms && go run cmd/server/main.go

# Frontend only
frontend-run: ## Run frontend dev server
	cd frontend && npm run dev

# Database
db-create: ## Create SQLite database file
	cd oms && go run cmd/server/main.go --init-db

# Git
git-check: ## Check git status for untracked/ignored files
	@echo "=== Git Status ===" && git status --short
	@echo "=== Git Ignored Files ===" && git ls-files -o --exclude-standard --ignored | while read f; do echo "Ignored: $$f"; done

# Full check before commit
pre-commit: check lint git-check ## Run all checks before commit
	@echo "All checks passed!"

# Development
dev: backend-run frontend-run ## Run both backend and frontend