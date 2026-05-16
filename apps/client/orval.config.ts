import { defineConfig } from "orval";

const openApiUrl = new URL(process.env.NEXT_PUBLIC_API_URL || "http://localhost:4002/");
openApiUrl.pathname = "/openapi.json"

export default defineConfig({
  tanstackQuery: {
    input: openApiUrl.toString(),
    output: {
      mode: 'tags-split',
      target: './src/dal/orval/tanstackQuery/',
      schemas: './src/dal/orval/tanstackQuery/schemas',
      client: 'react-query',
    },
  },
  fetch: {
    input: openApiUrl.toString(),
    output: {
      mode: 'tags-split',
      target: './src/dal/orval/fetch/',
      schemas: './src/dal/orval/fetch/schemas',
      client: 'fetch',
      override: {
        mutator: {
          path: './src/lib/api-fetcher.ts',
          default: true,
        },
        fetch: {
          includeHttpResponseReturnType: false,
        },
      },
    },
  },
  zod: {
    input: openApiUrl.toString(),
    output: {
      mode: 'tags-split',
      target: './src/dal/orval/zod/',
      schemas: './src/dal/orval/zod/schemas',
      client: 'zod',
    },
  },
})
