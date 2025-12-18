# Architecture Guide

## Table of Contents

- [Overview](#overview)
- [Architectural Principles](#architectural-principles)
- [Layer Architecture](#layer-architecture)
- [Domain-Driven Design](#domain-driven-design)
- [Module Pattern](#module-pattern)
- [Request Flow](#request-flow)
- [Error Handling](#error-handling)
- [Testing Strategy](#testing-strategy)
- [Best Practices](#best-practices)

## Overview

idiogo implements a **Clean Architecture** approach combined with **Domain-Driven Design (DDD)** principles. The architecture is designed to be:

- **Maintainable**: Clear separation of concerns
- **Testable**: Dependencies point inward, making testing easier
- **Flexible**: Easy to swap implementations
- **Scalable**: Modular design supports growth

## Architectural Principles

### 1. Dependency Rule

Dependencies can only point inward. Inner layers know nothing about outer layers.

```
┌────────────────────────────────────┐
│     Frameworks & Drivers           │  ← Outermost
│  (HTTP, DB, External Services)     │
├────────────────────────────────────┤
│     Interface Adapters             │
│  (REST Handlers, Presenters)       │
├────────────────────────────────────┤
│     Application Business Rules     │
│  (Use Cases, Services)             │
├────────────────────────────────────┤
│     Enterprise Business Rules      │  ← Innermost
│  (Entities, Domain Logic)          │
└────────────────────────────────────┘
```

### 2. Single Responsibility Principle

Each layer, module, and component has one reason to change.

### 3. Interface Segregation

Clients should not depend on interfaces they don't use. Define specific port interfaces.

### 4. Dependency Inversion

High-level modules should not depend on low-level modules. Both should depend on abstractions.

## Layer Architecture

### Domain Layer (`internal/domain/`)

**Purpose**: Contains the core business logic and rules.

**Characteristics**:
- No external dependencies (except standard library)
- Pure business logic
- Framework-agnostic
- Highly testable

**Components**:

```go
domain/
└── todo/
    ├── todo.go          # Entity (Aggregate Root)
    ├── service.go       # Domain Service (Business Logic)
    ├── repo.go          # Repository Implementation
    ├── handler.go       # HTTP Handler (Adapter)
    └── filters.go       # Query Objects
```

**Example Entity**:
```go
type Todo struct {
    entity.Base
    Title       string
    Description string
    Status      Status
}

// Business rules are encoded in the entity
func (t *Todo) Complete() error {
    if t.Status == StatusArchived {
        return ErrCannotCompleteArchivedTodo
    }
    t.Status = StatusCompleted
    return nil
}
```

### Application Layer (`internal/app/`)

**Purpose**: Orchestrates the flow of data and coordinates domain operations.

**Characteristics**:
- Application-specific business rules
- Depends on domain layer
- Independent of frameworks
- Coordinates use cases

**Components**:

```go
app/
├── serve/
│   ├── serve.go         # Application initialization
│   ├── modules.go       # Module composition
│   └── depends.go       # Dependency management
└── cron/
    └── cron.go          # Background job initialization
```

**Responsibilities**:
- Initialize dependencies
- Wire up modules
- Configure services
- Manage application lifecycle

### Infrastructure Layer (`internal/infra/`)

**Purpose**: Provides technical capabilities that support higher layers.

**Characteristics**:
- Framework and library integration
- External service connections
- Technical implementations
- Replaceable components

**Components**:

```go
infra/
└── db/
    ├── pg.go            # PostgreSQL connection
    └── migration/       # Database migrations
```

### Interface Layer (`internal/rest/`, `internal/port/`)

**Purpose**: Adapts external communications to internal use cases.

**Characteristics**:
- Protocol-specific code
- Request/response transformation
- Middleware and filters
- Port interfaces

**Components**:

```go
rest/
├── rest.go              # Server configuration
├── service.go           # REST utilities
├── handler.go           # Generic handlers
├── base.go              # Base response types
└── middleware/          # HTTP middlewares

port/
└── rest.go              # Port interface definitions
```

## Domain-Driven Design

### Entities

Objects with identity that persist over time.

```go
type Todo struct {
    entity.Base  // ID, timestamps, soft delete
    Title       string
    Description string
    Status      Status
}
```

**Characteristics**:
- Have unique identity (UUID)
- Mutable
- Contain business logic
- Track lifecycle (created, updated, deleted)

### Value Objects

Objects without identity, defined by their attributes.

```go
type Status string

const (
    StatusPending   Status = "pending"
    StatusCompleted Status = "completed"
)
```

**Characteristics**:
- Immutable
- No identity
- Defined by attributes
- Behavior without side effects

### Aggregates

A cluster of entities and value objects with defined boundaries.

```go
// Todo is an aggregate root
type Todo struct {
    entity.Base
    Title       string
    Description string
    Status      Status
    // Could have todo items as part of the aggregate
}
```

**Characteristics**:
- Has a root entity
- Enforces invariants
- Transactional boundary
- Referenced by ID from outside

### Repositories

Provide persistence abstraction for aggregates.

```go
type Repo struct {
    db *gorm.DB
}

func (r *Repo) Save(ctx context.Context, todo *Todo) error {
    return r.db.WithContext(ctx).Save(todo).Error
}

func (r *Repo) View(ctx context.Context, id uuid.UUID) (*Todo, error) {
    var todo Todo
    err := r.db.WithContext(ctx).First(&todo, "id = ?", id).Error
    return &todo, err
}
```

**Characteristics**:
- Collection-like interface
- Encapsulates persistence
- Returns domain objects
- Maintains aggregate boundaries

### Domain Services

Operations that don't naturally fit on entities.

```go
type Service struct {
    repo *Repo
}

func (s *Service) Create(ctx context.Context, req CreateReq) (*Todo, error) {
    // Orchestrate domain operations
    todo := &Todo{
        Title:       req.Title,
        Description: req.Description,
        Status:      StatusPending,
    }
    
    if err := s.repo.Save(ctx, todo); err != nil {
        return nil, err
    }
    
    return todo, nil
}
```

**Characteristics**:
- Stateless
- Operate on multiple aggregates
- Coordinate domain logic
- Don't contain state

## Module Pattern

Each domain is organized as a self-contained module:

```go
// Module encapsulates all layers of a domain
type Module[R any, S any] struct {
    Repo    R      // Repository layer
    Service S      // Service layer
    Router  Router // Handler layer
}

// Example instantiation
todoModule := rest.Module[*todo.Repo, *todo.Service]{
    Repo:    todo.NewRepo(db),
    Service: nil,
    Router:  nil,
}
todoModule.Service = todo.NewService(todoModule.Repo)
todoModule.Router = todo.NewHandler(*todoModule.Service)
```

**Benefits**:
- Encapsulation
- Reusability
- Clear dependencies
- Easy testing

## Request Flow

### HTTP Request Lifecycle

```
1. HTTP Request
   ↓
2. Fiber Framework
   ↓
3. Middlewares (I18n, IP Detection, Recovery)
   ↓
4. Router (Route matching)
   ↓
5. Handler Wrapper (rest.Handle)
   ↓
6. Request Parsing (WithBody, WithParams, WithQuery)
   ↓
7. Validation (WithValidation)
   ↓
8. Domain Service
   ↓
9. Repository
   ↓
10. Database
    ↓
11. Response Transformation
    ↓
12. HTTP Response
```

### Example Flow

```go
// 1. Route registration
group.Post("/",
    srv.Timeout(
        rest.Handle(
            rest.WithBody(
                rest.WithValidation(
                    srv.ValidateStruct(),
                    rest.CreateResponds(h.srv.Create)
                )
            )
        )
    )
)

// 2. Request comes in
// POST /todos
// {"title": "Buy milk", "description": "2% milk"}

// 3. WithBody parses JSON → CreateReq
// 4. WithValidation validates CreateReq
// 5. h.srv.Create(ctx, req) - domain service
// 6. repo.Save(ctx, todo) - persistence
// 7. CreateResponds transforms result
// 8. Response: 201 Created with todo data
```

### Handler Wrappers

Generic handler wrappers provide composable request processing:

```go
// WithBody parses request body
func WithBody[Req any, Res any](
    next func(c *fiber.Ctx, req Req) (Res, error),
) func(*fiber.Ctx) (Res, error)

// WithValidation validates request
func WithValidation[Req any, Res any](
    validator validation.Fn,
    next func(c *fiber.Ctx, req Req) (Res, error),
) func(*fiber.Ctx, Req) (Res, error)

// WithParams parses URL parameters
func WithParams[Req any, Res any](
    next func(c *fiber.Ctx, req Req) (Res, error),
) func(*fiber.Ctx) (Res, error)
```

**Composition**:
```go
// Combine wrappers
rest.Handle(
    rest.WithBody(              // Parse body
        rest.WithValidation(    // Validate
            srv.ValidateStruct(),
            rest.Data(          // Transform response
                h.srv.Update    // Call service
            )
        )
    )
)
```

## Error Handling

### Error Flow

```
Domain Error (business rule violation)
    ↓
Service Layer (wrap with context)
    ↓
Handler Layer (convert to HTTP status)
    ↓
Error Handler Middleware (format response)
    ↓
HTTP Response (structured error JSON)
```

### Error Types

**Domain Errors**:
```go
var (
    ErrNotFound     = errors.New("todo not found")
    ErrInvalidState = errors.New("invalid todo state")
)
```

**Service Layer**:
```go
func (s *Service) View(ctx context.Context, id uuid.UUID) (*Todo, error) {
    todo, err := s.repo.View(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to fetch todo: %w", err)
    }
    return todo, nil
}
```

**HTTP Layer**:
```go
func (s *Service) ErrorHandler() fiber.ErrorHandler {
    return func(c *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError
        
        if errors.Is(err, ErrNotFound) {
            code = fiber.StatusNotFound
        }
        
        return c.Status(code).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
}
```

## Testing Strategy

### Unit Tests

Test domain logic in isolation:

```go
func TestTodo_Complete(t *testing.T) {
    todo := &Todo{
        Status: StatusPending,
    }
    
    err := todo.Complete()
    assert.NoError(t, err)
    assert.Equal(t, StatusCompleted, todo.Status)
}
```

### Integration Tests

Test with real dependencies:

```go
func TestRepo_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := NewRepo(db)
    
    todo := &Todo{Title: "Test"}
    err := repo.Save(context.Background(), todo)
    
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, todo.ID)
}
```

### Handler Tests

Test HTTP endpoints:

```go
func TestHandler_Create(t *testing.T) {
    app := fiber.New()
    handler := NewHandler(mockService)
    handler.RegisterRoutes(mockRestService, app)
    
    req := httptest.NewRequest("POST", "/todos", 
        strings.NewReader(`{"title":"Test"}`))
    resp, _ := app.Test(req)
    
    assert.Equal(t, 201, resp.StatusCode)
}
```

## Best Practices

### 1. Keep Domain Pure

```go
// Good: Pure domain logic
func (t *Todo) Complete() error {
    if t.Status == StatusArchived {
        return ErrCannotCompleteArchivedTodo
    }
    t.Status = StatusCompleted
    return nil
}

// Bad: Framework dependency in domain
func (t *Todo) Complete(c *fiber.Ctx) error { // ❌
    // ...
}
```

### 2. Use Value Objects

```go
// Good: Type-safe value object
type Status string

const (
    StatusPending Status = "pending"
)

// Bad: Primitive obsession
status := "pending" // Just a string
```

### 3. Encapsulate Behavior

```go
// Good: Behavior with data
type Todo struct {
    Status Status
}

func (t *Todo) Archive() {
    t.Status = StatusArchived
}

// Bad: Anemic model
type Todo struct {
    Status Status
}
// Logic elsewhere
```

### 4. Define Clear Boundaries

```go
// Repository interface in domain
type Repository interface {
    Save(ctx context.Context, todo *Todo) error
    View(ctx context.Context, id uuid.UUID) (*Todo, error)
}

// Implementation in infrastructure
type Repo struct {
    db *gorm.DB
}
```

### 5. Use Dependency Injection

```go
// Good: Dependencies injected
func NewService(repo *Repo) *Service {
    return &Service{repo: repo}
}

// Bad: Hidden dependencies
func NewService() *Service {
    db := connectDB() // ❌ Hidden dependency
    return &Service{repo: NewRepo(db)}
}
```

### 6. Handle Context Properly

```go
// Good: Context passed through
func (s *Service) Create(ctx context.Context, req CreateReq) (*Todo, error) {
    return s.repo.Save(ctx, todo)
}

// Bad: Ignoring context
func (s *Service) Create(req CreateReq) (*Todo, error) {
    return s.repo.Save(context.Background(), todo) // ❌
}
```

### 7. Validate at Boundaries

```go
// Good: Validate at HTTP layer
type CreateReq struct {
    Title string `json:"title" validate:"required,min=3"`
}

// Use WithValidation wrapper in handler

// Bad: Validation scattered everywhere
```

### 8. Return Errors, Don't Panic

```go
// Good: Return errors
func (s *Service) Create(ctx context.Context, req CreateReq) (*Todo, error) {
    if err := s.repo.Save(ctx, todo); err != nil {
        return nil, err
    }
    return todo, nil
}

// Bad: Panic on errors
func (s *Service) Create(ctx context.Context, req CreateReq) *Todo {
    if err := s.repo.Save(ctx, todo); err != nil {
        panic(err) // ❌
    }
    return todo
}
```

---

This architecture guide provides the foundation for building maintainable, scalable Go applications with idiogo. Follow these principles and patterns to create clean, testable code that stands the test of time.
