# **NumerisBook - Invoice Management System**

## **Table of Contents**
1. [Architecture Overview](#architecture-overview)
2. [Project Structure](#project-structure)
3. [Dependency Injection](#dependency-injection)
4. [API Endpoints](#api-endpoints)
5. [Scalability](#scalability)
6. [Error Handling](#error-handling)

---

## **Architecture Overview**

### **Clean Architecture**
The project adheres to Clean Architecture principles, ensuring a clear separation of concerns across layers. The structure is organized as follows:

```
pkg/ 
├── controllers/ # HTTP request handlers 
├── services/ # Business logic 
├── repositories/ # Data access 
├── models/ # Domain models 
├── dtos/ # Data Transfer Objects 
├── interfaces/ # Layer interfaces 
└── common/ # Shared utilities
```


### **Key Components**
1. **Controllers**: Handle HTTP requests and responses.
2. **Services**: Implement business logic and orchestrate operations.
3. **Repositories**: Manage data persistence and retrieval.
4. **Models**: Define core business entities.
5. **DTOs**: Define request/response data structures.
6. **Interfaces**: Define contracts between layers.

### **Flow of Control**
The typical flow of control in the system is as follows:
```
HTTP Request → Controller → Service → Repository → Database
↑ ↑ ↑
└─ DTOs ─────┴── Models ┘
```


---

## **Project Structure**

### **Core Packages**

```
pkg/
├── controllers/
│ ├── interfaces/
│ └── implementations/
├── services/
│ ├── interfaces/
│ └── implementations/
├── repositories/
│ ├── interfaces/
│ └── implementations/
├── models/
├── dtos/
│ ├── request/
│ └── response/
└── common/
├── exceptions/
└── helper/
```


### **Key Features**
- **Clear Separation of Concerns**: Independent layers for modularity and maintainability.
- **Interface-Driven Development**: Facilitates testability and flexibility.
- **Dependency Injection**: Promotes clean, decoupled code.
- **Centralized Error Handling**: Ensures consistent and structured error responses.
- **Audit Logging**: Tracks significant system events.
- **Automated Testing**: Mocks and unit tests for all critical components.

---

## **Dependency Injection**

### **Implementation**
Dependency injection is implemented using constructor-based methods, enabling seamless integration and mockability for unit tests. Example:

```go
func NewInvoiceController(
    logger zerolog.Logger,
    invoiceService services_interfaces.InvoiceService,
    auditService services_interfaces.AuditService,
    reminderService services_interfaces.ReminderService,
    customerService services_interfaces.CustomerService,
) controller_interfaces.InvoiceController {
    return &invoiceController{
        logger:          logger,
        invoiceService:  invoiceService,
        auditService:    auditService,
        reminderService: reminderService,
        customerService: customerService,
    }
}
```

## API Endpoints

### Authentication
All endpoints require authentication via the x-customer-id header:

```bash
x-customer-id: <customer_id>
```

### Documentation
Full API documentation is available here: https://documenter.getpostman.com/view/14136605/2sAYBYe9So


## Scalability

### Database Design
- Efficient indexing
- Proper relationship modeling
- Connection pooling
- Query optimization

### Performance Considerations
1. **Caching**
   - Response caching
   - Database query caching
   - In-memory caching for frequently accessed data

2. **Async Processing**
   - Background job processing and message queues for heavy operations

### Monitoring and Logging
- Structured logging with zerolog
- Audit trail for important operations
- Performance metrics
- Error tracking




## Error Handling

### Centralized Error Management
``` go
// Common error responses
func ThrowBadRequestException(ctx gin.Context, message string) {
ctx.JSON(http.StatusBadRequest, BuildErrorResponse(message))
}
```

### Error Types
1. Validation errors
2. Business logic errors
3. Database errors
4. External service errors
