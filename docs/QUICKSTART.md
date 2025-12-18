# Quick Start Guide

Get up and running with idiogo in minutes!

## Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose (for database)
- Git

## Installation

### Option 1: Use as GitHub Template

1. Click "Use this template" on GitHub
2. Clone your new repository
3. Update module name in `go.mod`
4. Update imports throughout the project

### Option 2: Clone and Customize

```bash
# Clone the repository
git clone https://github.com/salihguru/idiogo.git my-project
cd my-project

# Update go.mod
# Change "github.com/salihguru/idiogo" to your module path
vim go.mod

# Update imports in all Go files
find . -type f -name "*.go" -exec sed -i '' 's/github.com\/salihguru\/idiogo/your-module-path/g' {} +

# Initialize as new git repository
rm -rf .git
git init
git add .
git commit -m "Initial commit from idiogo template"
```

## Configuration

1. Create configuration file:

```bash
cp deployments/config.yml config.yaml
```

2. Edit `config.yaml`:

```yaml
rest:
  host: 0.0.0.0
  port: 4041

db:
  host: localhost    # or idiogo-pg for Docker
  port: "5432"
  user: idiogo
  pass: idiogo
  name: idiogo
  ssl_mode: disable
  migrate: true      # Auto-run migrations
  debug: true        # SQL query logging

i18n:
  locales: [en, tr]
  default: en
  dir: "./assets/locales"
```

## Start the Database

```bash
docker compose -f deployments/compose.yml up -d idiogo-pg
```

Verify database is running:
```bash
docker ps
```

## Run the Application

### Development Mode

```bash
go run cmd/serve/main.go
```

You should see:
```
idiogo api is running on 0.0.0.0:4041
```

### Using Docker Compose

```bash
docker compose -f deployments/compose.yml up -d
```

This starts both database and application.

## Test the API

### Create a Todo

```bash
curl -X POST http://localhost:4041/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Todo",
    "description": "Testing idiogo API"
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "My First Todo",
  "description": "Testing idiogo API",
  "status": "pending",
  "created_at": "2025-12-18T10:00:00Z"
}
```

### List Todos

```bash
curl http://localhost:4041/todos
```

### Get a Todo

```bash
curl http://localhost:4041/todos/{id}
```

### Update a Todo

```bash
curl -X PATCH http://localhost:4041/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

### Delete a Todo

```bash
curl -X DELETE http://localhost:4041/todos/{id}
```

## Next Steps

### Add Your First Domain

1. **Create domain directory**:
```bash
mkdir -p internal/domain/product
```

2. **Create entity** (`internal/domain/product/product.go`):
```go
package product

import "github.com/your-module/idiogo/pkg/entity"

type Product struct {
    entity.Base
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

3. **Create repository** (`internal/domain/product/repo.go`):
```go
package product

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

func (r *Repo) Save(ctx context.Context, p *Product) error {
    return r.db.WithContext(ctx).Save(p).Error
}

func (r *Repo) View(ctx context.Context, id uuid.UUID) (*Product, error) {
    var p Product
    err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
    return &p, err
}
```

4. **Create service** (`internal/domain/product/service.go`):
```go
package product

import "context"

type Service struct {
    repo *Repo
}

func NewService(repo *Repo) *Service {
    return &Service{repo: repo}
}

type CreateReq struct {
    Name  string  `json:"name" validate:"required"`
    Price float64 `json:"price" validate:"required,gt=0"`
}

func (s *Service) Create(ctx context.Context, req CreateReq) (*Product, error) {
    p := &Product{
        Name:  req.Name,
        Price: req.Price,
    }
    if err := s.repo.Save(ctx, p); err != nil {
        return nil, err
    }
    return p, nil
}
```

5. **Create handler** (`internal/domain/product/handler.go`):
```go
package product

import (
    "github.com/gofiber/fiber/v2"
    "github.com/your-module/idiogo/internal/port"
    "github.com/your-module/idiogo/internal/rest"
)

type Handler struct {
    srv Service
}

func NewHandler(srv Service) *Handler {
    return &Handler{srv: srv}
}

func (h *Handler) RegisterRoutes(srv port.RestService, router fiber.Router) {
    group := router.Group("/products")
    
    group.Post("/",
        srv.Timeout(rest.Handle(rest.WithBody(rest.WithValidation(srv.ValidateStruct(),
            rest.CreateResponds(h.srv.Create))))))
}
```

6. **Register module** in `internal/app/serve/modules.go`:
```go
import "github.com/your-module/idiogo/internal/domain/product"

func newModules(d *Depends) Modules {
    // ... existing todo module
    
    // Add product module
    productModule := rest.Module[*product.Repo, *product.Service]{
        Repo:    product.NewRepo(d.DB),
        Service: nil,
        Router:  nil,
    }
    productModule.Service = product.NewService(productModule.Repo)
    productModule.Router = product.NewHandler(*productModule.Service)
    
    return Modules{
        Todo:    todoModule.Router,
        Product: productModule.Router,
    }
}
```

7. **Update Modules struct** in `internal/app/serve/modules.go`:
```go
type Modules struct {
    Todo    rest.Router
    Product rest.Router
}

func (m Modules) Routers() []rest.Router {
    return []rest.Router{
        m.Todo,
        m.Product,
    }
}
```

8. **Restart and test**:
```bash
# Stop the server (Ctrl+C)
# Start again
go run cmd/serve/main.go

# Test new endpoint
curl -X POST http://localhost:4041/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "price": 29.99
  }'
```

## Development Tips

### Enable Hot Reload

Use [Air](https://github.com/cosmtrek/air):

```bash
# Install
go install github.com/cosmtrek/air@latest

# Run
air
```

### View Logs

```bash
# Docker Compose
docker compose -f deployments/compose.yml logs -f

# Application logs
tail -f /path/to/log/file
```

### Database Management

```bash
# Connect to database
docker exec -it idiogo-db psql -U idiogo -d idiogo

# View tables
\dt

# Query todos
SELECT * FROM todos;
```

### Run Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/domain/todo/...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

## Troubleshooting

### Port Already in Use

Change port in `config.yaml`:
```yaml
rest:
  port: 8080  # Use different port
```

### Database Connection Failed

Check database is running:
```bash
docker ps | grep idiogo-pg
```

Verify connection settings in `config.yaml`:
```yaml
db:
  host: localhost  # or idiogo-pg for Docker network
  port: "5432"
```

### Module Import Errors

Ensure you've updated all imports after changing module name:
```bash
go mod tidy
```

### Migration Errors

Drop and recreate database:
```bash
docker compose -f deployments/compose.yml down -v
docker compose -f deployments/compose.yml up -d
```

## What's Next?

- Read [Architecture Documentation](docs/ARCHITECTURE.md)
- Check [Contributing Guidelines](CONTRIBUTING.md)
- Explore the example Todo domain
- Build your application!

## Getting Help

- 📚 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/salihguru/idiogo/issues)
- 💬 [Discussions](https://github.com/salihguru/idiogo/discussions)

Happy coding! 🚀
