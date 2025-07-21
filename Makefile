# Resume Optimizer Build and Deploy Scripts

.PHONY: build-all build-client build-auth build-processor deploy-all deploy-client deploy-services setup-db clean

# Build all services
build-all: build-client build-auth build-processor

# Build individual services
build-client:
	cd client && npm run build
	docker build -t resume-optimizer-client:latest client/

build-auth:
	docker build -t resume-optimizer-auth:latest services/auth/

build-processor:
	docker build -t resume-optimizer-processor:latest services/resume-processor/

# Deploy to Kubernetes
deploy-all: deploy-services deploy-client

deploy-client:
	helm upgrade --install resume-optimizer-client deployment/helm/client/ \
		--namespace resume-optimizer-client \
		--create-namespace

deploy-services:
	helm upgrade --install resume-optimizer-services deployment/helm/services/ \
		--namespace resume-optimizer-controller \
		--create-namespace

# Database setup
setup-db:
	./deployment/scripts/setup-db.sh

# Development
dev-client:
	cd client && npm start

dev-backend:
	docker-compose up auth-service resume-processor postgres

# Clean up
clean:
	docker system prune -f
	helm uninstall resume-optimizer-client -n resume-optimizer-client || true
	helm uninstall resume-optimizer-services -n resume-optimizer-controller || true