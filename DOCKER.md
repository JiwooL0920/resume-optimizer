# Docker Development Setup

This guide explains how to run the Resume Optimizer using Docker Compose for local development.

## Prerequisites

- Docker and Docker Compose installed
- PostgreSQL running (either locally or in a Kubernetes cluster with port-forward)
- Make sure ports 3000, 8080, 8081 are available

## Quick Start

1. **Set up environment:**
   ```bash
   make setup
   ```
   This creates `.env.docker` from the template.

2. **Edit credentials:**
   ```bash
   # Edit .env.docker with your actual database credentials
   vim .env.docker  # or your preferred editor
   ```

3. **Start everything:**
   ```bash
   make dev
   ```

4. **Access the application:**
   - Frontend: http://localhost:3000
   - Auth API: http://localhost:8080
   - Resume API: http://localhost:8081

## Environment Configuration

Copy `.env.docker.example` to `.env.docker` and update these key settings:

### Database (Required)
```bash
DB_HOST=host.docker.internal    # For external PostgreSQL
DB_PORT=5432
DB_USER=your-db-username        # Update this
DB_PASSWORD=your-db-password    # Update this  
DB_NAME=your-db-name           # Update this
```

### Google OAuth (Optional)
```bash
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### Other Settings
```bash
JWT_SECRET=your-jwt-secret-key-change-this-in-production
REDIRECT_URL=http://localhost:3000/auth/callback
```

## Available Commands

### Development
- `make dev` - Start full development environment
- `make setup` - Set up environment files
- `make health` - Check service health
- `make restart` - Restart all services
- `make rebuild` - Rebuild and restart services

### Control
- `make down` - Stop all services
- `make clean` - Stop services and clean up containers

### Logs
- `make docker-logs` - View all service logs
- `make docker-logs-client` - Frontend logs only
- `make docker-logs-auth` - Auth service logs only  
- `make docker-logs-resume` - Resume processor logs only

### Testing
- `make test-upload` - Test resume upload
- `make test-list` - List uploaded resumes

## Service Architecture

The Docker Compose setup includes:

1. **Frontend (React + Nginx)** - Port 3000
   - Serves the React application
   - Proxies API requests to backend services

2. **Auth Service (Go)** - Port 8080  
   - Handles user authentication
   - Google OAuth integration
   - JWT token management

3. **Resume Processor (Go)** - Port 8081
   - Resume upload and management
   - File processing and storage
   - Resume optimization (future)

4. **External PostgreSQL**
   - Uses your existing PostgreSQL (k8s cluster or local)
   - Connected via host.docker.internal

## Database Setup

This setup assumes you have PostgreSQL running externally (e.g., in your Kubernetes cluster with port-forward active).

### Using Kubernetes PostgreSQL:
```bash
# Make sure your PostgreSQL port-forward is active
kubectl port-forward -n cnpg-system service/postgresql-cluster-rw 5432:5432

# Update .env.docker with these settings:
DB_HOST=host.docker.internal
DB_PORT=5432
DB_USER=app  # or your actual username
DB_PASSWORD=your-actual-password
DB_NAME=app  # or your actual database name
```

### Using Local PostgreSQL:
```bash
# Update .env.docker:
DB_HOST=host.docker.internal  # or localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-local-password  
DB_NAME=resume_optimizer
```

## Troubleshooting

### Services won't start
1. Check if .env.docker exists and has correct values
2. Verify PostgreSQL is accessible
3. Check if ports 3000, 8080, 8081 are available
4. Run `make docker-logs` to see detailed logs

### Database connection issues
1. Verify PostgreSQL port-forward is active: `lsof -i :5432`
2. Test connection manually: `psql "postgres://user:pass@localhost:5432/dbname"`
3. Check DB_* environment variables in `.env.docker`

### Frontend can't reach backend
1. Check that all services are healthy: `make health`
2. Test APIs directly: `curl http://localhost:8080/health`
3. Check nginx proxy configuration in `client/nginx.conf`

### Build issues
1. Clean Docker cache: `make clean`
2. Rebuild everything: `make rebuild`
3. Check Docker logs: `make docker-logs`

## Security Notes

- ⚠️ **Never commit `.env.docker`** - it contains sensitive credentials
- 🔐 **Change JWT_SECRET** in production
- 🔑 **Use strong database passwords**
- 🛡️ **Configure proper Google OAuth redirect URLs**

## Development Workflow

```bash
# Initial setup (one time)
make setup
# Edit .env.docker with your credentials

# Daily development
make dev          # Start everything
make health       # Check it's working
make test-upload  # Test functionality

# When making changes
make rebuild      # Rebuild and restart

# When done
make down         # Stop services
```
