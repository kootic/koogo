# koogo

koogo is a production-ready boilerplate for Go API projects with built-in observability, database integration, and API documentation. It provides a solid foundation for building scalable, maintainable, and well-tested Go applications.

## Features

- 🚀 Production-ready Go API boilerplate
- 📊 Built-in OpenTelemetry integration for observability
- 📚 Swagger/OpenAPI documentation
- 🗄️ PostgreSQL database integration with Bun ORM
- 📝 Structured logging with zap
- 🐳 Docker Compose for local development
- 🧪 Comprehensive testing setup with integration tests
- 🔄 Database migrations with Atlas

## Under development

- ⏲️ Cron jobs
- 🚩 Feature flags with PostHog
- 🚨 Error tracking with PostHog
- 🛑 Rate limiting

## Prerequisites

- Go 1.24.2 or later
- Docker and Docker Compose
- PostgreSQL 17 (optional, provided via Docker)

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

3. Configure environment variables:
   Create a `.env` file in the project root with the following variables:

All environment variables are prefixed with `KOO_` (e.g., `KOO_APP_ENV`) to prevent conflicts with other applications and clearly identify variables specific to this application.

```sh
# Application
KOO_APP_ENV=local  # Options: local, dev, staging, prod
KOO_APP_NAME=koogo
KOO_APP_VERSION=0.0.1
KOO_APP_PORT=8080
KOO_APP_LOG_LEVEL=debug  # Options: debug, info, warn, error

# Database
KOO_DATABASE_HOST=localhost
KOO_DATABASE_PORT=5432
KOO_DATABASE_USERNAME=postgres
KOO_DATABASE_PASSWORD=postgres
KOO_DATABASE_DBNAME=koogo

# Swagger
KOO_SWAGGER_ENABLED=true
KOO_SWAGGER_USERNAME=swagger
KOO_SWAGGER_PASSWORD=swagger

# OpenTelemetry
KOO_OTEL_ENABLED=true
KOO_OTEL_EXPORTER=otlp-grpc  # Options: console, otlp-grpc, none
KOO_OTEL_OTLP_ENDPOINT=localhost:4317

# Grafana Cloud (optional)
GRAFANA_CLOUD_OTLP_ENDPOINT=your-grafana-cloud-otlp-endpoint
GRAFANA_CLOUD_INSTANCE_ID=your-grafana-cloud-instance-id
GRAFANA_CLOUD_API_KEY=your-grafana-cloud-api-key
```

For different environments, you can create environment-specific `.env` files:

- `.env.local` - Local development
- `.env.development` - Development environment
- `.env.production` - Production environment
- `.env.test` - Test environment

Note: `.env.test` is tracked in git to setup the test environment, while other `.env` files are gitignored.

4. Start the service:

```sh
# Using environment variables
KOO_APP_ENV=local go run cmd/koogo/main.go

# Or using VS Code launch configuration
# The launch.json file is already configured with the necessary environment variables
```

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
│   ├── repository/          # Data access layer
│   │   └── dbrepo/          # Database repository implementation
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

Integration tests are disabled by default. To run them, set the `RUN_INTEGRATION_TESTS` environment variable to `true`:

```sh
RUN_INTEGRATION_TESTS=true task test
```

Integration tests:

- Create a unique test database for each run
- Start the full application server
- Send real HTTP requests to test endpoints
- Clean up resources after completion

### Database Migrations

The project uses [Atlas](https://atlasgo.io/) for database migrations. Common commands:

```sh
# Check migration status
task atlas:status

# Generate migrations based on declared state
task atlas:diff name=add_user_table

# Create a manual migration
task atlas:manual name=change_configuration

# Apply pending migrations
task atlas:apply
```

Migration files are stored in `internal/repository/dbrepo/migrations/`.

#### Declarative Migrations

The project uses a declarative approach to database migrations:

1. Define the desired database state in SQL schema files under `internal/repository/dbrepo/migrations/sqlschema/`
2. Use `task atlas:diff` to generate migration files that transform the current state to the desired state
3. Review and apply the generated migrations

Use the main entry point (`sqlschema/main.sql`) to organize your schema definitions. This file serves as the central import hub where you can:

1. Import individual entity schema files using Atlas import directives
2. Control the order of schema application (important for dependencies)
3. Keep entity definitions modular by placing each in its own file

For example (`internal/repository/dbrepo/migrations/sqlschema/main.sql`):

```sql
-- atlas:import koo_users.sql
-- atlas:import koo_pets.sql
```

This approach:

- Makes it easy to understand the final database structure
- Reduces manual migration writing
- Helps prevent migration conflicts
- Allows for easy schema versioning

### Linting

```sh
task lint
```

### Generating API Documentation

```sh
task generate # This also generates any other code that relies on go:generate
```

## Observability

The project uses OpenTelemetry for distributed tracing, metrics, and logging. By default, it's configured to use the local OpenTelemetry Collector running in Docker.

To view traces and metrics:

1. The OpenTelemetry Collector is available at `http://localhost:4317`
2. If using Grafana Cloud, configure the environment variables to send data there

## API Documentation

When enabled, the Swagger UI is available at `http://localhost:8080/swagger/` with basic authentication.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
