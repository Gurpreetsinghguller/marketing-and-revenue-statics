# Marketing & Revenue Statistics API

Scalable backend service for tracking marketing campaign performance, user engagement, and revenue analytics.

Built with Go, designed using clean architecture principles, and optimized for performance, scalability, and production use cases.

---

## Overview

This system provides REST APIs to:

- Manage marketing campaigns
- Track impressions, clicks, and conversions
- Analyze campaign performance
- Generate aggregated analytics reports
- Monitor user engagement funnels

Intended users include marketers, analysts, and administrators.

---

## Key Features

### Authentication and Authorization

- JWT-based authentication
- Email and password login
- Role-based access control (RBAC)

Roles:

| Role      | Permissions                          |
|-----------|--------------------------------------|
Admin       | Full system access                   |
Marketer    | Manage own campaigns and reports     |
Analyst     | Read-only analytics access           |

---

### Campaign Management

- Create, update, and delete campaigns
- Activate or pause campaigns
- Public campaign previews
- Search and filtering by:
  - Status
  - Date range
  - Channel
  - Creator

---

### Event Tracking

Supports ingestion of marketing events:

- Impression
- Click
- Conversion

Designed for high-volume event ingestion and time-series storage.

---

### Analytics and Reporting

Aggregated performance metrics:

| Metric            | Formula                                      |
|-------------------|----------------------------------------------|
CTR                | (Clicks / Impressions) × 100                  |
CPC                | Spend / Clicks                                |
ROI                | ((Revenue − Spend) / Spend) × 100             |
Conversion Rate    | (Conversions / Clicks) × 100                  |

Supports:

- Daily, weekly, and monthly summaries
- Campaign-level analytics
- User-level insights
- Filterable reports

---

### User Engagement Tracking

- Session behavior tracking
- Funnel analysis (Ad → Landing → Signup → Purchase)
- Drop-off detection between stages

---

## Architecture and Design

The system is designed to handle:

- Write-heavy workloads (event ingestion)
- Read-heavy workloads (analytics queries)
- Role-based data access
- Horizontal scalability

Key design considerations:

- Separation of ingestion and reporting pipelines
- Pre-aggregation of metrics
- Efficient querying of large datasets
- Cost-optimized storage for event logs

See Architecture.md for detailed explanation and design decisions.

---

## Technology Stack

- Language: Go
- API: REST
- Authentication: JWT (HMAC)
- Configuration: YAML-based
- Testing: Go testing framework
- Documentation: OpenAPI

---

## Setup Instructions

See the full developer setup guide:

Developer Setup Instructions: setup.md

Includes:

- Environment setup
- Configuration
- Running the server
- Testing
- Useful commands

---

## Running the Service

After setup:

```bash
go run ./cmd
```

API base URL:

```
http://localhost:8080/api/v1
```

Health check endpoint:

```
GET /api/v1/health
```

---

## API Documentation

OpenAPI documentation available at:

```
http://localhost:8080/api/v1/docs
```

You can also import:

```
api/openapi.yaml
```

into Postman.

---

## Testing

Run unit tests:

```bash
go test -v ./...
```

Generate coverage:

```bash
go test -cover ./...
```

---

## Project Structure

```
cmd/            Application entrypoints
internal/       Core business logic
config/         Configuration files
shared/         Secrets and shared resources
api/            OpenAPI specification
```

---

## Future Enhancements

- Real-time analytics pipeline
- Streaming aggregation
- Export support for large datasets
- Cost-optimized storage strategies
- Alerting (CTR drop, budget exceeded)
- Rate limiting

---

## Engineering Highlights

This project demonstrates:

- Clean API design
- Scalable architecture
- Production-grade configuration management
- Security best practices
- Performance-oriented design

---

## Author

Gurpreet Singh  
Software Engineer — Backend and Distributed Systems
