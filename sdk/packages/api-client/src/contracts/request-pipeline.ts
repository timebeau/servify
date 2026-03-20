export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

export interface ApiRequestContext {
  method: HttpMethod;
  path: string;
  headers: Record<string, string>;
  query?: Record<string, string>;
  body?: unknown;
  idempotencyKey?: string;
  metadata?: Record<string, unknown>;
}

export interface ApiResponseContext<TBody = unknown> {
  status: number;
  headers: Record<string, string>;
  body?: TBody;
}

export interface IdempotencyKeyProvider {
  generate(context: ApiRequestContext): string;
}

export interface ApiRequestMiddleware {
  name: string;
  onRequest?(context: ApiRequestContext): Promise<ApiRequestContext> | ApiRequestContext;
  onResponse?<TBody>(
    context: ApiRequestContext,
    response: ApiResponseContext<TBody>,
  ): Promise<ApiResponseContext<TBody>> | ApiResponseContext<TBody>;
}
