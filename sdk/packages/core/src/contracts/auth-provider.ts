export interface AuthContext {
  audience?: string;
  scopes?: string[];
  metadata?: Record<string, unknown>;
}

export interface AuthToken {
  accessToken: string;
  expiresAt?: string;
  refreshToken?: string;
}

export interface AuthProvider {
  getToken(context?: AuthContext): Promise<AuthToken | null>;
  refreshToken?(context?: AuthContext): Promise<AuthToken | null>;
}
