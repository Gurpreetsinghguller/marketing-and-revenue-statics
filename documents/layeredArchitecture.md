## Software Architecture

### Layered Architecture with Clean Code Principles

The system follows a **layered architecture pattern** with inverted dependencies (Onion Architecture) to maximize testability, maintainability, and independence from external frameworks:

```
┌─────────────────────────────────────────┐
│        HTTP Request                     │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│   Boundary Objects (HTTP Handlers)      │  ← Presentation Layer
│   - Validate requests                   │
│   - Transform DTOs                      │
│   - Route to use cases                  │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│   Control Objects (Use Cases)           │  ← Business Logic Layer
│   - Orchestrate domain operations       │
│   - Enforce business rules              │
│   - Coordinate repositories             │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│   Domain (Entity Objects)               │  ← Domain Layer
│   - Campaign, User, Event, Metrics      │
│   - Domain logic and constraints        │
│   - Repository interfaces               │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│   Data Access Layer (Repositories)      │  ← Persistence Layer
│   - Implementation of repository        │
│     interfaces                          │
│   - Database operations                 │
│   - Data transformation                 │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│        Data Storage Layer               │
│        (In-memory, DB, Cloud)           │
└─────────────────────────────────────────┘
```
