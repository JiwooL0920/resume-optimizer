# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

I want to create a resume optimizing webb app using AI.
User will login, and they will upload their resumes in the settings
then in the main page, they can choose which resume to use
and then they can enter job description url
it will ask AI model (they can choose which model to use, and provide their own API key in the settings), and then it will create a preview page with the optimized version newly created by AI to fit the job description and increase their chance of getting noticed by HR.
user can toggle button to keep the generated resume within 1 page
user can highlight whichever part they want to edit, and can give comments. after adding all comments, they can click 'Edit' button to ask AI again to make changes according to their feedback
once they're happy, the file can be downloaded, and it is auto-saved for the user to retrieve whenever they want
u

This is a resume optimizer project that appears to be in its initial setup phase. The repository currently contains only a basic README.md file with the project name.

## Project Structure

```
resume-optimizer/
├── client/                          # React TypeScript frontend
│   ├── src/                        # React source code
│   ├── public/                     # Static assets
│   ├── package.json                # Node.js dependencies
│   ├── tailwind.config.js          # TailwindCSS configuration
│   └── Dockerfile                  # Frontend container
├── services/
│   ├── auth/                       # Authentication microservice
│   │   ├── main.go                 # Go server entry point
│   │   ├── internal/handlers/      # HTTP handlers
│   │   ├── internal/middleware/    # Authentication middleware
│   │   ├── internal/config/        # Configuration management
│   │   └── Dockerfile              # Service container
│   └── resume-processor/           # Resume processing microservice
│       ├── main.go                 # Go server entry point
│       ├── internal/handlers/      # API handlers
│       └── Dockerfile              # Service container
├── shared/
│   ├── models/                     # Common Go data models
│   │   ├── user.go                 # User and auth models
│   │   └── resume.go               # Resume and optimization models
│   └── database/
│       └── migrations/             # SQL migration files
├── deployment/
│   ├── helm/                       # Kubernetes Helm charts
│   │   ├── client/                 # Frontend deployment charts
│   │   └── services/               # Backend services charts
│   └── scripts/
│       └── setup-db.sh             # Database initialization script
├── docker-compose.yml              # Local development setup
├── Makefile                        # Build and deployment commands
└── .env.example                    # Environment variables template

## Development Setup

- React frontend with TypeScript, TailwindCSS, and Redux
- Go backend services with Gin framework (microservices architecture)
- Kubernetes deployment with Helm charts
- Cloud Native PostgreSQL (CNPG) in cluster namespace 'cnpg-system'
- Google OAuth authentication
- Docker containerization for all services

## Commands

### Build
- `make build-all` - Build all services
- `make build-client` - Build React frontend
- `make build-auth` - Build auth service
- `make build-processor` - Build resume processor service

### Development
- `make dev-client` - Start React development server
- `make dev-backend` - Start backend services with docker-compose
- `npm start --prefix client` - Start React app from root directory
- `go run -C services/auth main.go` - Start auth service directly
- `go run -C services/resume-processor main.go` - Start resume processor directly

### Database
- `make setup-db` - Create database and run migrations on CNPG cluster
- Database: `resume_optimizer` on CNPG cluster in `cnpg-system` namespace

### Deployment
- `make deploy-all` - Deploy all services to Kubernetes
- `make deploy-client` - Deploy client to `resume-optimizer-client` namespace
- `make deploy-services` - Deploy backend services to `resume-optimizer-controller` namespace

## Architecture

### Microservices
- **client/** - React frontend with nginx
- **services/auth/** - Authentication service (Google OAuth, JWT)
- **services/resume-processor/** - Resume optimization and AI processing
- **shared/models/** - Common data models
- **shared/database/** - Database migrations and schemas

### Database Schema
- **users** - User profiles and OAuth data
- **resumes** - User uploaded resumes
- **optimization_sessions** - AI optimization requests and results
- **feedback** - User feedback on optimized resumes
- **user_api_keys** - Encrypted user API keys for AI services

## Notes

This project appears to be in early development with no specific architecture, dependencies, or build system established yet. Future documentation should be updated as the project structure and tooling are defined.