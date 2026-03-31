import { clearToken, getToken } from '@/utils/auth';

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  params?: object;
  data?: unknown;
  headers?: Record<string, string>;
  responseType?: 'json' | 'blob' | 'text';
};

function buildQueryString(params?: object) {
  if (!params) {
    return '';
  }

  const search = new URLSearchParams();
  Object.entries(params as Record<string, unknown>).forEach(([key, value]) => {
    if (
      value !== undefined &&
      value !== null &&
      value !== '' &&
      (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean')
    ) {
      search.set(key, String(value));
    }
  });

  const query = search.toString();
  return query ? `?${query}` : '';
}

async function parseErrorMessage(response: Response) {
  try {
    const payload = await response.json();
    if (typeof payload?.message === 'string' && payload.message) {
      return payload.message;
    }
    if (typeof payload?.error === 'string' && payload.error) {
      return payload.error;
    }
  } catch {
    // ignore
  }

  return `${response.status} ${response.statusText}`.trim();
}

export async function request<T = any>(url: string, options: RequestOptions = {}): Promise<T> {
  const {
    method = 'GET',
    params,
    data,
    headers,
    responseType = 'json',
  } = options;
  const token = getToken();
  const response = await fetch(`${url}${buildQueryString(params)}`, {
    method,
    headers: {
      ...(data !== undefined ? { 'Content-Type': 'application/json' } : {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
    body: data !== undefined ? JSON.stringify(data) : undefined,
  });

  if (response.status === 401) {
    clearToken();
    window.location.replace('/login');
    throw new Error('登录已失效，请重新登录');
  }

  if (response.status === 403) {
    throw new Error('权限不足，无法执行此操作');
  }

  if (response.status >= 500) {
    throw new Error('服务器错误，请稍后重试');
  }

  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  if (responseType === 'blob') {
    return (await response.blob()) as T;
  }

  if (responseType === 'text') {
    return (await response.text()) as T;
  }

  return (await response.json()) as T;
}
