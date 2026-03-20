export interface ServerAuthContext {
  audience?: string;
  scopes?: string[];
  tenantId?: string;
  metadata?: Record<string, unknown>;
}

export interface ServerAccessToken {
  token: string;
  expiresAt?: string;
  tokenType?: 'Bearer' | 'Basic' | 'ApiKey';
}

export interface ServerAuthHeaders {
  authorization?: string;
  apiKey?: string;
  headers?: Record<string, string>;
}

export interface ServerAuthProvider {
  getHeaders(context?: ServerAuthContext): Promise<ServerAuthHeaders>;
  invalidate?(context?: ServerAuthContext): Promise<void>;
}

export interface ApiKeyAuthProvider extends ServerAuthProvider {
  readonly kind: 'api_key';
  readonly keyName: string;
}

export interface BearerTokenAuthProvider extends ServerAuthProvider {
  readonly kind: 'bearer';
  getAccessToken(context?: ServerAuthContext): Promise<ServerAccessToken | null>;
}
