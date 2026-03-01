## Developer Setup

### Prerequisites
- Go 1.25.7 (see [go.mod](go.mod))
- Git

### Clone and Install Dependencies
```bash
git clone https://github.com/Gurpreetsinghguller/marketing-and-revenue-statics.git
cd marketing-and-revenue-statics
go mod tidy
```

### Configure JWT Secret
The API uses an HMAC JWT secret for signing and verification.

Create a secret file:
```bash
mkdir -p shared
echo -n "your-strong-secret" > shared/secret
```

Example secret file content:
```text
your-strong-secret
```

Load the secret into the environment for the API server:
```bash
export JWT_SECRET="$(cat shared/secret)"
```

### Configure Application YAML
Create or update [config/config.yml](config/config.yml):

```yaml
server:
	port: "8080"

log:
	level: "info"

auth:
	secret_file: "shared/secret"

rate_limit:
	max_requests: 100
	window_seconds: 60
```

### Run the API Server
```bash
go run ./cmd
```

The API is available at:
```
http://localhost:8080/api/v1/health        (Health check)
http://localhost:8080/api/v1/docs          (OpenAPI spec)
http://localhost:8080/api/v1/auth/register (Example endpoint)
```

Notes:
- Server port and log level are loaded from `config/config.yml`.
- `shared/secret` is ignored by git.


### Useful Commands
Format code:
```bash
gofmt -w ./cmd ./internal
```

Build binaries:
```bash
go build -o bin/statistics ./cmd/statistics
go build -o bin/tokengen ./cmd/tokengen
```

### OpenAPI Spec
The OpenAPI spec is served at `http://localhost:8080/api/v1/docs` when the server is running, or import [api/openapi.yaml](api/openapi.yaml) into Postman to generate requests for all endpoints.
