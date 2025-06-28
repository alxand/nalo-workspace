# Nalo Workspace API

A comprehensive Go Fiber + GORM backend API for managing daily tasks with JWT authentication and role-based authorization.

[![Go Report Card](https://goreportcard.com/badge/github.com/alxand/nalo-workspace)](https://goreportcard.com/report/github.com/alxand/nalo-workspace)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/alxand/nalo-workspace/ci.yml?branch=main&label=CI)](https://github.com/alxand/nalo-workspace/actions/workflows/ci.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/alexanderadade/nalo-workspace)](https://hub.docker.com/r/alexanderadade/nalo-workspace)

## Features

- 🔐 **JWT Authentication** - Secure token-based authentication
- 👥 **User Management** - User registration, login, and profile management
- 🛡️ **Role-Based Authorization** - Admin, Manager, and User roles
- 📝 **Daily Task Management** - Create, read, update, and delete daily tasks
- 🔒 **Data Isolation** - Users can only access their own tasks
- 📊 **Admin Dashboard** - User management for administrators
- 🗄️ **PostgreSQL Database** - Robust data persistence
- 📚 **Swagger Documentation** - Interactive API documentation
- 🐳 **Docker Support** - Easy deployment with Docker Compose
- 🧪 **Comprehensive Testing** - Unit tests with mocks

## Tech Stack

- **Framework**: Fiber (Go web framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT with bcrypt password hashing
- **Validation**: Go Playground Validator
- **Logging**: Zap logger
- **Documentation**: Swagger/OpenAPI
- **Testing**: Go testing with testify
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL
- Docker & Docker Compose (optional)

### Environment Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd nalo_workspace
```

2. Create a `.env` file:
```bash
# Server Configuration
PORT=3000
HOST=0.0.0.0
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=60s

# Database Configuration
DB_DRIVER=postgres
DSN=postgres://username:password@localhost:5432/nalo_workspace?sslmode=disable
DB_MAX_CONNS=10

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRATION=24h

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

### Using Docker (Recommended)

1. Start the application with Docker Compose:
```bash
docker compose up --build
```

2. Create an admin user:
```bash
make seed
```

3. Access the API:
- API: http://localhost:3000
- Swagger Docs: http://localhost:3000/swagger/

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Run migrations and create admin user:
```bash
make seed
```

3. Start the application:
```bash
make run
```

## Authentication & Authorization

### User Roles

- **Admin**: Full access to all features including user management
- **Manager**: Extended access to team-related features
- **User**: Basic access to personal daily tasks

### Authentication Flow

1. **Register** a new user account
2. **Login** to receive a JWT token
3. **Use the token** in the Authorization header for protected endpoints
4. **Refresh** the token when needed

### API Endpoints

#### Public Endpoints
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

#### Protected Endpoints (Require JWT)
- `GET /api/v1/auth/profile` - Get current user profile
- `POST /api/v1/auth/refresh` - Refresh JWT token

#### Daily Task Endpoints (Require JWT)
- `POST /api/v1/dailytask` - Create a new daily task
- `GET /api/v1/dailytask/:date` - Get tasks for a specific date
- `PUT /api/v1/dailytask/:id` - Update a task
- `DELETE /api/v1/dailytask/:id` - Delete a task

#### Admin Endpoints (Require Admin Role)
- `GET /api/v1/admin/users` - List all users
- `GET /api/v1/admin/users/:id` - Get user by ID
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user

### Example Usage

#### 1. Register a new user
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### 3. Create a daily task (with JWT token)
```bash
curl -X POST http://localhost:3000/api/v1/dailytask \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "day": "Monday",
    "date": "2024-01-15T00:00:00Z",
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-15T17:00:00Z",
    "status": "in_progress",
    "score": 8,
    "productivity_score": 85
  }'
```

## Project Structure

```
nalo_workspace/
├── cmd/
│   ├── api/          # Main application entry point
│   └── seed/         # Database seeding script
├── internal/
│   ├── api/          # HTTP handlers
│   │   ├── auth/     # Authentication handlers
│   │   ├── dailytask/ # Daily task handlers
│   │   └── user/     # User management handlers
│   ├── config/       # Configuration management
│   ├── container/    # Dependency injection container
│   ├── domain/       # Domain models and interfaces
│   ├── pkg/          # Shared packages
│   │   ├── errors/   # Error handling
│   │   ├── logger/   # Logging utilities
│   │   ├── middleware/ # HTTP middleware
│   │   └── validation/ # Validation utilities
│   ├── repository/   # Data access layer
│   └── server/       # Server setup and routing
├── docs/             # Swagger documentation
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Development

### Running Tests
```bash
make test
make test-coverage
```

### Code Quality
```bash
make lint
make fmt
```

### Generate Documentation
```bash
make swagger
```

### Development Setup
```bash
make dev-setup
```

## Security Features

- **Password Hashing**: Bcrypt with configurable cost
- **JWT Tokens**: Secure token-based authentication
- **Role-Based Access Control**: Fine-grained permissions
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM with parameterized queries
- **CORS Configuration**: Configurable cross-origin requests
- **Helmet Middleware**: Security headers

## Database Schema

### Users Table
- `id` (Primary Key)
- `email` (Unique)
- `username` (Unique)
- `password` (Hashed)
- `first_name`
- `last_name`
- `role` (admin/user/manager)
- `is_active`
- `last_login`
- `created_at`
- `updated_at`

### Daily Tasks Table
- `id` (Primary Key)
- `user_id` (Foreign Key)
- `day`
- `date`
- `start_time`
- `end_time`
- `status`
- `score`
- `productivity_score`
- `created_at`
- `updated_at`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue on GitHub.
