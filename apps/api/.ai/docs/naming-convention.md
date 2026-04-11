# Naming Convention (API)

## Database
- Tables: snake_case, plural (e.g., `users`, `user_profiles`)
- Columns: snake_case (e.g., `created_at`, `user_id`)
- Primary keys: `id` (UUID or bigint)

## API Endpoints
- RESTful conventions: `/resources`, `/resources/:id`
- Use kebab-case for URL paths (e.g., `/user-profiles`)
- HTTP methods: GET (read), POST (create), PUT/PATCH (update), DELETE (remove)

## Code
- Classes: PascalCase (e.g., `UserService`)
- Functions: camelCase (e.g., `getUserById`)
- Variables: camelCase
- Constants: UPPER_SNAKE_CASE
- Files: kebab-case (e.g., `user-service.ts`)

## Error Codes
- Use consistent error code prefixes by module
