# Contributing to idiogo

First off, thank you for considering contributing to idiogo! It's people like you that make idiogo such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

**Bug Report Template:**
```markdown
**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Execute '...'
3. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Environment:**
 - OS: [e.g., macOS, Linux, Windows]
 - Go version: [e.g., 1.25.5]
 - idiogo version: [e.g., v1.0.0]

**Additional context**
Add any other context about the problem here.
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful** to most idiogo users
- **List some examples** of how this enhancement would be used

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the coding standards** outlined below
3. **Write clear commit messages** following conventional commits
4. **Add tests** if you're adding functionality
5. **Update documentation** if needed
6. **Ensure tests pass** before submitting

#### Pull Request Process

1. Update the README.md with details of changes if applicable
2. Update the documentation in the `docs/` folder if needed
3. The PR will be merged once you have the sign-off of at least one maintainer

## Development Setup

### Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose
- Git

### Setup Steps

1. Fork and clone the repository:
```bash
git clone https://github.com/YOUR-USERNAME/idiogo.git
cd idiogo
```

2. Install dependencies:
```bash
go mod download
```

3. Start the development database:
```bash
docker compose -f deployments/compose.yml up -d idiogo-pg
```

4. Run the tests:
```bash
go test ./...
```

5. Run the application:
```bash
go run cmd/serve/main.go
```

## Coding Standards

### Go Style Guide

- Follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- Use `gofmt` to format your code
- Run `go vet` to catch common errors
- Use meaningful variable and function names
- Keep functions small and focused

### Code Formatting

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run security scan
gosec ./...
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `todo`, `user`)
- **Files**: snake_case (e.g., `service.go`, `handler.go`)
- **Interfaces**: describe behavior (e.g., `Service`, `Repository`)
- **Structs**: PascalCase (e.g., `TodoService`, `UserHandler`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported

### Project Structure

When adding new domains:
```
internal/domain/newdomain/
├── entity.go      # Domain entities
├── service.go     # Business logic
├── repo.go        # Repository implementation
├── handler.go     # HTTP handlers
└── filters.go     # Query filters (if needed)
```

### Testing

- Write unit tests for business logic
- Write integration tests for repositories
- Use table-driven tests where applicable
- Mock external dependencies
- Aim for >80% code coverage

Example test:
```go
func TestService_Create(t *testing.T) {
    tests := []struct {
        name    string
        req     CreateReq
        want    *Todo
        wantErr bool
    }{
        {
            name: "valid request",
            req: CreateReq{
                Title:       "Test Todo",
                Description: "Test Description",
            },
            wantErr: false,
        },
        {
            name: "empty title",
            req: CreateReq{
                Title: "",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(todo): add status filter to list endpoint
fix(auth): correct JWT token expiration handling
docs: update installation instructions
refactor(repo): simplify query builder
```

### Documentation

- Add godoc comments to exported functions, types, and packages
- Update README.md for user-facing changes
- Update architecture docs for structural changes
- Include code examples where helpful

Example documentation:
```go
// Service handles business logic for todo operations.
// It encapsulates all todo-related use cases and coordinates
// with the repository layer for data persistence.
type Service struct {
    repo *Repo
}

// Create creates a new todo item with the given parameters.
// It validates the request, creates the entity, and persists it.
//
// Example:
//   req := CreateReq{Title: "Buy milk", Description: "2% milk"}
//   todo, err := service.Create(ctx, req)
//
// Returns an error if validation fails or database operation fails.
func (s *Service) Create(ctx context.Context, req CreateReq) (*Todo, error) {
    // Implementation
}
```

## Architecture Guidelines

### Domain-Driven Design

- Keep business logic in the domain layer
- Domain entities should be rich models, not anemic
- Use value objects for concepts without identity
- Repository interfaces belong to the domain
- Services orchestrate domain operations

### Dependency Rules

- Domain layer has no external dependencies
- Application layer depends on domain
- Infrastructure depends on domain interfaces
- No circular dependencies

### Error Handling

- Return errors, don't panic
- Wrap errors with context
- Use custom error types for domain errors
- Handle errors at appropriate levels

```go
// Good
if err := repo.Save(ctx, todo); err != nil {
    return nil, fmt.Errorf("failed to save todo: %w", err)
}

// Bad
if err := repo.Save(ctx, todo); err != nil {
    panic(err)
}
```

### Validation

- Validate at boundaries (HTTP handlers)
- Use struct tags for declarative validation
- Return structured validation errors
- Separate validation from business logic

## Questions?

Feel free to:
- Open an issue for discussion
- Ask in pull request comments
- Reach out to maintainers

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- GitHub contributors page

Thank you for contributing to idiogo! 🎉
