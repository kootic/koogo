# koogo

koogo is a production-ready boilerplate for Go API projects with built-in observability, database
integration, and API documentation. It provides a solid foundation for building scalable,
maintainable, and well-tested Go applications.

## Features

- 🚀 Production-ready Go API boilerplate
- 📊 Built-in OpenTelemetry integration for observability
- 📚 Swagger/OpenAPI documentation
- 🗄️ PostgreSQL database integration with Bun ORM
- 🔄 Automated database migrations based on Bun models with Atlas
- 📝 Structured logging with zap
- 🐳 Docker Compose for local development
- 🧪 Comprehensive testing setup with integration tests

## Getting Started

1. Clone the repository:

```sh
git clone https://github.com/kootic/koogo.git
cd koogo
```

2. Set up the development environment:

```sh
# This will install dependencies, start local infrastructure, and set up git hooks
./scripts/boot.sh
```

3. Configure environment variables: Create a `cmd/koogo/.env` with the following variables:

All environment variables are prefixed with `KOO_` (e.g., `KOO_APP_ENV`) to prevent conflicts with
other applications and clearly identify variables specific to this application.

```sh
# Application
KOO_APP_ENV=local  # Options: local, dev, staging, prod
KOO_APP_NAME=koogo
KOO_APP_VERSION=0.0.1
KOO_APP_PORT=8080
KOO_APP_LOG_LEVEL=debug  # Options: debug, info, warn, error
KOO_APP_READ_TIMEOUT_SECONDS=15  # HTTP read timeout (default: 15)
KOO_APP_WRITE_TIMEOUT_SECONDS=15  # HTTP write timeout (default: 15)
KOO_APP_IDLE_TIMEOUT_SECONDS=120  # Connection idle timeout (default: 120)
KOO_APP_BODY_LIMIT_MB=4  # Max request body size in MB (default: 4)

# Database
KOO_DB_HOST=localhost
KOO_DB_PORT=5432
KOO_DB_USERNAME=postgres
KOO_DB_PASSWORD=postgres
KOO_DB_DATABASE=koogo
KOO_DB_MAX_CONNS=25  # Maximum connections in pool (default: 25)
KOO_DB_MIN_CONNS=5  # Minimum connections in pool (default: 5)
KOO_DB_MAX_CONN_LIFETIME_MINUTES=60  # Max connection lifetime (default: 60)
KOO_DB_MAX_CONN_IDLE_TIME_MINUTES=30  # Max connection idle time (default: 30)
KOO_DB_CONNECTION_TIMEOUT_SECONDS=10  # Connection timeout (default: 10)

# Swagger
KOO_SWAGGER_ENABLED=true
KOO_SWAGGER_USERNAME=swagger
KOO_SWAGGER_PASSWORD=swagger

# OpenTelemetry
KOO_OTEL_ENABLED=true
KOO_OTEL_EXPORTER=otlp-grpc  # Options: console, otlp-grpc, none; our own environment variable to control which exporter to use
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 # Standard OpenTelemetry environment variable used by the SDK
OTEL_EXPORTER_OTLP_INSECURE=true # Standard OpenTelemetry environment variable used by the SDK

# Grafana Cloud (optional)
GRAFANA_CLOUD_OTLP_ENDPOINT=your-grafana-cloud-otlp-endpoint
GRAFANA_CLOUD_INSTANCE_ID=your-grafana-cloud-instance-id
GRAFANA_CLOUD_API_KEY=your-grafana-cloud-api-key
```

4. Start the service:

Run:

```sh
cd cmd/koogo && go run main.go start
```

Or using VS Code launch configuration provided in `.vscode/launch.json`.

## Project Structure

```
koogo/
├── cmd/                     # Application entry points
│   └── koogo/               # Main application binary
├── deployment/              # Deployment configuration files
├── internal/                # Private application code
│   ├── app/                 # Application core
│   ├── config/              # Application configuration management
│   ├── dto/                 # Data transfer objects for request/response
│   ├── handler/             # HTTP handlers
│   ├── jobs/                # CLI job system (e.g., migrations)
│   ├── repo/                # Data access layer
│   │   └── postgres/        # PostgreSQL repository implementation
│   │       ├── bun/         # Bun ORM models (schema source of truth)
│   │       └── migrations/  # Database migration files
│   ├── server/              # HTTP server setup
│   ├── service/             # Business logic
│   └── tests/               # Integration tests
├── pkg/                     # Public packages
│   ├── kooctx/              # Context utilities
│   ├── koodb/               # Database client providers
│   ├── koohttp/             # HTTP utilities
│   ├── koolog/              # Logging utilities
│   └── kootel/              # OpenTelemetry utilities
├── scripts/                 # Utility scripts
│   └── boot.sh              # Development environment setup
└── swagger/                 # API documentation
```

## Development

### Running Tests

```sh
task test
```

### Integration Tests

Integration tests are disabled by default. To run them, set the `RUN_INTEGRATION_TESTS` environment
variable to `true`:

```sh
RUN_INTEGRATION_TESTS=true task test
```

Integration tests:

- Create a unique test database for each run
- Start the full application server
- Send real HTTP requests to test endpoints
- Clean up resources after completion

### Database Migrations

The project uses [Atlas](https://atlasgo.io/) with
[atlas-provider-bun](https://github.com/ariga/atlas-provider-bun) for database migrations. The
schema is defined using Bun ORM models in `internal/repo/postgres/bun/`.

Common commands:

```sh
# Check migration status
task atlas:status

# Generate migrations based on Bun models
task atlas:diff name=add_user_table

# Create a manual migration
task atlas:manual name=change_configuration

# Apply pending migrations
task atlas:apply
```

Migration files are stored in `internal/repo/postgres/migrations/`.

#### Declarative Migrations with Bun Models

The project uses Bun ORM models as the source of truth for the database schema:

1. Define or modify Bun models in `internal/repo/postgres/bun/`
2. Use `task atlas:diff` to generate migration files that transform the current database state to
   match your models
3. Review and apply the generated migrations

This approach:

- Uses Go structs as the single source of truth for both code and database schema
- Automatically generates migrations from model changes
- Reduces manual migration writing
- Helps prevent migration conflicts

### Linting

```sh
task lint
```

### Generating API Documentation

```sh
task generate # This also generates any other code that relies on go:generate
```

## Observability

The project uses OpenTelemetry for distributed tracing, metrics, and logging. By default, it's
configured to use the local OpenTelemetry Collector running in Docker.

To view traces and metrics:

1. The OpenTelemetry Collector is available at `http://localhost:4317`
2. If using Grafana Cloud, configure the environment variables to send data there

## API Documentation

When enabled, the Swagger UI is available at `http://localhost:8080/swagger/` with basic
authentication.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
