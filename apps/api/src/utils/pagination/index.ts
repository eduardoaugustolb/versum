const DEFAULT_PAGE = 1;
const DEFAULT_LIMIT = 10;
const MAX_LIMIT = 50;

export interface PaginationQuery {
  page: number;
  limit: number;
}

export function parsePagination(
  query: Record<string, string | undefined>,
): PaginationQuery {
  let page = Number(query.page);
  let limit = Number(query.limit);

  // valida page
  if (!Number.isInteger(page) || page <= 0) {
    page = DEFAULT_PAGE;
  }

  // valida limit
  if (!Number.isInteger(limit) || limit <= 0) {
    limit = DEFAULT_LIMIT;
  }

  // proteção contra abuso
  if (limit > MAX_LIMIT) {
    limit = MAX_LIMIT;
  }

  return { page, limit };
}
