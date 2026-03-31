import { useMemo } from 'react';

export function navigateTo(path: string) {
  window.location.assign(path);
}

export function goBack() {
  window.history.back();
}

export function useDetailParams(): { id?: string } {
  return useMemo(() => {
    const segments = window.location.pathname.split('/').filter(Boolean);
    const id = segments.at(-1);
    return id ? { id } : {};
  }, []);
}

/** 读取 URL query string 中的指定参数 */
export function useQueryParam(key: string): string | null {
  return useMemo(() => {
    const params = new URLSearchParams(window.location.search);
    return params.get(key);
  }, []);
}
