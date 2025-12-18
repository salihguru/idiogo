# idiogo

> **Idiomatic Go + Domain-Driven Design** - A production-ready Go template combining idiomatic Go practices with DDD architecture patterns.

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Development](#development)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

**idiogo** is a modern, production-ready Go boilerplate that combines **idiomatic Go practices** with **Domain-Driven Design (DDD)** principles. It provides a solid foundation for building scalable, maintainable REST APIs with clean architecture.

This template is designed to help you kickstart Go projects with best practices baked in, including:

- Clean DDD architecture
- Type-safe request/response handling
- Built-in validation and i18n support
- Database migrations with GORM
- Docker-ready deployment
- Middleware patterns

## Useful Links

- [`golang code review comments`](https://go.dev/wiki/CodeReviewComments)
- [`effective go docs`](https://go.dev/doc/effective_go)
- [`golang styleguide by Google`](https://google.github.io/styleguide/go/decisions)

## ✨ Features

### Core Features

- 🏗️ **Domain-Driven Design**: Clean separation of concerns with domain, application, and infrastructure layers
- 🚀 **Fiber Framework**: High-performance HTTP framework built on Fasthttp
- 🗃️ **GORM ORM**: Powerful ORM with PostgreSQL support and auto-migrations
- ✅ **Validation**: Comprehensive request validation with go-playground/validator
- 🌍 **Internationalization**: Multi-language support with go-i18n
- 🔄 **Type-Safe Handlers**: Generic handler patterns for type-safe request/response handling
- 📦 **Dependency Injection**: Clean dependency management pattern
- 🐳 **Docker Support**: Production-ready Docker and Docker Compose configurations

### Developer Experience

- 🔧 **Hot Reload Ready**: Easy integration with development tools
- 📝 **Structured Logging**: Clean, readable log output
- 🎯 **Graceful Shutdown**: Proper resource cleanup on termination
- 🔐 **Security**: Built-in security headers and proxy detection
- 📊 **Database Migrations**: Automatic schema management with GORM

## 🏛️ Architecture

idiogo follows a **Layered DDD Architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────┐
│            cmd/ (Entry Points)              │
│  • serve: REST API server                   │
│  • cron: Background jobs                    │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│        internal/app (Application)           │
│  • Application initialization               │
│  • Module composition                       │
│  • Dependency wiring                        │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│       internal/domain (Business Logic)      │
│  • Entities & Aggregates                    │
│  • Services (Business Logic)                │
│  • Repository Interfaces                    │
│  • Handlers (HTTP)                          │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│        internal/infra (Infrastructure)      │
│  • Database connections                     │
│  • External service integrations            │
│  • Repository implementations               │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          pkg/ (Shared Packages)             │
│  • Reusable utilities                       │
│  • Common types & helpers                   │
└─────────────────────────────────────────────┘
```

### Key Principles

1. **Domain Layer** (`internal/domain/`): Contains business logic, entities, and domain services
2. **Application Layer** (`internal/app/`): Orchestrates domain logic and handles use cases
3. **Infrastructure Layer** (`internal/infra/`): Technical implementations (DB, external APIs)
4. **Interface Layer** (`internal/rest/`, `internal/port/`): HTTP handlers and protocol adapters
5. **Shared Kernel** (`pkg/`): Reusable utilities and common abstractions

## 📁 Project Structure

```
idiogo/
├── cmd/                          # Application entry points
│   ├── serve/                    # REST API server
│   │   ├── main.go
│   │   └── Dockerfile
│   └── cron/                     # Background jobs
│       └── main.go
├── internal/                     # Private application code
│   ├── app/                      # Application layer
│   │   ├── serve/               # Server initialization
│   │   └── cron/                # Cron initialization
│   ├── domain/                   # Business logic layer
│   │   └── todo/                # Example domain (Todo)
│   │       ├── todo.go          # Entity
│   │       ├── service.go       # Business logic
│   │       ├── repo.go          # Repository implementation
│   │       ├── handler.go       # HTTP handlers
│   │       └── filters.go       # Query filters
│   ├── infra/                    # Infrastructure layer
│   │   └── db/                  # Database
│   │       ├── pg.go            # PostgreSQL connection
│   │       └── migration/       # DB migrations
│   ├── rest/                     # REST API implementation
│   │   ├── rest.go              # Server setup
│   │   ├── service.go           # REST service utilities
│   │   ├── handler.go           # Generic handlers
│   │   └── middleware/          # HTTP middlewares
│   ├── port/                     # Ports (interfaces)
│   │   └── rest.go              # REST port interface
│   └── config/                   # Configuration
│       ├── config.go            # Config structures
│       └── bind.go              # Config loader
├── pkg/                          # Public shared packages
│   ├── entity/                  # Base entities
│   ├── i18np/                   # Internationalization
│   ├── validation/              # Validation utilities
│   ├── query/                   # Query builders
│   ├── state/                   # State management
│   ├── list/                    # List/pagination
│   └── server/                  # Server utilities
├── assets/                       # Static assets
│   └── locales/                 # Translation files
│       ├── en.toml
│       └── tr.toml
├── deployments/                  # Deployment configurations
│   ├── compose.yml              # Docker Compose
│   └── config.yml               # Application config
├── docs/                         # Documentation
├── tmp/                          # Temporary files (gitignored)
├── go.mod                        # Go modules
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.25.5 or higher
- **Docker**: Latest version (for containerized deployment)
- **PostgreSQL**: 15+ (or use Docker Compose)

### Installation

1. **Clone the repository** (or use as a GitHub template):

```bash
git clone https://github.com/salihguru/idiogo.git my-project
cd my-project
```

2. **Install dependencies**:

```bash
go mod download
```

3. **Update module name**:

Replace `github.com/salihguru/idiogo` with your module path in:

- `go.mod`
- All import statements

```bash
# Use a script or IDE refactoring tools
find . -type f -name "*.go" -exec sed -i '' 's/github.com\/salihguru\/idiogo/your-module-path/g' {} +
```

4. **Configure the application**:

Copy the example configuration:

```bash
cp deployments/config.yml config.yaml
```

Edit `config.yaml` with your settings.

5. **Start the database**:

```bash
docker compose -f deployments/compose.yml up -d idiogo-pg
```

6. **Run migrations** (automatic on startup with `migrate: true` in config)

7. **Start the server**:

```bash
go run cmd/serve/main.go
```

The API will be available at `http://localhost:4041`

## ⚙️ Configuration

Configuration is managed through YAML files. See [deployments/config.yml](deployments/config.yml) for an example.

### Configuration Structure

```yaml
rest:
  host: 0.0.0.0           # Server host
  port: 4041               # Server port

db:
  host: localhost          # Database host
  port: "5432"             # Database port
  user: idiogo             # Database user
  pass: idiogo             # Database password
  name: idiogo             # Database name
  ssl_mode: disable        # SSL mode (disable, require, verify-ca, verify-full)
  migrate: true            # Auto-run migrations on startup
  debug: false             # Enable SQL query logging

i18n:
  locales:                 # Supported locales
    - en
    - tr
  default: en              # Default locale
  dir: "./assets/locales"  # Locale files directory
```

### Environment Variables

You can override configuration values with environment variables (when implemented):

```bash
export DB_HOST=localhost
export DB_PORT=5432
export REST_PORT=8080
```

## 🛠️ Development

### Adding a New Domain

1. **Create domain directory**:

```bash
mkdir -p internal/domain/yourdomein
```

2. **Define the entity** (`internal/domain/yourdomain/entity.go`):

```go
package yourdomain

import "github.com/salihguru/idiogo/pkg/entity"

type YourEntity struct {
    entity.Base
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

3. **Create the repository** (`internal/domain/yourdomain/repo.go`):

```go
package yourdomain

import (
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Repo struct {
    db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
    return &Repo{db: db}
}

func (r *Repo) Save(ctx context.Context, entity *YourEntity) error {
    return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Repo) View(ctx context.Context, id uuid.UUID) (*YourEntity, error) {
    var entity YourEntity
    err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
    return &entity, err
}
```

4. **Implement the service** (`internal/domain/yourdomain/service.go`):

```go
package yourdomain

import "context"

type Service struct {
    repo *Repo
}

func NewService(repo *Repo) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateReq) (*YourEntity, error) {
    entity := &YourEntity{
        Name:        req.Name,
        Description: req.Description,
    }
    if err := s.repo.Save(ctx, entity); err != nil {
        return nil, err
    }
    return entity, nil
}
```

5. **Create HTTP handlers** (`internal/domain/yourdomain/handler.go`):

```go
package yourdomain

import (
    "github.com/gofiber/fiber/v2"
    "github.com/salihguru/idiogo/internal/port"
    "github.com/salihguru/idiogo/internal/rest"
)

type Handler struct {
    srv Service
}

func NewHandler(srv Service) *Handler {
    return &Handler{srv: srv}
}

func (h *Handler) RegisterRoutes(srv port.RestService, router fiber.Router) {
    group := router.Group("/yourdomains")
    
    group.Post("/",
        srv.Timeout(rest.Handle(rest.WithBody(rest.WithValidation(srv.ValidateStruct(),
            rest.CreateResponds(h.srv.Create))))))
}
```

6. **Register the module** in `internal/app/serve/modules.go`:

```go
yourDomainRepo := yourdomain.NewRepo(d.DB)
yourDomainService := yourdomain.NewService(yourDomainRepo)
yourdomainModule := rest.Module[*yourdomain.Repo, *yourdomain.Service]{
    Repo:    yourDomainRepo,
    Service: yourDomainService,
    Router:  yourdomain.NewHandler(*yourdomainModule.Service),
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/domain/todo/...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Security scan
gosec ./...
```

## 📚 API Documentation

### Example: Todo API

#### Create Todo

```http
POST /todos
Content-Type: application/json

{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending",
  "created_at": "2025-12-18T10:00:00Z"
}
```

#### List Todos

```http
GET /todos?page=1&limit=10&status=pending
```

#### Get Todo

```http
GET /todos/{id}
```

#### Update Todo

```http
PATCH /todos/{id}
Content-Type: application/json

{
  "title": "Buy groceries and fruits",
  "status": "completed"
}
```

#### Delete Todo

```http
DELETE /todos/{id}
```

### Response Format

All API responses follow a consistent format:

**Success Response:**

```json
{
  "data": { /* response data */ }
}
```

**Error Response:**

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "title",
        "message": "title is required"
      }
    ]
  }
}
```

## 🧪 Testing

The project includes examples for:

- Unit tests
- Integration tests
- Repository tests

Example test structure:

```go
func TestService_Create(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    repo := NewRepo(db)
    service := NewService(repo)
    
    // Test
    result, err := service.Create(context.Background(), CreateReq{
        Title: "Test Todo",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Test Todo", result.Title)
}
```

## 🐳 Deployment

### Docker Compose (Recommended for Development)

```bash
# Start all services
docker compose -f deployments/compose.yml up -d

# View logs
docker compose -f deployments/compose.yml logs -f

# Stop services
docker compose -f deployments/compose.yml down
```

### Docker Build

```bash
# Build image
docker build -f cmd/serve/Dockerfile -t idiogo:latest .

# Run container
docker run -p 4041:4041 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  idiogo:latest
```

### Production Deployment

For production deployment:

1. Set environment-specific configuration
2. Use secrets management for sensitive data
3. Enable SSL/TLS for database connections
4. Configure reverse proxy (nginx, traefik)
5. Set up monitoring and logging
6. Enable health checks

Example production config:

```yaml
rest:
  host: 0.0.0.0
  port: 4041

db:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  pass: ${DB_PASSWORD}
  name: ${DB_NAME}
  ssl_mode: require
  migrate: false
  debug: false

i18n:
  locales: [en, tr]
  default: en
  dir: "/app/assets/locales"
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [go-playground/validator](https://github.com/go-playground/validator) - Validation
- [go-i18n](https://github.com/nicksnyder/go-i18n) - Internationalization
- [rescode](https://github.com/restayway/rescode) - Response code management

## 📞 Support

- 📧 Email: <contact@salih.guru>
- 🐛 Issues: [GitHub Issues](https://github.com/salihguru/idiogo/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/salihguru/idiogo/discussions)

---

**Built with ❤️ using idiomatic Go and DDD principles**
