# ADR 007: React Query for State Management

## Status

Accepted

## Context

We needed to choose a state management solution for our React frontend that would:
- Handle server state efficiently
- Provide caching and background updates
- Support real-time data synchronization
- Minimize boilerplate code
- Integrate well with TypeScript
- Handle loading and error states

The main options considered were:
1. Redux + RTK Query
2. React Query (TanStack Query)
3. SWR
4. Apollo Client
5. Custom hooks with Context API

## Decision

We have chosen to use React Query (TanStack Query) as our primary state management solution, particularly for server state management.

## Consequences

### Advantages

1. **Server State Management**
   - Automatic background data refetching
   - Smart caching with configurable cache time
   - Optimistic updates
   - Automatic retry logic
   - Deduplication of requests

2. **Developer Experience**
   - Minimal boilerplate
   - Intuitive hooks-based API
   - Excellent TypeScript support
   - Built-in devtools
   - Easy integration with any data fetching library

3. **Performance**
   - Efficient cache management
   - Automatic garbage collection
   - Request deduplication
   - Parallel queries support
   - Lazy loading capabilities

4. **Features**
   - Infinite queries for pagination
   - Mutations with rollback
   - Prefetching capabilities
   - Offline support
   - Real-time updates support

### Disadvantages

1. **Learning Curve**
   - New patterns for developers used to traditional state management
   - Understanding cache invalidation strategies
   - Mastering advanced features like optimistic updates

2. **Client State Management**
   - Not designed for client-only state
   - Requires additional solution for UI state (using React's built-in state management)

3. **Bundle Size**
   - Adds about 12.4kB min+gzip to the bundle
   - Additional size when including devtools

## Alternatives Considered

### Redux + RTK Query
- Pros: Familiar ecosystem, good tooling, handles all state types
- Cons: More boilerplate, steeper learning curve, larger bundle size

### SWR
- Pros: Simpler API, lighter weight, good stale-while-revalidate strategy
- Cons: Fewer features, less mature ecosystem, less TypeScript-focused

### Apollo Client
- Pros: Excellent GraphQL integration, robust caching, mature ecosystem
- Cons: Overkill for REST APIs, complex setup, larger bundle size

### Custom Hooks + Context
- Pros: No dependencies, complete control, smallest bundle size
- Cons: No built-in caching, manual implementation of features, maintenance burden

## Implementation Notes

Our implementation includes:
- QueryClient configuration with default options
- Integration with Axios for data fetching
- Custom hooks for common data patterns
- Error boundary integration
- Testing utilities setup

## References

- [TanStack Query Documentation](https://tanstack.com/query/latest)
- [React Query GitHub Repository](https://github.com/TanStack/query)
- [Practical React Query](https://tkdodo.eu/blog/practical-react-query) 