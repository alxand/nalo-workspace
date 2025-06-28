# Configuration Guide

This document describes all configuration options available for the Nalo Workspace API.

## Environment Variables

Create a `.env` file in the root directory with the following variables:

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Port number for the HTTP server |
| `HOST` | `0.0.0.0` | Host address for the HTTP server |
| `READ_TIMEOUT` | `30s` | HTTP read timeout |
| `WRITE_TIMEOUT` | `30s` | HTTP write timeout |
| `IDLE_TIMEOUT` | `60s` | HTTP idle timeout |

### Database Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DRIVER` | `postgres` | Database driver (`postgres` or `sqlite`) |
| `DSN` | - | Database connection string (required) |
| `TEST_DSN` | `:memory:` | Test database connection string |
| `DB_MAX_CONNS` | `10` | Maximum database connections |

### JWT Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | - | JWT signing secret (required) |
| `JWT_EXPIRATION` | `24h` | JWT token expiration time |

### Logging Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `LOG_FORMAT` | `json` | Log format (`json` or `console`) |

## Example Configuration

```env
# Server Configuration
PORT=3000
HOST=0.0.0.0
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=60s

# Database Configuration
DB_DRIVER=postgres
DSN=host=localhost user=postgres password=password dbname=nalo_workspace port=5432 sslmode=disable TimeZone=UTC
TEST_DSN=:memory:

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

## Database Connection Strings

### PostgreSQL

```
host=localhost user=username password=password dbname=database port=5432 sslmode=disable TimeZone=UTC
```

### SQLite

```
file:./data.db?cache=shared&_foreign_keys=on
```

For testing (in-memory):
```
:memory:
```

## Production Considerations

1. **JWT Secret**: Use a strong, randomly generated secret in production
2. **Database**: Use connection pooling and proper SSL configuration
3. **Logging**: Use structured logging (JSON format) in production
4. **Environment**: Set appropriate timeouts for your use case
5. **Security**: Never commit `.env` files to version control

## Docker Environment

When using Docker, you can pass environment variables through:

1. **Docker Compose** (recommended):
   ```yaml
   environment:
     - PORT=3000
     - DSN=host=db user=postgres password=password dbname=nalo_workspace
   ```

2. **Docker run**:
   ```bash
   docker run -e PORT=3000 -e DSN="..." nalo-workspace
   ```

3. **Environment file**:
   ```bash
   docker run --env-file .env nalo-workspace
   ``` 