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
