# Scalability Rules (API)

## Performance
- Avoid N+1 queries; use eager loading or batch queries
- Implement caching for expensive operations
- Use pagination for large datasets
- Optimize database queries with proper indexes

## Architecture
- Design with separation of concerns
- Use dependency injection to improve testability
- Implement stateless services when possible
- Design APIs to be horizontally scalable

## Database
- Use appropriate database for each use case
- Implement connection pooling
- Use read replicas for read-heavy operations
- Design schema with scaling in mind

## Async Processing
- Offload long-running tasks to background jobs
- Use message queues for async communication
- Implement proper retry mechanisms
- Handle failures gracefully with exponential backoff
