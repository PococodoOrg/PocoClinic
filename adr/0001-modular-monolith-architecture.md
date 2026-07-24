# ADR-0001: Modular Monolith Architecture

## Status
Accepted

## Context
When designing the PocoClinic EMR system, we needed to choose an architecture that would:
- Be simple to deploy and maintain for non-profit organizations with limited technical resources
- Allow for future scalability and potential migration to microservices if needed
- Maintain clear boundaries between different domains of the application
- Enable rapid development and easy testing
- Ensure data consistency across the system
- Keep operational complexity low

## Decision
We have decided to implement a modular monolith architecture using vertical slices and the mediator pattern. The architecture will have the following characteristics:

1. **Vertical Slices**:
   - Each business capability (e.g., patient demographics, appointments, prescriptions) will be a separate module
   - Modules will have their own domain models, services, and handlers
   - Cross-cutting concerns will be handled through shared infrastructure code

2. **Mediator Pattern**:
   - Commands and queries will be handled through a central mediator
   - Each request will be encapsulated in a command/query object
   - Handlers will be specific to each command/query
   - This enables loose coupling between modules while maintaining clear request/response flows

3. **Module Structure**:
   ```
   internal/
   ├── features/
   │   ├── patients/
   │   │   ├── commands/
   │   │   ├── queries/
   │   │   ├── domain/
   │   │   └── handlers/
   │   └── ...
   ├── shared/
   │   ├── auth/
   │   ├── validation/
   │   └── infrastructure/
   └── api/
       └── routes/
   ```

4. **Communication Pattern**:
   - Inter-module communication will happen through the mediator
   - Direct dependencies between modules are prohibited
   - Shared types and interfaces will be placed in the shared domain package

## Consequences

### Positive
- Simpler deployment and operations compared to microservices
- Clear module boundaries while maintaining monolithic benefits
- Easier testing and debugging
- Strong consistency guarantees
- Lower initial complexity
- Easier to refactor and maintain
- Clear upgrade path to microservices if needed in the future

### Negative
- Need to be disciplined about maintaining module boundaries
- Risk of modules becoming too coupled over time
- All modules must use the same technology stack
- Scaling must be done for the entire application
- Need to carefully manage the size of the codebase

### Mitigations
- Regular architectural reviews to ensure module boundaries are maintained
- Clear documentation and guidelines for inter-module communication
- Automated tests to verify module independence
- Code analysis tools to detect unwanted dependencies between modules 