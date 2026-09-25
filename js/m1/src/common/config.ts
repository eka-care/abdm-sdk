export type Environment = "sandbox" | "production";

/** NHA serves M1 from three API families, each with its own base URL in production. */
export interface BaseUrls {
  /** ABHA number APIs: /enrollment, /profile. */
  abha: string;
  /** PHR app APIs: ABHA address from a mobile number. */
  phrApp: string;
  /** PHR web APIs: ABHA address verification. */
  phrWeb: string;
}

export const BASE_URLS: Record<Environment, BaseUrls> = {
  sandbox: {
    abha: "https://abhasbx.abdm.gov.in/abha/api/v3",
    phrApp: "https://abhasbx.abdm.gov.in/abha/api/v3/phr/app",
    phrWeb: "https://abhasbx.abdm.gov.in/abha/api/v3/phr/web",
  },
  production: {
    abha: "https://abha.abdm.gov.in/api/abha/v3",
    phrApp: "https://phr.abdm.gov.in/api/phr/app/v3",
    phrWeb: "https://phr.abdm.gov.in/api/phr/web/v3",
  },
};

export type AccessTokenSource = string | (() => string | Promise<string>);

export interface M1Config {
  environment?: Environment;
  /** ABDM gateway session token, or a function returning a fresh one per call. */
  accessToken: AccessTokenSource;
  /** Override any base URL, e.g. to route through a proxy. */
  baseUrls?: Partial<BaseUrls>;
  /** Override NHA's public keys (PEM) if they rotate. */
  publicKeys?: { abha?: string; phr?: string };
  fetch?: typeof fetch;
}

export interface CallOptions {
  /** Session token for this call only. */
  accessToken?: string;
  requestId?: string;
  signal?: AbortSignal;
  headers?: Record<string, string>;
}
