## Architecture


- We have used domain driven design so that your core logic should revolve around the domain.

- We have used layered architecture to separate concerns and improve maintainability. The layers are as follows:

 Layered Architecture (Decomposition) 
HTTP Request
    ↓
Handler Layer (Presentation)
    ↓
UseCase Layer (Business Logic)
    ↓
Persistence Layer (Data Access)


## SOLID Principles

We adhere to the following SOLID principles in our architecture:

- **Single Responsibility Principle**: Each class and module has a single, well-defined responsibility.
- **Open/Closed Principle**: Our code is open for extension but closed for modification.
- **Liskov Substitution Principle**: Subtypes are substitutable for their base types without breaking functionality.
- **Interface Segregation Principle**: Clients depend on specific, focused interfaces rather than broad ones.
- **Dependency Injection**: Dependencies are injected rather than created internally, promoting loose coupling and testability.

## Onion Architecture

We have adopted Onion Architecture (also known as Ports and Adapters) to ensure our domain logic remains independent of external concerns. This approach places the domain model at the core, surrounded by application services, followed by infrastructure and presentation layers. This inversion of dependencies ensures that external frameworks and technologies depend on our business logic, not vice versa.