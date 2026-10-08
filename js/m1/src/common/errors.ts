export class AbdmError extends Error {
  constructor(
    readonly status: number,
    readonly body: unknown,
    readonly requestId: string,
  ) {
    super(`ABDM request failed with ${status}: ${typeof body === "string" ? body : JSON.stringify(body)}`);
    this.name = "AbdmError";
  }
}
