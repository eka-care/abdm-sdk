import { randomUUID } from "node:crypto";
import { ABHA_BASE_URLS, type AccessTokenSource, type CallOptions, type M1Config } from "./config.js";
import { AbdmError } from "./errors.js";

export interface RequestInit {
  body?: unknown;
  headers?: Record<string, string>;
  query?: Record<string, string>;
}

export interface HttpClient {
  request<T>(method: string, path: string, init?: RequestInit, opts?: CallOptions): Promise<T>;
}

export const bearer = (token: string) => (token.startsWith("Bearer ") ? token : `Bearer ${token}`);

/** Sends ABHA requests with the headers NHA requires on every call. */
export class AbhaHttpClient implements HttpClient {
  private readonly baseUrl: string;
  private readonly accessToken: AccessTokenSource;
  private readonly fetch: typeof fetch;

  constructor(config: M1Config) {
    if (!config.accessToken) throw new Error("accessToken is required");
    this.baseUrl = config.baseUrl ?? ABHA_BASE_URLS[config.environment ?? "sandbox"];
    this.accessToken = config.accessToken;
    this.fetch = config.fetch ?? globalThis.fetch;
  }

  async request<T>(method: string, path: string, init: RequestInit = {}, opts: CallOptions = {}): Promise<T> {
    const url = new URL(this.baseUrl + path);
    for (const [key, value] of Object.entries(init.query ?? {})) url.searchParams.set(key, value);

    const requestId = opts.requestId ?? randomUUID();
    const res = await this.fetch(url, {
      method,
      signal: opts.signal,
      headers: {
        "Content-Type": "application/json",
        Authorization: bearer(opts.accessToken ?? (await this.resolveToken())),
        "REQUEST-ID": requestId,
        TIMESTAMP: new Date().toISOString(),
        ...init.headers,
        ...opts.headers,
      },
      body: init.body === undefined ? undefined : JSON.stringify(init.body),
    });

    const body = await parseBody(res);
    if (!res.ok) throw new AbdmError(res.status, body, requestId);
    return body as T;
  }

  private async resolveToken(): Promise<string> {
    return typeof this.accessToken === "function" ? this.accessToken() : this.accessToken;
  }
}

async function parseBody(res: Response): Promise<unknown> {
  const type = res.headers.get("content-type") ?? "";
  if (type.includes("json")) return res.json();
  if (!type || type.startsWith("text/")) return res.text();
  return res.arrayBuffer();
}
