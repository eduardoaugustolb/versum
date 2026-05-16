type CustomRequestInit = RequestInit & { baseURL?: string }

export default async function apiFetcher<T>(
  url: string,
  options?: CustomRequestInit,
): Promise<T> {
  const { baseURL = process.env.NEXT_PUBLIC_API_URL, ...fetchOptions } = options ?? {}
  const fullUrl = baseURL
    ? `${baseURL.replace(/\/+$/, "")}${url}`
    : url

  const response = await fetch(fullUrl, {
    credentials: 'include',
    ...fetchOptions,
  } as RequestInit)

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}))
    const error = new Error(errorBody.message ?? `HTTP ${response.status}`)
    ;(error as any).response = { status: response.status, data: errorBody }
    throw error
  }

  if (response.status === 204 || response.status === 205) {
    return undefined as T
  }

  return response.json()
}
