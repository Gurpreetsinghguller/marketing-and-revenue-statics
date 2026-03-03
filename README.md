# Marketing & Revenue Statistics API

A production-grade backend service for tracking marketing campaign performance, user engagement, and revenue analytics. Built with Go and architected using Clean Architecture, Domain-Driven Design, and SOLID principles.


## Core Features

### Authentication and Authorization
- JWT-based authentication with email and password login
- Role-based access control (RBAC) with three tiers:
  - **Admin**: Full system access
  - **Marketer**: Manage and view own campaigns and reports
  - **Analyst**: Read-only analytics access

### Campaign Management
- Create, update, and delete campaigns with full lifecycle management
- Activate, pause, or complete campaigns
- Public campaign previews for stakeholder sharing
- Advanced filtering and search by status, date range, channel, and creator

### Event Tracking
- High-volume ingestion of marketing events: Impression, Click, Conversion
- Time-series event storage optimized for analytics queries
- Extensible event model for future event types

### Analytics and Reporting
- Aggregated performance metrics (CTR, CPC, ROI, Conversion Rate)
- Daily, weekly, and monthly report summaries
- Campaign-level and user-level analytics with filterable reports

### User Engagement Tracking
- Session behavior and funnel analysis (Ad → Landing → Signup → Purchase)
- Drop-off detection to identify conversion bottlenecks
- User-level engagement insights

---
## Layered Architecture Here
[Layered-Architecture-Here](development/layeredArchitecture.md)

### Object Classification

- **Entity Objects**: Domain models representing core business concepts
  - `Campaign`: Marketing campaigns with budgets, channels, and lifecycle states
  - `User`: Application users with authentication credentials and roles
  - `Event`: Marketing events (impressions, clicks, conversions)
  - `Metrics`: Aggregated analytics and performance data

- **Boundary Objects**: API handlers that interact with external HTTP clients
  - Validate incoming requests
  - Transform and serialize responses
  - Extract context from HTTP headers (user ID, authorization)

- **Control Objects**: Use case orchestrators that execute business logic
  - `CampaignUseCase`: Campaign management operations
  - `AuthUseCase`: User authentication and registration
  - `AnalyticsUseCase`: Report generation and metrics aggregation
  - `EventUseCase`: Event ingestion and processing

---

## Design Patterns

### 1. **Repository Pattern**
**Purpose**: Abstract data access logic and provide a collection-like interface to domain objects.

**Implementation**: 
- `PersistenceDB` interface defines CRUD operations (`Create`, `Read`, `Update`, `Delete`, `List`)
- Concrete repositories (`CampaignRepository`, `UserRepository`, `EventRepository`) implement repository interfaces
- Use cases depend on repository interfaces, not implementations
- Enables seamless switching between JSON file storage, PostgreSQL, MongoDB, or other backends

**Benefit**: The business layer is decoupled from storage mechanism. Swapping databases requires only changing the repository implementation, not the use case logic.

```
UseCase → RepositoryInterface ← CampaignRepository (implements)
                               ← GraphQLAdapter (future)
                               ← PostgreSQLAdapter (future)
```

### 2. **Strategy Design Pattern**
**Purpose**: Define a family of algorithms, encapsulate each one, and make them interchangeable.

**Implementation**:
- Campaign filtering strategies based on user role:
  - Admin strategy: Return all campaigns
  - Marketer strategy: Return only own campaigns
  - Analyst strategy: Return public campaigns only
- Event processing strategies for different event types:
  - Impression strategy: Count events, update reach metrics
  - Click strategy: Track conversion path, update CTR
  - Conversion strategy: Process revenue, update ROI

**Benefit**: New filtering logic or event types can be added without modifying existing use cases. Strategies are swappable at runtime based on context (user role, event type).

### 3. **Factory Pattern**
**Purpose**: Create objects without specifying exact classes, promoting loose coupling.

**Implementation**:
- Factory functions in `main.go`:
  - `NewCampaignUseCase(repo)`: Creates campaign use case with dependency injection
  - `NewAuthUseCase(userRepo)`: Creates auth use case
  - `NewEventUseCase()`: Creates event use case
- Centralized object construction in application bootstrap phase
- Simplifies testing via mock factories

**Benefit**: Construction logic is centralized, making it easy to swap implementations. Reduces tight coupling between components.

### 4. **Dependency Injection Pattern**
**Purpose**: Provide dependencies to objects rather than having them create dependencies internally.

**Implementation**:
- Constructor-based injection in handlers and use cases:
  ```go
  func NewCampaignHandler(uc CampaignUseCaseInterface) *CampaignHandler {
      return &CampaignHandler{usecase: uc}
  }
  ```
- Interface-based dependencies: Handlers depend on `UseCaseInterface`, not concrete implementation
- All dependencies are injected in `main.go` during bootstrap

**Benefit**: Components are loosely coupled and highly testable. Mocks can be easily injected for unit testing.

### 5. **Adapter/Adapter Pattern**
**Purpose**: Convert the interface of a class into another interface clients expect.

**Benefit**: Decouples domain logic from HTTP specifics. HTTP layer can change without affecting business logic.

### 6. **Middleware/Chain of Responsibility Pattern**
**Purpose**: Pass requests along a chain of handlers, each deciding whether to process and forward.

**Implementation**:
- `AuthMiddleware`: Validates JWT tokens and adds user context to request
- `RateLimitMiddleware`: Enforces rate limits per user
- Middleware chain in router setup processes requests before reaching handlers

**Benefit**: Cross-cutting concerns (auth, logging, rate limiting) are separated from handler logic and can be reused across endpoints.

## SOLID Design Principles

### Single Responsibility Principle (SRP)
Each component has a single, well-defined responsibility:
- **Handlers**: Parse HTTP requests and format responses
- **Use Cases**: Orchestrate business logic
- **Repositories**: Handle data persistence
- **Domain Models**: Represent business entities and rules

Example: `CampaignHandler` only handles HTTP concerns; business logic is in `CampaignUseCase`.

### Open/Closed Principle (OCP)
Code is open for extension but closed for modification:
- New event types can be added via strategy pattern without modifying existing event processing
- New authorization roles can be added as new strategies without changing existing code
- New persistence backends can be added by implementing `PersistenceDB` interface

### Liskov Substitution Principle (LSP)
Subtypes are substitutable for their base types without breaking functionality:
- All repository implementations (`CampaignRepository`, `UserRepository`, `EventRepository`) adhere to `RepositoryInterface`
- Any implementation can be swapped without affecting use cases
- Mock repositories used in tests are substitutable for real repositories

### Interface Segregation Principle (ISP)
Clients depend on specific, focused interfaces rather than broad ones:
- `CampaignUseCaseInterface` defines only campaign-related operations
- `UserRepo` defines only user-related data operations
- `PersistenceDB` provides generic CRUD interface without forcing implementation of unused methods

### Dependency Inversion Principle (DIP)
High-level modules depend on abstractions, not low-level modules:
- Use cases depend on repository interfaces, not concrete implementations
- Handlers depend on use case interfaces, not concrete implementations
- Dependencies flow downward toward abstractions

---

## Four Pillars of Architecture

### 1. **Abstraction**
Hide implementation complexity behind well-defined interfaces.

**Application**:
- Repository interfaces abstract database operations
- Use case interfaces abstract business logic
- Domain models abstract real-world entities
- HTTP handlers abstract protocol-specific concerns

**Benefit**: Internal complexity is hidden; components interact via clean contracts.

### 2. **Encapsulation**
Bundle data and behavior together; expose only necessary operations.

**Application**:
- Domain models encapsulate entity data and invariants
- Use cases encapsulate business logic with clear method signatures
- Repositories encapsulate data access strategies
- Middleware encapsulates cross-cutting concerns (auth, rate limiting)

**Benefit**: Implementation details are hidden; only essential operations are exposed.

### 3. **Decomposition**
Break the system into smaller, manageable, independent components.

**Application**:
- Separate concerns across layers: handlers, use cases, repositories
- Organize features into modules: campaigns, events, analytics, auth, profile
- Each module contains handler, use case, and interface files
- Domain models are separate from persistence models

**Benefit**: Large system is easier to understand, test, and maintain. Changes in one module don't ripple through others.

### 4. **Generalization**
Build abstractions that apply to multiple contexts, avoiding duplication.

**Application**:
- `PersistenceDB` interface is generic CRUD, usable by all repositories
- Middleware pattern generalizes request processing logic
- Rate limiter is a generic utility applicable to any endpoint
- Validators are reusable across multiple handlers

**Benefit**: Code reuse is maximized; common patterns are centralized. Changes to shared logic benefit all consumers.

---

## Non-Functional Requirements (NFR) & Architecture Resilience

### Event Source Abstraction
**Requirement**: If events come from a message broker (Kafka, RabbitMQ) instead of HTTP in the future, the system should adapt without changing persistence or business logic.

**How Achieved**:
- Event handler accepts events regardless of source (HTTP, broker, batch file)
- Event use case depends on repository interface, not specific event source
- New event adapters (KafkaAdapter, FileImportAdapter) can be added by implementing the event repository interface

**Independence**: Persistence layer and business logic are unaffected by event source changes.

### Persistence Layer Flexibility
**Requirement**: Switching from in-memory or JSON file storage to a cloud-hosted/distributed database should not impact the core business logic.

**How Achieved**:
- All repositories implement the `PersistenceDB` interface
- Use cases depend on repository interfaces, not concrete implementations
- Business logic in use cases contains zero database-specific logic
- New persistence backends (PostgreSQL, MongoDB, DynamoDB, Firestore) can be added as implementations of `PersistenceDB`

**Independence**: Entire persistence layer can be replaced; use cases and domain logic remain untouched.

### HTTP Interface Layer Independence
**Requirement**: If API protocol changes (REST → GraphQL, REST → gRPC), handlers should be replaceable without affecting business logic.

**How Achieved**:
- Handlers are adapters that convert protocol-specific requests to domain operations
- Use cases accept domain objects, not HTTP requests
- Business logic is protocol-agnostic
- New handlers (GraphQL handlers, gRPC handlers) can be added without modifying use cases

**Independence**: Protocol layer is completely decoupled from business logic.

### Strategy for NFR Achievement
The system achieves these NFRs through:

1. **Dependency Inversion**: High-level logic depends on abstractions, not low-level details
2. **Interface-based Design**: Components interact via interfaces, making implementations swappable
3. **Layered Architecture**: Concerns are separated; changes in one layer don't cascade
4. **Domain-Driven Design**: Business logic is isolated from infrastructure concerns



## Development Setup
[here](documents/setup.md)

---

## Future Enhancements

- **Event Broker Integration**: Support for Kafka/RabbitMQ for high-volume event ingestion
- **Distributed Storage**: Integration with cloud databases (PostgreSQL, DynamoDB, Firestore)
- **Event Sourcing**: Event-driven architecture for audit trails
- **CQRS Pattern**: Separate read and write models for optimized queries

## Technology Stack

- Language: Go
- API: REST
- Authentication: JWT (HMAC)
- Configuration: YAML-based
- Testing: Go testing framework
- Documentation: OpenAPI


## Author

Gurpreet Singh  
Senior Software Engineer — Backend and Distributed Systems
