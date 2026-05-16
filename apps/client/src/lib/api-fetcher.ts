type CustomRequestInit = RequestInit & { baseURL?: string };

type ErrorResponse = {
  status: number;
  data: unknown;
};

class ApiError extends Error {
  response: ErrorResponse;

  constructor(message: string, response: ErrorResponse) {
    super(message);
    this.response = response;
  }
}

const isAbsoluteUrl = (input: string) =>
  input.startsWith("http://") || input.startsWith("https://");

export default async function apiFetcher<T>(
  url: string,
  options?: CustomRequestInit,
): Promise<T> {
  const { baseURL, ...fetchOptions } = options ?? {};
  const resolvedBaseURL = baseURL ?? process.env.NEXT_PUBLIC_API_URL;
  const fullUrl =
    resolvedBaseURL && !isAbsoluteUrl(url)
      ? `${resolvedBaseURL.replace(/\/+$/, "")}${url}`
      : url;

  const response = await fetch(fullUrl, {
    credentials: "include",
    ...fetchOptions,
  });

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}));
    throw new ApiError(errorBody.message ?? `HTTP ${response.status}`, {
      status: response.status,
      data: errorBody,
    });
  }

  if (response.status === 204 || response.status === 205) {
    return undefined as T;
  }

  return response.json();
}
