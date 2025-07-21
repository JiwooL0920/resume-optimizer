# Resume Optimizer Build and Deploy Scripts

.PHONY: help build-all build-client build-auth build-processor deploy-all deploy-client deploy-services setup-db clean \
	docker-build docker-up docker-down down docker-logs docker-logs-auth docker-logs-resume docker-logs-client \
	health test-upload test-list setup dev restart rebuild

# Default target
help: ## Show this help message
	@echo "Resume Optimizer - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build all services
build-all: build-client build-auth build-processor ## Build all services for K8s deployment

# Build individual services
build-client: ## Build React client Docker image
	cd client && npm run build
	docker build -t resume-optimizer-client:latest client/

build-auth: ## Build auth service Docker image
	docker build -t resume-optimizer-auth:latest services/auth/

build-processor: ## Build resume processor Docker image
	docker build -t resume-optimizer-processor:latest services/resume-processor/

# Deploy to Kubernetes
deploy-all: deploy-services deploy-client ## Deploy all services to K8s

deploy-client: ## Deploy client to K8s
	helm upgrade --install resume-optimizer-client deployment/helm/client/ \
		--namespace resume-optimizer-client \
		--create-namespace

deploy-services: ## Deploy backend services to K8s
	helm upgrade --install resume-optimizer-services deployment/helm/services/ \
		--namespace resume-optimizer-controller \
		--create-namespace

# Database setup
setup-db: ## Setup database
	./deployment/scripts/setup-db.sh

# Development
dev-client: ## Start React client in development mode
	cd client && npm start

dev-backend: ## Start backend services with Docker Compose (deprecated - use docker-up)
	docker-compose up auth-service resume-processor postgres

# Docker Compose commands (for local development with external k8s PostgreSQL)
docker-build: ## Build Docker images with Docker Compose
	@echo "Building Docker images..."
	@docker-compose --env-file .env.docker build

docker-up: ## Start services with Docker Compose
	@echo "Starting Resume Optimizer with Docker Compose..."
	@if [ ! -f .env.docker ]; then \
		echo "Error: .env.docker file not found!"; \
		echo "Please copy .env.docker.example to .env.docker and update with your credentials:"; \
		echo "  cp .env.docker.example .env.docker"; \
		echo "  # Then edit .env.docker with your actual database credentials"; \
		exit 1; \
	fi
	@docker-compose --env-file .env.docker up -d
	@echo "Waiting for services to start..."
	@sleep 5
	@echo ""
	@echo "Services started successfully!"
	@$(MAKE) health
	@echo ""
	@echo "Available endpoints:"
	@echo "  Frontend (React):       http://localhost:3000"
	@echo "  Auth Service:           http://localhost:8080"
	@echo "  Resume Processor:       http://localhost:8081"

docker-down: ## Stop Docker Compose services
	@echo "Stopping Docker services..."
	@docker-compose --env-file .env.docker down

down: docker-down ## Alias for docker-down (stop services)

docker-logs: ## View logs from all Docker services
	@docker-compose --env-file .env.docker logs -f

docker-logs-auth: ## View logs from auth service only
	@docker-compose --env-file .env.docker logs -f auth-service

docker-logs-resume: ## View logs from resume processor service only
	@docker-compose --env-file .env.docker logs -f resume-processor

docker-logs-client: ## View logs from client service only
	@docker-compose --env-file .env.docker logs -f client

# Testing and health checks
health: ## Check service health
	@echo "Checking service health..."
	@echo "Auth Service:"
	@curl -s http://localhost:8080/health 2>/dev/null | jq '.' || echo "  ❌ Service not responding"
	@echo ""
	@echo "Resume Processor Service:"
	@curl -s http://localhost:8081/health 2>/dev/null | jq '.' || echo "  ❌ Service not responding"

test-upload: ## Test resume upload functionality
	@echo "Testing resume upload..."
	@if [ ! -f sample_resume.pdf ]; then \
		echo "Creating sample PDF for testing..."; \
		curl -s -o sample_resume.pdf "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"; \
	fi
	@curl -X POST http://localhost:8081/api/v1/resumes/upload -F "file=@sample_resume.pdf" | jq '.'

test-list: ## Test listing resumes
	@echo "Testing resume list..."
	@curl -s http://localhost:8081/api/v1/resumes/ | jq '.'

# Setup and development commands
setup: ## Set up development environment (copy env template)
	@if [ ! -f .env.docker ]; then \
		echo "Setting up development environment..."; \
		cp .env.docker.example .env.docker; \
		echo "✅ Created .env.docker from template"; \
		echo "⚠️  Please edit .env.docker with your actual credentials before running 'make dev'"; \
		echo "   - Update database credentials (DB_USER, DB_PASSWORD, DB_NAME)"; \
		echo "   - Add Google OAuth credentials if needed"; \
	else \
		echo "✅ .env.docker already exists"; \
	fi

dev: docker-up ## Start development environment with Docker Compose
	@echo "Development environment ready!"

restart: docker-down docker-up ## Restart all Docker services

rebuild: docker-down docker-build docker-up ## Rebuild and restart services

# Clean up
clean: ## Clean up Docker containers and K8s deployments
	@echo "Cleaning up..."
	@docker-compose --env-file .env.docker down --volumes --remove-orphans 2>/dev/null || true
	@docker system prune -f
	@helm uninstall resume-optimizer-client -n resume-optimizer-client 2>/dev/null || true
	@helm uninstall resume-optimizer-services -n resume-optimizer-controller 2>/dev/null || true
