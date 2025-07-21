# 🧪 Testing Guide - Item PDP Service

## Prerequisites

### 1. Install Go
```bash
# macOS (using Homebrew)
brew install go

# Or download from https://golang.org/dl/
```

### 2. Install Development Tools
```bash
# Install all development dependencies
make deps

# Or install manually:
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### 3. Download Dependencies
```bash
# Download all Go modules
go mod download
go mod tidy
```

## 🚀 Quick Start

### Run All Tests
```bash
# Simple test run
make test

# With verbose output
make test-verbose

# With race detection
go test -v -race ./...
```

### Check Test Coverage
```bash
# Generate coverage report
make test-coverage

# Check against 80% target
make test-coverage-target

# Run our custom coverage script
./test_coverage.sh
```

## 📊 Test Categories

### 1. Unit Tests (Domain Layer)
```bash
# Test business logic
go test -v ./internal/domain/item/...

# Expected Output:
# ✅ TestNewItem
# ✅ TestItem_UpdatePrice  
# ✅ TestItem_UpdateInventory
# ✅ TestItem_AddImage
# ✅ TestNewSKU
# ✅ TestNewPrice
# ✅ TestNewCategory
# ✅ TestDomainError_Is
```

### 2. Application Layer Tests
```bash
# Test use cases with mocks
go test -v ./internal/application/usecase/...

# Test HTTP handlers
go test -v ./internal/application/http/handlers/...

# Expected Output:
# ✅ TestItemUseCase_CreateItem
# ✅ TestItemUseCase_GetItemByID
# ✅ TestItemHandler_CreateItem
# ✅ TestItemHandler_GetItem
```

### 3. Infrastructure Tests
```bash
# Test repository with SQL mocking
go test -v ./internal/infrastructure/persistence/...

# Expected Output:
# ✅ TestPostgresItemRepository_Save
# ✅ TestPostgresItemRepository_FindByID
# ✅ TestPostgresItemRepository_Update
```

### 4. Integration Tests
```bash
# Test full API workflows
go test -v ./test/integration/...

# Test configuration
go test -v ./test/...

# Expected Output:
# ✅ TestItemServiceIntegration
# ✅ TestItemServiceErrorHandling
# ✅ TestConfigLoad
```

## 🎯 Coverage Analysis

### Target: 80% Coverage

Our comprehensive test suite should achieve **80%+ coverage** across:

- **Domain Layer**: ~95% (Business logic critical)
- **Application Layer**: ~90% (Use cases and handlers)  
- **Infrastructure Layer**: ~85% (Repository operations)
- **Integration Layer**: ~80% (End-to-end workflows)

### Coverage Breakdown
```bash
# Detailed coverage by function
make test-coverage-func

# Coverage by package
go tool cover -func=coverage.out

# HTML report (opens in browser)
make test-coverage && open coverage.html
```

## 🔧 Development Workflow

### Test-Driven Development
```bash
# Watch mode - auto-run tests on changes
make test-watch

# Run specific test
go test -v -run TestItemUseCase_CreateItem ./internal/application/usecase/

# Run tests for specific package
go test -v ./internal/domain/item/
```

### Quality Checks
```bash
# Run linter
make lint

# Format code
make format

# Security scan
make security

# Full CI pipeline
make ci
```

## 🧪 Test Structure Overview

```
test/
├── integration/
│   └── item_integration_test.go     # Full API workflow tests
├── testutils/
│   └── test_helpers.go              # Reusable test utilities
└── config_test.go                   # Configuration tests

internal/
├── domain/item/
│   ├── entity_test.go               # Business logic tests
│   ├── value_objects_test.go        # Value object validation
│   └── errors_test.go               # Domain error handling
├── application/
│   ├── usecase/
│   │   └── item_usecase_test.go     # Use case logic with mocks
│   └── http/handlers/
│       └── item_handler_test.go     # HTTP endpoint tests
└── infrastructure/
    └── persistence/
        └── postgres_item_repository_test.go  # Repository tests
```

## 📈 Expected Test Results

### ✅ Successful Run Output
```
=== RUN   TestNewItem
=== RUN   TestNewItem/valid_item_creation
=== RUN   TestNewItem/empty_name_should_fail
--- PASS: TestNewItem (0.00s)
    --- PASS: TestNewItem/valid_item_creation (0.00s)
    --- PASS: TestNewItem/empty_name_should_fail (0.00s)

=== RUN   TestItemUseCase_CreateItem
=== RUN   TestItemUseCase_CreateItem/successful_creation
=== RUN   TestItemUseCase_CreateItem/duplicate_SKU_error
--- PASS: TestItemUseCase_CreateItem (0.01s)

🧪 Running tests with coverage...
📊 Coverage report generated: coverage.html
📈 Current coverage: 87.3%
🎯 Target coverage: 80%
✅ Coverage target met! (87.3% >= 80%)

🎉 Test suite summary:
  - All tests passing
  - Coverage above 80% threshold
  - Ready for production!
```

## 🚨 Troubleshooting

### Common Issues

1. **Import Errors**
   ```bash
   go mod tidy
   go mod download
   ```

2. **Mock Errors**
   ```bash
   # Ensure testify is available
   go get github.com/stretchr/testify
   ```

3. **Database Mock Issues**
   ```bash
   # Ensure sqlmock is available  
   go get github.com/DATA-DOG/go-sqlmock
   ```

4. **Coverage Below Target**
   ```bash
   # See detailed breakdown
   go tool cover -func=coverage.out | grep -v "100.0%"
   
   # Add tests for uncovered functions
   ```

## 🎯 Coverage Goals by Component

| Component | Target | Priority |
|-----------|--------|----------|
| Domain Entities | 95%+ | Critical |
| Value Objects | 95%+ | Critical |
| Use Cases | 90%+ | High |
| HTTP Handlers | 85%+ | High |
| Repository | 85%+ | High |
| Configuration | 80%+ | Medium |
| Integration | 80%+ | Medium |

## 📝 Adding New Tests

### For New Features
1. Start with domain tests (TDD approach)
2. Add use case tests with mocks
3. Add handler tests for API endpoints
4. Add integration tests for workflows
5. Update coverage target if needed

### Test Naming Convention
- `TestComponentName_MethodName` for unit tests
- `TestComponentName_MethodName_Scenario` for specific scenarios
- Use descriptive scenario names: `successful_creation`, `invalid_input`, `database_error`

## 🔄 Continuous Integration

The test suite is designed to run in CI/CD pipelines:

```bash
# CI Pipeline command
make ci

# This runs:
# 1. go mod tidy
# 2. Tests with coverage check (80% threshold)  
# 3. Linting with golangci-lint
# 4. Exits with non-zero code if any step fails
```

---

**Ready to test?** Run `make test-coverage` to get started! 🚀 