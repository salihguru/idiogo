# Frequently Asked Questions (FAQ)

## General Questions

### What is idiogo?

idiogo is a production-ready Go boilerplate template that combines **idiomatic Go** practices with **Domain-Driven Design (DDD)** principles. It provides a solid foundation for building scalable REST APIs with clean architecture.

### Why use idiogo?

- ✅ **Time-saving**: Skip initial project setup and focus on business logic
- ✅ **Best practices**: Follow established Go and DDD patterns
- ✅ **Production-ready**: Includes Docker, CI/CD, and deployment configurations
- ✅ **Well-documented**: Comprehensive guides and examples
- ✅ **Maintainable**: Clean architecture makes code easy to understand and modify
- ✅ **Testable**: Dependency injection and layer separation simplify testing

### Is idiogo suitable for my project?

idiogo is ideal for:
- REST APIs and microservices
- CRUD applications with business logic
- Projects requiring clean architecture
- Teams wanting DDD structure
- Applications that need to scale

Not recommended for:
- Simple scripts or CLI tools
- Projects with no database
- GraphQL-first applications (without REST)
- Serverless functions (may be too heavy)

## Getting Started

### How do I create a new project from idiogo?

**Option 1: GitHub Template (Recommended)**
1. Click "Use this template" on GitHub
2. Create your repository
3. Clone and customize

**Option 2: Clone directly**
```bash
git clone https://github.com/salihguru/idiogo.git my-project
cd my-project
# Update module name and imports
```

See [Quick Start Guide](QUICKSTART.md) for detailed instructions.

### What are the system requirements?

- Go 1.25.5 or higher
- Docker (for database and containerized deployment)
- PostgreSQL 15+ (can run in Docker)

### How do I change the module name?

1. Update `go.mod`:
```go
module github.com/yourusername/yourproject
```

2. Update all imports:
```bash
find . -type f -name "*.go" -exec sed -i '' 's/github.com\/salihguru\/idiogo/github.com\/yourusername\/yourproject/g' {} +
```

3. Run `go mod tidy`

## Architecture

### What is Domain-Driven Design (DDD)?

DDD is an approach to software development that focuses on:
- Modeling complex business logic
- Using a ubiquitous language
- Organizing code around business domains
- Clear separation between layers

See [Architecture Guide](ARCHITECTURE.md) for details.

### What are the different layers?

1. **Domain Layer**: Business logic and entities
2. **Application Layer**: Use case orchestration
3. **Infrastructure Layer**: Technical implementations (DB, external services)
4. **Interface Layer**: HTTP handlers and adapters

### Why is the repository in the domain layer?

While repository implementations contain infrastructure code, the repository pattern and interfaces belong to the domain. In idiogo, we place the concrete implementation in the domain layer for simplicity, but you can move it to infrastructure if preferred.

### Can I use a different database?

Yes! Replace the GORM PostgreSQL driver with another:

```go
// For MySQL
import "gorm.io/driver/mysql"

// For SQLite
import "gorm.io/driver/sqlite"
```

Update the connection string in `internal/infra/db/pg.go`.

### How do I add middleware?

Create middleware in `internal/rest/middleware/`:

```go
func MyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Middleware logic
        return c.Next()
    }
}
```

Register in `internal/rest/rest.go`:

```go
s.app.Use(s.srv.Recover(), s.srv.I18n(), MyMiddleware())
```

## Development

### How do I add a new domain/module?

1. Create domain directory structure
2. Define entity, repository, service, and handler
3. Register module in `internal/app/serve/modules.go`

See [Quick Start Guide - Add Your First Domain](QUICKSTART.md#add-your-first-domain) for step-by-step instructions.

### How do I add validation rules?

Use struct tags with `go-playground/validator`:

```go
type CreateReq struct {
    Name  string `json:"name" validate:"required,min=3,max=100"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=130"`
}
```

Available validators: [validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10)

### How do I add custom error responses?

Define custom errors in your domain:

```go
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalidStatus = errors.New("invalid status")
)
```

Handle in error handler (`internal/rest/service.go`):

```go
if errors.Is(err, domain.ErrNotFound) {
    code = fiber.StatusNotFound
}
```

### How do I enable hot reload?

Use [Air](https://github.com/cosmtrek/air):

```bash
# Install
go install github.com/cosmtrek/air@latest

# Run
air

# Or with Makefile
make dev
```

### How do I add translations?

1. Add translation key in `assets/locales/en.toml`:
```toml
[errors]
not_found = "Resource not found"
```

2. Add translation in `assets/locales/tr.toml`:
```toml
[errors]
not_found = "Kaynak bulunamadı"
```

3. Use in code:
```go
message := i18n.Localize(ctx, "errors.not_found")
```

## Testing

### How do I write tests?

**Unit tests** for services:
```go
func TestService_Create(t *testing.T) {
    // Setup
    service := NewService(mockRepo)
    
    // Test
    result, err := service.Create(ctx, req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

**Integration tests** for repositories:
```go
func TestRepo_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := NewRepo(db)
    
    err := repo.Save(ctx, entity)
    assert.NoError(t, err)
}
```

### How do I mock dependencies?

Use interfaces and dependency injection:

```go
type MockRepo struct {
    mock.Mock
}

func (m *MockRepo) Save(ctx context.Context, todo *Todo) error {
    args := m.Called(ctx, todo)
    return args.Error(0)
}
```

Or use [testify/mock](https://github.com/stretchr/testify#mock-package).

### How do I run tests?

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/domain/todo/...

# Or use Makefile
make test
make test-coverage
```

## Deployment

### How do I deploy to production?

See [Deployment section in README](../README.md#-deployment) for:
- Docker deployment
- Docker Compose
- Environment configuration
- Production checklist

### Should I use Docker or build binaries?

**Docker** (recommended):
- ✅ Consistent environment
- ✅ Easy deployment
- ✅ Includes all dependencies
- ❌ Slightly larger size

**Binary**:
- ✅ Smaller size
- ✅ Direct execution
- ❌ Requires manual dependency management
- ❌ Environment-specific builds

### How do I manage environment variables?

Create environment-specific config files:

```yaml
# config.prod.yaml
db:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  pass: ${DB_PASSWORD}
```

Or use environment variables directly in code.

### How do I handle database migrations?

idiogo uses GORM auto-migration. To disable:

```yaml
db:
  migrate: false
```

For production, consider:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- Manual SQL migration scripts

### What about secrets management?

For production:
- Use Kubernetes Secrets
- AWS Secrets Manager
- HashiCorp Vault
- Azure Key Vault
- Environment variables (less secure)

**Never commit secrets to Git!**

## Troubleshooting

### Port already in use

Change port in `config.yaml`:
```yaml
rest:
  port: 8080
```

### Database connection failed

1. Check database is running:
```bash
docker ps | grep idiogo-pg
```

2. Verify connection settings in config
3. Check network connectivity
4. Review database logs:
```bash
docker logs idiogo-db
```

### Module import errors

Run:
```bash
go mod tidy
go mod download
```

### Validation not working

Ensure:
1. Struct tags are correct
2. `WithValidation` wrapper is used
3. Validator is initialized in app

### CORS errors

Add CORS middleware:
```go
import "github.com/gofiber/fiber/v2/middleware/cors"

app.Use(cors.New())
```

## Performance

### How do I improve performance?

1. **Enable database connection pooling** (already configured)
2. **Add caching** (Redis, in-memory cache)
3. **Use database indexes** on frequently queried fields
4. **Optimize queries** (select specific fields, use joins)
5. **Add pagination** to list endpoints
6. **Enable compression** in Fiber
7. **Profile your code** with pprof

### Should I use caching?

Consider caching when:
- Data is read frequently
- Data changes infrequently
- Database queries are expensive
- Response time is critical

## Contributing

### How can I contribute?

See [CONTRIBUTING.md](../CONTRIBUTING.md) for:
- Code of conduct
- Development setup
- Coding standards
- Pull request process

### I found a bug!

Please [open an issue](https://github.com/salihguru/idiogo/issues) with:
- Description of the bug
- Steps to reproduce
- Expected vs actual behavior
- Environment details

### Can I suggest a feature?

Yes! [Open a feature request](https://github.com/salihguru/idiogo/issues) describing:
- What you want to achieve
- Why it's useful
- Possible implementation

## License

### What license does idiogo use?

MIT License - you can use it freely in personal and commercial projects.

### Can I use idiogo commercially?

Yes! The MIT license allows commercial use.

### Do I need to credit idiogo?

Not required, but appreciated! You can:
- Keep the LICENSE file
- Mention idiogo in your README
- Star the repository on GitHub

---

**Still have questions?** 

- 📖 Check the [Documentation](.)
- 💬 Open a [Discussion](https://github.com/salihguru/idiogo/discussions)
- 🐛 Report an [Issue](https://github.com/salihguru/idiogo/issues)
