# Item PDP Service

A **Product Detail Page (PDP) microservice** built with **Domain-Driven Design (DDD)** principles using **Go 1.21** and modern cloud-native technologies.

## 📋 Project Overview

This service provides comprehensive **item/product management** functionality including CRUD operations, inventory management, search capabilities, and image handling through a RESTful API.

### **Core Business Domain**
- **Item Management**: Complete lifecycle management of product items
- **Inventory Tracking**: Real-time stock level management and availability
- **Categorization**: Hierarchical product categorization with slug support
- **Multi-Currency Pricing**: Flexible pricing with currency support
- **Image Management**: Multiple product images with primary designation
- **Dynamic Attributes**: Extensible key-value attributes (color, size, brand, etc.)
- **Status Lifecycle**: Draft → Active → Inactive → Archived workflow

## 🏗️ Architecture

### **3-Layer DDD Architecture**
```
├── 🏢 Application Layer
│   ├── DTOs (Data Transfer Objects)
│   ├── HTTP Handlers & Routes  
│   ├── Middleware (CORS, Logging, Validation)
│   └── Use Cases (Business Orchestration)
│
├── 🎯 Domain Layer  
│   ├── Entities (Item, rich business models)
│   ├── Value Objects (SKU, Price, Category, etc.)
│   ├── Domain Events
│   ├── Business Rules & Validation
│   └── Repository Interfaces
│
└── 🔧 Infrastructure Layer
    ├── PostgreSQL Repository Implementation
    ├── Database Connection & Migration Management
    ├── Configuration Management (Viper)
    └── External Service Integrations
```

### **Design Patterns Used**
- **Repository Pattern** for data persistence abstraction
- **Use Case Pattern** for application logic orchestration  
- **Value Objects** for domain model integrity
- **Dependency Injection** using Uber FX
- **Clean Architecture** with dependency inversion

## ✨ Features & API Endpoints

### **Core Item Management**
- `POST /api/v1/items` - Create new item
- `GET /api/v1/items/{id}` - Get item by ID
- `GET /api/v1/items/sku/{sku}` - Get item by SKU
- `PUT /api/v1/items/{id}` - Update item
- `DELETE /api/v1/items/{id}` - Delete item

### **Inventory Management**
- `PATCH /api/v1/items/{id}/inventory` - Update stock levels
- `GET /api/v1/items/available` - Get all available items

### **Status Management**
- `PATCH /api/v1/items/{id}/activate` - Activate item
- `PATCH /api/v1/items/{id}/deactivate` - Deactivate item

### **Search & Filtering**
- `GET /api/v1/items/search?query=...` - Full-text search
- `GET /api/v1/items/category/{category}` - Filter by category
- Advanced filtering by status, availability, price range

### **Image Management**
- `POST /api/v1/items/{id}/images` - Add product images
- Support for primary image designation and alt text

### **Health & Monitoring**
- `GET /health` - Service health check

## 🗃️ Database Schema

### **Items Table**
```sql
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price_amount BIGINT NOT NULL,        -- Stored in cents
    price_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    category_name VARCHAR(100) NOT NULL,
    category_slug VARCHAR(100) NOT NULL,
    inventory_quantity INTEGER NOT NULL DEFAULT 0,
    images JSONB DEFAULT '[]'::jsonb,    -- Flexible image storage
    attributes JSONB DEFAULT '{}'::jsonb, -- Dynamic attributes
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

### **Performance Optimizations**
- **Indexes**: SKU, category, status, inventory, timestamps
- **Partial Index**: Available items (status='active' AND inventory > 0)
- **Full-text Search**: GIN index for name/description search
- **Automatic Timestamps**: Trigger-based updated_at management

## 🛠️ Tech Stack

### **Core Technologies**
- **Go 1.21** - Primary language with modern features
- **Gin Web Framework** - Fast HTTP router and middleware
- **PostgreSQL** - Primary database with JSONB support
- **Uber FX** - Dependency injection framework

### **Key Dependencies**
- **Database**: `lib/pq` (PostgreSQL driver)
- **Configuration**: `spf13/viper` (YAML/env config)
- **Logging**: `rs/zerolog` (Structured logging)
- **Validation**: `go-playground/validator` (Request validation)
- **UUID**: `google/uuid` (ID generation)

### **Development & Testing**
- **Testing**: `stretchr/testify` with `go-sqlmock`
- **Code Coverage**: Built-in Go coverage tools
- **Linting**: `golangci-lint` configuration
- **Hot Reload**: `air` for development

## 📂 Project Structure

```
item-pdp-service/
├── cmd/api/main.go                 # Application entry point
├── internal/                       # Private application code
│   ├── application/                # Application layer
│   │   ├── dto/                    # Data Transfer Objects
│   │   ├── http/                   # HTTP handling
│   │   │   ├── handlers/           # HTTP handlers
│   │   │   ├── middleware/         # HTTP middleware
│   │   │   └── routes/             # Route configuration
│   │   └── usecase/                # Use case implementations
│   ├── domain/item/                # Domain layer
│   │   ├── entity.go               # Item entity
│   │   ├── value_objects.go        # Domain value objects
│   │   ├── repository.go           # Repository interface
│   │   ├── events.go               # Domain events
│   │   └── errors.go               # Domain errors
│   └── infrastructure/             # Infrastructure layer
│       ├── config/                 # Configuration management
│       ├── database/               # Database connection
│       └── persistence/            # Repository implementations
├── migrations/                     # Database migrations
├── configs/config.yaml             # Default configuration
├── docker-compose.yml              # Local development setup
├── Dockerfile                      # Container definition
├── Makefile                        # Development commands
└── docs/                           # Additional documentation
```

## 🚀 Getting Started

### **Prerequisites**
- Go 1.21 or higher
- PostgreSQL 13+ 
- Docker & Docker Compose (optional)
- Make (optional, for convenience commands)

### **Quick Start**
```bash
# 1. Clone the repository
git clone <repository-url>
cd item-pdp-service

# 2. Start PostgreSQL (using Docker)
docker-compose up -d postgres

# 3. Install dependencies
go mod download

# 4. Run database migrations
make migrate-up

# 5. Start the service
make run

# Service will be available at http://localhost:8080
```

### **Using Make Commands**
```bash
# Development setup
make setup          # Install dev dependencies
make deps           # Install all required tools

# Testing
make test           # Run all tests  
make test-coverage  # Run tests with coverage report
make test-watch     # Watch mode for development

# Code quality  
make lint           # Run linter
make format         # Format code
make security       # Security scan

# Database
make migrate-up     # Apply migrations
make migrate-down   # Rollback migrations
make db-setup       # Setup local database

# Docker
make docker-build   # Build Docker image
make docker-run     # Run in container
make docker-up      # Start all services
```

## 🧪 Testing

### **Test Coverage**
- **Domain Layer**: 86.7% coverage with comprehensive entity tests
- **Use Case Layer**: 36.8% coverage with business logic tests  
- **Repository Layer**: 57.6% coverage with database integration tests
- **HTTP Handler Layer**: 32.5% coverage with API endpoint tests
- **Total**: 74 test cases across all architectural layers

### **Test Types**
- **Unit Tests**: Domain entities and value objects
- **Integration Tests**: Database repository layer
- **API Tests**: HTTP handler endpoints
- **Mock Tests**: Use case business logic

### **Running Tests**
```bash
# Run all tests
make test

# With coverage report
make test-coverage

# Watch mode for development
make test-watch

# Verbose output
make test-verbose
```

## 📊 Configuration

### **Environment Variables**
```bash
# Application
APP_NAME=item-pdp-service
APP_VERSION=1.0.0
APP_ENVIRONMENT=development

# Server
SERVER_HOST=0.0.0.0  
SERVER_PORT=8080

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_NAME=item_pdp_db

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### **Configuration Files**
- `configs/config.yaml` - Default configuration
- `env.example` - Environment variable template
- Support for multiple environments (dev/staging/production)

## 🐳 Docker Support

### **Development**
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### **Production**
```bash
# Build image
docker build -t item-pdp-service .

# Run container
docker run -p 8080:8080 item-pdp-service
```

## 📈 Production Considerations

### **Performance Features**
- Connection pooling with configurable limits
- Database query optimization with proper indexes
- Graceful shutdown with configurable timeouts
- Structured logging for monitoring and debugging

### **Scalability**
- Stateless design for horizontal scaling
- Database connection management
- Configurable timeouts and limits
- Health check endpoint for load balancers

### **Security**
- Input validation on all endpoints
- CORS middleware configuration
- SQL injection prevention through prepared statements
- Comprehensive error handling without information leakage

## 📚 Additional Documentation

- [`TESTING_GUIDE.md`](./TESTING_GUIDE.md) - Comprehensive testing documentation
- [`flaw.md`](./flaw.md) - Code quality and architecture analysis
- [`scratchbook.md`](./scratchbook.md) - Development notes and insights

## 🎯 Code Review Challenge: Finding & Fixing Flaws

### **Objective**
This codebase contains **10 intentional flaws** designed as a training exercise for code review skills, security awareness, and architectural best practices. Your mission is to identify all flaws and create merge requests to fix them.

### **Types of Flaws to Find**

#### **🔒 Security Vulnerabilities **
- Common security issues that could expose the application to various attack vectors

#### **🏗️ Architecture Anti-Patterns **
- Design patterns that violate Domain-Driven Design and Clean Architecture principles

#### **⚡ Performance Issues **
- Code patterns that could impact application performance and scalability

### **Finding the Flaws**

#### **🕵️ Detection Strategy**
Use systematic code review techniques and available development tools to identify issues

#### **📁 Areas to Focus**
Review all architectural layers and their implementations for potential issues

#### **🔍 Red Flags to Look For**
- Code patterns that deviate from security, architecture, and performance best practices

### **Submission Process**

Create appropriate branches and merge requests to fix identified issues following standard development practices.

### **📞 Getting Help**

Refer to the project documentation and use the available development tools for guidance

---

**Happy Bug Hunting! 🐛🔍**

Remember: The goal is learning, not just finding flaws. Understanding *why* each issue is problematic and *how* to prevent similar issues is more valuable than speed.

## 🤝 Development Workflow

### **Git Workflow**
1. Create feature branch from `main`
2. Implement feature with tests
3. Ensure all tests pass: `make test`
4. Check code quality: `make lint`
5. Submit pull request with description

### **Code Standards**
- Go formatting with `gofmt`
- Linting with `golangci-lint`
- Test coverage maintenance
- Meaningful commit messages
- Documentation updates

---

**Built with ❤️ using Domain-Driven Design, Clean Architecture, and Go best practices.** 