HTTP Request → Router → Middleware → Handler → UseCase → Database
                                                  ↑
                              (business logic & DB operations)


                              1. Usecase Layer (Core Business Logic)
Implemented 6 usecase packages with database abstraction:

usecase.go - User registration, login, token verification
usecase.go - CRUD operations for campaigns with filtering
usecase.go - User profile management
usecase.go - Event tracking and aggregation
usecase.go - Analytics reports and metrics calculation
usecase.go - User engagement and funnel analysis
2. Handler Layer (HTTP Request Processing)
Implemented 6 handler packages that call usecase methods:

handler.go - RegisterHandler, LoginHandler with validation
handler.go - CRUD handlers with query filtering & pagination
handler.go - Profile GET/UPDATE handlers
handler.go - Event tracking & retrieval handlers
handler.go - Report generation handlers
handler.go - Engagement & funnel handlers
3. Middleware Layer (Cross-cutting Concerns)
Fully implemented in middleware.go:

AuthMiddleware - JWT token validation & user context extraction
RoleMiddleware - Role-based access control (Admin, Marketer, Analyst)
RateLimitMiddleware - Request rate limiting per user/IP
CORSMiddleware - Cross-origin resource sharing
LoggingMiddleware - Request/response logging
Key Design Decisions:
✅ Handler decoupling from DB - Handlers call usecase methods, not DB directly
✅ Consistent error handling - Proper HTTP status codes and error messages
✅ Context-based user info - User ID passed through request context and headers
✅ Input validation - Validation in both handlers and usecases
✅ Flexible filtering - Query parameters for filtering in GET operations
✅ Authorization checks - Both authentication and authorization middleware in place

No Breaking Changes:
All changes are backward compatible with the existing router structure in http.go. The router already imports all handlers and middleware.




----------------------------
 Layered Architecture (Decomposition) 
HTTP Request
    ↓
Handler Layer (Presentation)
    ↓
UseCase Layer (Business Logic)
    ↓
Persistence Layer (Data Access)