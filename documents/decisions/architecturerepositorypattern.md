Repository Pattern
Purpose: Abstracts data access logic and provides a collection-like interface for accessing domain objects.

Key Idea: The application doesn't know WHERE or HOW data is stored (database, files, memory, API). It just uses a consistent interface.

In your code:



Benefits:

UseCase layer doesn't know about JSON file storage
Easy to swap implementations (JSON → PostgreSQL → MongoDB)
Testable: mock the PersistenceDB interface
