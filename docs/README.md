# Documentation

Welcome to the idiogo documentation! This guide will help you understand, use, and contribute to idiogo.

## 📚 Table of Contents

### Getting Started
- **[Quick Start Guide](QUICKSTART.md)** - Get up and running in minutes
- **[FAQ](FAQ.md)** - Frequently asked questions and troubleshooting

### Core Documentation
- **[Architecture Guide](ARCHITECTURE.md)** - Deep dive into the architecture and design patterns
- **[API Reference](API.md)** - Complete API endpoint documentation

### Development
- **[Contributing Guidelines](../CONTRIBUTING.md)** - How to contribute to idiogo
- **[Code of Conduct](../CONTRIBUTING.md#code-of-conduct)** - Community guidelines

### Project Information
- **[README](../README.md)** - Project overview and features
- **[License](../LICENSE)** - MIT License details

## 🎯 Quick Navigation

### I want to...

#### Learn about idiogo
→ Start with the [README](../README.md) for an overview  
→ Read the [Architecture Guide](ARCHITECTURE.md) to understand the structure

#### Start a new project
→ Follow the [Quick Start Guide](QUICKSTART.md)  
→ Check [FAQ](FAQ.md#how-do-i-create-a-new-project-from-idiogo) for setup help

#### Add features to my project
→ See [Quick Start - Add Your First Domain](QUICKSTART.md#add-your-first-domain)  
→ Review the example Todo domain in `internal/domain/todo/`

#### Understand the API
→ Read the [API Reference](API.md)  
→ Try the example requests in [Quick Start - Test the API](QUICKSTART.md#test-the-api)

#### Deploy my application
→ Check [Deployment](../README.md#-deployment) in README  
→ Review Docker configurations in `deployments/`

#### Contribute
→ Read [Contributing Guidelines](../CONTRIBUTING.md)  
→ Check open [Issues](https://github.com/salihguru/idiogo/issues)

## 📖 Key Concepts

### Domain-Driven Design (DDD)

idiogo implements DDD principles:

- **Entities**: Objects with identity (e.g., Todo, User)
- **Value Objects**: Immutable objects without identity (e.g., Status)
- **Aggregates**: Cluster of entities with defined boundaries
- **Repositories**: Persistence abstraction
- **Services**: Business logic that doesn't fit in entities

Learn more in the [Architecture Guide](ARCHITECTURE.md#domain-driven-design).

### Clean Architecture

The project follows clean architecture with clear layers:

```
Domain (Business Logic)
    ↓
Application (Use Cases)
    ↓
Infrastructure (Technical Details)
    ↓
Interfaces (HTTP, gRPC, etc.)
```

Dependencies point inward - inner layers know nothing about outer layers.

### Module Pattern

Each domain is organized as a self-contained module:

```go
Module {
    Repo    // Data access
    Service // Business logic
    Handler // HTTP interface
}
```

This pattern promotes:
- **Encapsulation**: Each module is self-contained
- **Reusability**: Modules can be easily reused or extracted
- **Testability**: Easy to mock dependencies

## 🏗️ Project Structure

```
idiogo/
├── cmd/                  # Application entry points
│   ├── serve/           # REST API server
│   └── cron/            # Background jobs
├── internal/            # Private application code
│   ├── app/            # Application initialization
│   ├── domain/         # Business logic (DDD domains)
│   ├── infra/          # Infrastructure (DB, external services)
│   ├── rest/           # REST API implementation
│   ├── port/           # Port interfaces
│   └── config/         # Configuration
├── pkg/                 # Public shared packages
│   ├── entity/         # Base entities
│   ├── validation/     # Validation utilities
│   ├── i18np/          # Internationalization
│   └── ...             # Other utilities
├── assets/              # Static assets (locales, etc.)
├── deployments/         # Deployment configs
├── docs/                # Documentation (you are here!)
└── vendor/              # Vendored dependencies
```

## 🔧 Configuration

idiogo uses YAML configuration files:

```yaml
rest:
  host: 0.0.0.0
  port: 4041

db:
  host: localhost
  port: "5432"
  user: idiogo
  pass: idiogo
  name: idiogo
  ssl_mode: disable
  migrate: true
  debug: false

i18n:
  locales: [en, tr]
  default: en
  dir: "./assets/locales"
```

See example in `deployments/config.yml`.

## 🧪 Testing

idiogo supports multiple testing levels:

**Unit Tests**: Test business logic in isolation
```bash
go test ./internal/domain/todo/...
```

**Integration Tests**: Test with real dependencies
```bash
go test -tags=integration ./...
```

**Full Test Suite**: Run all tests with coverage
```bash
make test-coverage
```

See [Architecture Guide - Testing Strategy](ARCHITECTURE.md#testing-strategy).

## 🚀 Deployment

idiogo can be deployed in multiple ways:

### Docker Compose (Development/Simple Production)
```bash
docker compose -f deployments/compose.yml up -d
```

### Docker (Production)
```bash
docker build -f cmd/serve/Dockerfile -t idiogo:latest .
docker run -p 4041:4041 idiogo:latest
```

### Binary (Native)
```bash
make build
./bin/serve
```

See [README - Deployment](../README.md#-deployment) for details.

## 📦 Dependencies

Key dependencies:

- **[Fiber](https://gofiber.io/)**: High-performance HTTP framework
- **[GORM](https://gorm.io/)**: ORM for database operations
- **[validator](https://github.com/go-playground/validator)**: Request validation
- **[go-i18n](https://github.com/nicksnyder/go-i18n)**: Internationalization
- **[uuid](https://github.com/google/uuid)**: UUID generation

## 🛠️ Development Tools

Recommended tools:

- **[Air](https://github.com/cosmtrek/air)**: Hot reload for Go apps
- **[golangci-lint](https://golangci-lint.run/)**: Linter aggregator
- **[gosec](https://github.com/securego/gosec)**: Security scanner
- **[Docker](https://www.docker.com/)**: Containerization
- **[Make](https://www.gnu.org/software/make/)**: Build automation

Install tools:
```bash
make install-tools
```

## 📝 Examples

### Creating a Todo

```bash
curl -X POST http://localhost:4041/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
```

### Listing Todos

```bash
curl "http://localhost:4041/todos?page=1&limit=10&status=pending"
```

More examples in the [API Reference](API.md).

## 🤝 Contributing

We welcome contributions! Here's how:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write/update tests
5. Submit a pull request

See [Contributing Guidelines](../CONTRIBUTING.md) for details.

## 📄 License

idiogo is licensed under the MIT License. See [LICENSE](../LICENSE) for details.

## 🆘 Getting Help

- 📖 **Documentation**: You're reading it!
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/salihguru/idiogo/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/salihguru/idiogo/discussions)
- ⭐ **Star us**: [GitHub Repository](https://github.com/salihguru/idiogo)

## 🗺️ Roadmap

Planned features (community input welcome!):

- [ ] GraphQL support
- [ ] gRPC support
- [ ] Authentication/Authorization module
- [ ] Rate limiting middleware
- [ ] Observability (metrics, tracing)
- [ ] More database adapters
- [ ] CLI code generator
- [ ] Additional example domains

See [Issues](https://github.com/salihguru/idiogo/issues) for planned work.

## 📊 Diagrams

### Request Flow

```
HTTP Request
    ↓
Middleware (I18n, Recovery, etc.)
    ↓
Router
    ↓
Handler Wrapper (Parse, Validate)
    ↓
Domain Service
    ↓
Repository
    ↓
Database
```

### Module Structure

```
Module
├── Entity (Domain Model)
├── Repository (Data Access)
├── Service (Business Logic)
└── Handler (HTTP Interface)
```

More diagrams in [Architecture Guide](ARCHITECTURE.md).

## 🎓 Learning Resources

### Go Resources
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### DDD Resources
- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)
- [DDD Reference](https://www.domainlanguage.com/ddd/reference/)

### Architecture Resources
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

---

**Happy coding with idiogo!** 🚀

If you find idiogo useful, please ⭐ [star it on GitHub](https://github.com/salihguru/idiogo)!
