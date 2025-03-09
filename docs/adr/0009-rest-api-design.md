# ADR 009: REST API Design Pattern

## Status

Accepted

## Context

We needed to establish a consistent API design pattern that would:
- Support efficient client-server communication
- Maintain clear and consistent endpoints
- Enable easy versioning and evolution
- Support our security requirements
- Enable efficient caching
- Facilitate testing and documentation

## Decision

We have chosen to implement a RESTful API design with the following key characteristics:
- Resource-oriented endpoints
- Standard HTTP methods
- JSON as the primary data format
- JWT-based authentication
- Versioned endpoints
- Consistent error handling

## Consequences

### Advantages

1. **Consistency**
   - Predictable URL patterns
   - Standard HTTP status codes
   - Uniform request/response formats
   - Clear error structures
   - Consistent authentication flow

2. **Client Integration**
   - Easy to consume with React Query
   - Supports caching strategies
   - Clear contract for frontend development
   - Facilitates parallel development

3. **Maintainability**
   - Self-documenting URLs
   - Easy to version
   - Clear separation of concerns
   - Testable endpoints
   - Swagger/OpenAPI support

4. **Performance**
   - Efficient caching
   - Supports pagination
   - Allows partial responses
   - Enables bulk operations
   - Rate limiting support

### Disadvantages

1. **Overhead**
   - More endpoints to maintain
   - Additional documentation needed
   - Version management complexity
   - Cache invalidation challenges

2. **Flexibility**
   - Less flexible than GraphQL for varying data needs
   - Over/under-fetching possibilities
   - Schema evolution complexity

## API Design Patterns

### URL Structure
```
/api/v1/[resource]/[identifier]/[sub-resource]
```

### HTTP Methods
- GET: Retrieve resources
- POST: Create resources
- PUT: Update resources (full)
- PATCH: Update resources (partial)
- DELETE: Remove resources

### Response Format
```json
{
  "data": {
    // Resource data
  },
  "meta": {
    "page": 1,
    "total": 100
  },
  "error": null
}
```

### Error Format
```json
{
  "data": null,
  "meta": {},
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

## Implementation Guidelines

1. **Resource Naming**
   - Use plural nouns for collections
   - Use concrete names over abstract concepts
   - Keep URLs lowercase
   - Use hyphens for multi-word resources

2. **Query Parameters**
   - Use for filtering: `?status=active`
   - Use for sorting: `?sort=name`
   - Use for pagination: `?page=1&limit=10`
   - Use for searching: `?q=search_term`

3. **Status Codes**
   - 200: Success
   - 201: Created
   - 400: Bad Request
   - 401: Unauthorized
   - 403: Forbidden
   - 404: Not Found
   - 422: Unprocessable Entity
   - 429: Too Many Requests
   - 500: Server Error

4. **Security**
   - JWT in Authorization header
   - CORS configuration
   - Rate limiting
   - Input validation
   - Output sanitization

## Example Endpoints

### Patients Resource
```
GET    /api/v1/patients
POST   /api/v1/patients
GET    /api/v1/patients/:id
PUT    /api/v1/patients/:id
DELETE /api/v1/patients/:id
GET    /api/v1/patients/:id/appointments
POST   /api/v1/patients/:id/appointments
```

## References

- [REST API Design Best Practices](https://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api)
- [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines)
- [JSON:API Specification](https://jsonapi.org/)
- [REST API Security Checklist](https://github.com/shieldfy/API-Security-Checklist) 