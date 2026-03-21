export const API_VERSIONS = ["v1"] as const;

export type ApiVersion = (typeof API_VERSIONS)[number];

export const DEFAULT_API_VERSION: ApiVersion = "v1";
