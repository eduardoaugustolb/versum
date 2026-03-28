# Naming Convention

## Overview

This project uses **snake_case for database columns** and **camelCase for API responses**.

## Database (PostgreSQL/Drizzle ORM)

- Table names: `snake_case` (e.g., `bible_books`, `user_sessions`)
- Column names: `snake_case` (e.g., `picture_url`, `created_at`, `total_chapters`)

### Drizzle Schema Example

```typescript
export const books = pgTable(
  "bible_books",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    niceName: varchar("nice_name", { length: 100 }).notNull(),
    totalChapters: integer("total_chapters").notNull(),
  },
);
```

**Key:** The Drizzle property name is `camelCase`, but the database column name (first argument) is `snake_case`.

## API Responses (OpenAPI/JSON)

All API responses MUST use **camelCase** for property names.

### Correct

```json
{
  "pictureUrl": "https://...",
  "createdAt": "2024-01-01T00:00:00Z",
  "totalChapters": 50,
  "hasNextPage": true
}
```

### Incorrect

```json
{
  "picture_url": "https://...",
  "created_at": "2024-01-01T00:00:00Z",
  "total_chapters": 50,
  "has_next_page": true
}
```

## View Models

View Models MUST use **camelCase** for all properties. They are responsible for transforming ORM data to API responses.

### Correct (pagination.view-model.ts)

```typescript
export interface PaginationInput {
  page: number;
  limit: number;
  totalItems: number;
}

export class PaginationViewModel {
  readonly currentPage: number;
  readonly totalPages: number;
  readonly totalItems: number;
  readonly itemsPerPage: number;
  readonly hasNextPage: boolean;
  readonly hasPrevPage: boolean;

  toJSON() {
    return {
      currentPage: this.currentPage,
      totalPages: this.totalPages,
      totalItems: this.totalItems,
      itemsPerPage: this.itemsPerPage,
      hasNextPage: this.hasNextPage,
      hasPrevPage: this.hasPrevPage,
    };
  }
}
```

### Incorrect

```typescript
// NEVER use snake_case in view models or API responses
readonly current_page: number;  // WRONG
readonly total_items: number;   // WRONG
```

## OpenAPI Schemas

All schema definitions in `openapi.yaml` MUST use **camelCase** for property names.

```yaml
User:
  type: object
  properties:
    id:
      type: string
      format: uuid
    pictureUrl:
      type: string
      nullable: true
    createdAt:
      type: string
      format: date-time

PaginationViewModel:
  type: object
  properties:
    currentPage:
      type: integer
    totalPages:
      type: integer
    totalItems:
      type: integer
    itemsPerPage:
      type: integer
    hasNextPage:
      type: boolean
    hasPrevPage:
      type: boolean
```

## Checklist

When creating or updating code:

- [ ] Database column names: `snake_case`
- [ ] Drizzle schema properties: `camelCase`
- [ ] View Model properties: `camelCase`
- [ ] OpenAPI schema properties: `camelCase`
- [ ] API response JSON: `camelCase`
