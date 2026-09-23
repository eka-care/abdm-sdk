export type Environment = "sandbox" | "production";

export const ABHA_BASE_URLS: Record<Environment, string> = {
  sandbox: "https://abhasbx.abdm.gov.in/abha/api",
  production: "https://abha.abdm.gov.in/abha/api",
};

export type AccessTokenSource = string | (() => string | Promise<string>);

export interface M1Config {
  environment?: Environment;
  /** ABDM gateway session token, or a function returning a fresh one per call. */
  accessToken: AccessTokenSource;
  /** Override the ABHA host, e.g. to route through a proxy. */
  baseUrl?: string;
  fetch?: typeof fetch;
}

export interface CallOptions {
  /** Session token for this call only. */
  accessToken?: string;
  requestId?: string;
  signal?: AbortSignal;
  headers?: Record<string, string>;
}
