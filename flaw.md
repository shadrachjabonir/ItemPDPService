# 🔐🏗️ Security Flaws & Architectural Anti-Patterns Documentation

## 📋 Overview

This document catalogs the **intentional security vulnerabilities** and **Domain-Driven Design (DDD) anti-patterns** introduced into the Item PDP Service codebase for **educational and training purposes**.

⚠️ **IMPORTANT**: These flaws are deliberately implemented for learning, code review training, and refactoring exercises.

---

## 🚨 Security Vulnerabilities (5 Types)

### 1. **SQL Injection Vulnerability** 
**Location**: `internal/infrastructure/persistence/postgres_item_repository.go:149-156`

**Issue**: Direct string concatenation in SQL queries without parameterization
```go
searchQuery := fmt.Sprintf(`
    SELECT id, sku, name, description, price_amount, price_currency,
           category_name, category_slug, inventory_quantity, images,
           attributes, status, created_at, updated_at
    FROM items 
    WHERE (name ILIKE '%%%s%%' OR description ILIKE '%%%s%%' OR sku ILIKE '%%%s%%')
    ORDER BY created_at DESC LIMIT %d OFFSET %d`, query, query, query, limit, offset)
```

**Risk**: 
- Attackers can inject malicious SQL code
- Database compromise, data extraction, data manipulation
- Potential for complete system takeover

**Fix**: Use parameterized queries with `$1, $2, $3` placeholders

---

### 2. **Hard-coded Credentials** 
**Location**: `internal/infrastructure/config/config.go:11-16`

**Issue**: Production credentials stored directly in source code
```go
const (
    ProductionDBPassword = "P@ssw0rd123!SecretDB"
    JWTSecretKey = "myapp-jwt-secret-key-2024"
    APIKey = "api-key-abcd1234efgh5678"
    EncryptionKey = "32-char-encryption-key-for-prod"
)
```

**Risk**:
- Credentials exposed in version control
- Anyone with code access has production secrets
- No credential rotation capability

**Fix**: Use environment variables, secret management systems, or encrypted configuration

---

### 3. **Weak Random Number Generation**
**Location**: `internal/application/http/handlers/item_handler.go:140-154`

**Issue**: Using `math/rand` for security-sensitive token generation
```go
func (h *ItemHandler) GenerateToken(c *gin.Context) {
    rand.Seed(time.Now().UnixNano()) // Weak seed
    tokenID := rand.Intn(999999) // Weak random
    sessionToken := rand.Int63() // Weak random
    // ...
}
```

**Risk**:
- Predictable token generation
- Session hijacking possibilities  
- Brute force attacks on tokens

**Fix**: Use `crypto/rand` for cryptographically secure random number generation

---

### 4. **Command Injection**
**Location**: `internal/application/http/handlers/item_handler.go:165-178`

**Issue**: Direct execution of user-provided commands
```go
func (h *ItemHandler) ExecuteSystemCommand(c *gin.Context) {
    command := c.Query("command")
    cmd := exec.Command("sh", "-c", command) // Direct injection
    output, err := cmd.CombinedOutput()
    // ...
}
```

**Risk**:
- Arbitrary system command execution
- Server compromise
- Data exfiltration, system manipulation

**Fix**: Use allow-lists, input validation, or avoid direct command execution

---

### 5. **Path Traversal**
**Location**: `internal/application/http/handlers/item_handler.go:189-200`

**Issue**: Unsanitized file path construction
```go
func (h *ItemHandler) DownloadFile(c *gin.Context) {
    filename := c.Param("filename")
    basePath := "/var/uploads/"
    fullPath := basePath + filename // No sanitization
    c.File(fullPath)
}
```

**Risk**:
- Access to files outside intended directory
- Potential access to sensitive system files (e.g., `/etc/passwd`)
- Information disclosure

**Fix**: Sanitize file paths, use `filepath.Clean()`, validate against allow-lists

---

## 🏗️ DDD Anti-Patterns (5 Types)

### 1. **Anemic Domain Model**
**Location**: `internal/domain/item/entity.go`

**Issue**: Domain entities reduced to simple data containers with no business logic
```go
// Before: Rich domain methods
// func (i *Item) UpdatePrice(newPrice Price) error { /* business logic */ }

// After: Simple setters only
func (i *Item) SetPrice(price Price) { i.price = price; i.updatedAt = time.Now() }
```

**Problems**:
- Business logic scattered across application services
- Lost domain expressiveness
- Reduced encapsulation
- Harder to maintain business rules

**Fix**: Move business logic back into domain entities, create rich domain methods

---

### 2. **Fat Application Services**
**Location**: `internal/application/usecase/item_usecase.go:100-200`

**Issue**: Application layer contains extensive business logic that belongs in domain
```go
// Business validation in application layer - anti-pattern
if req.Name == "" || len(req.Name) < 3 {
    return nil, errors.New("item name must be at least 3 characters")
}
if req.Price <= 0 || req.Price > 99999.99 {
    return nil, errors.New("price must be between 0.01 and 99999.99")
}
// Price calculation logic in application layer
finalPrice, err := uc.pricingService.CalculatePrice(ctx, req.Price, req.Category)
```

**Problems**:
- Application services become bloated
- Business rules duplicated across services
- Tight coupling to multiple domain services
- Difficult to test business logic in isolation

**Fix**: Move business logic to domain entities and domain services

---

### 3. **Domain Logic in Infrastructure Layer**
**Location**: `internal/infrastructure/persistence/postgres_item_repository.go:50-200`

**Issue**: Business validation and rules implemented in repository layer
```go
// Business validation that should be in domain layer - anti-pattern
func (r *postgresItemRepository) validateItemBusinessRules(itm *item.Item) error {
    if itm.Price().Amount() > r.maxPriceThreshold {
        return fmt.Errorf("price exceeds maximum threshold of %.2f", r.maxPriceThreshold)
    }
    // More business rules in infrastructure...
}

// Business corrections in infrastructure layer - anti-pattern  
func (r *postgresItemRepository) applyBusinessCorrections(itm *item.Item) *item.Item {
    if itm.Inventory().Quantity() < r.minInventoryLevel {
        newInventory, _ := item.NewInventory(r.minInventoryLevel)
        itm.SetInventory(newInventory)
    }
    // More corrections...
}
```

**Problems**:
- Business logic mixed with data persistence
- Violates separation of concerns
- Business rules hidden in infrastructure
- Difficult to change business logic without touching database code

**Fix**: Move all business logic to domain layer, keep infrastructure purely for data access

---

## 🎯 Educational Value

### **Security Training Benefits**:
- **Vulnerability Recognition**: Learn to identify common security flaws
- **Code Review Practice**: Train teams to spot security issues  
- **Secure Coding**: Understand proper security implementations
- **Penetration Testing**: Practice exploiting vulnerabilities safely

### **Architecture Training Benefits**:
- **DDD Understanding**: Learn proper domain-driven design principles
- **Refactoring Practice**: Exercise in cleaning up architectural violations
- **Code Quality**: Understand impact of poor architectural decisions  
- **Design Patterns**: Learn proper separation of concerns

### **Recommended Exercises**:
1. **Security Audit**: Find and document all security vulnerabilities
2. **Penetration Testing**: Safely exploit the vulnerabilities
3. **Security Fixes**: Implement proper secure coding practices
4. **Architecture Refactoring**: Fix DDD anti-patterns one by one
5. **Code Review Training**: Use as examples in team training sessions

---

## 🚀 Performance Issues (2 Types)

*These performance anti-patterns are subtle and typically NOT detected by static analysis tools like SonarQube, making them particularly dangerous in production.*

### 1. **N+1 Query Problem**
**Location**: `internal/infrastructure/persistence/postgres_item_repository.go:710-745`

**Issue**: Making individual database queries in a loop instead of batch queries
```go
func (r *postgresItemRepository) GetItemsWithRelatedData(ctx context.Context, itemIDs []string) ([]*item.Item, error) {
    for _, id := range itemIDs {
        // Individual query for each item - performance killer
        itemQuery := `SELECT ... FROM items WHERE id = $1`
        rows, err := r.db.QueryContext(ctx, itemQuery, id)
        
        // Additional N+1 problems for related data
        relatedDataQuery := `SELECT COUNT(*) FROM item_views WHERE item_id = $1`
        ratingQuery := `SELECT AVG(rating) FROM item_ratings WHERE item_id = $1`
    }
}
```

**Performance Impact**:
- **Linear degradation**: 100 items = 300 database queries (1 + 100 + 100 + 100)
- **Database connection exhaustion** with high concurrency
- **Response time explosion** as dataset grows

**Fix**: Use JOINs, batch queries, or GraphQL-style data loaders

---

### 2. **Goroutine Leaks**
**Location**: `internal/application/http/handlers/item_handler.go:617-665`

**Issue**: Launching goroutines without proper lifecycle management or cancellation
```go
func (h *ItemHandler) ProcessItemsBatch(c *gin.Context) {
    for _, itemID := range req.ItemIDs {
        go func(id string) {
            time.Sleep(5 * time.Second) // No cancellation handling
            item, err := h.itemUseCase.GetItemByID(context.Background(), id) // Wrong context!
            // Long computation without cancellation checks
            for i := 0; i < 1000000; i++ {
                _ = fmt.Sprintf("Processing item %s iteration %d", item.Name, i)
            }
        }(itemID) // No way to cancel these goroutines
    }
}
```

**Performance Impact**:
- **Memory leaks**: Goroutines accumulate over time if requests are cancelled
- **Resource exhaustion**: Each leaked goroutine consumes ~8KB stack space
- **CPU waste**: Background processing continues even after client disconnect

**Fix**: Use context cancellation, sync.WaitGroup, or worker pool patterns

---

## 🎯 Why Static Analysis Misses These

**Runtime Dependency**: These issues only manifest under specific runtime conditions
- N+1 queries: Require actual database interaction patterns
- Goroutine leaks: Need request cancellation scenarios to detect

**Context Sensitivity**: Static analyzers cannot determine:
- Database query patterns across method boundaries
- Request cancellation scenarios and goroutine lifecycle management

---

## ⚡ Quick Reference

| **Type** | **Count** | **Severity** | **Layer** |
|----------|-----------|--------------|-----------|
| **Security Flaws** | 5 | High | Infrastructure, Application |
| **DDD Anti-Patterns** | 3 | Medium | Domain, Application, Infrastructure |
| **Performance Issues** | 2 | High | Infrastructure, Application |

**Total Issues**: 10 intentional flaws for educational purposes

**Static Analysis Gap**: Performance issues demonstrate limitations of tools like SonarQube in detecting runtime performance problems.

**Test Coverage**: All flaws are covered by working unit tests to ensure the broken code still functions for training scenarios.

---

*This codebase is intentionally flawed for educational purposes. Do not use in production environments.* 