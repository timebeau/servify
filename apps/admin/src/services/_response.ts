type UnknownRecord = Record<string, unknown>;

function isRecord(value: unknown): value is UnknownRecord {
  return typeof value === 'object' && value !== null;
}

function toFiniteNumber(value: unknown, fallback: number): number {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value);
    if (Number.isFinite(parsed)) {
      return parsed;
    }
  }
  return fallback;
}

function unwrapCollectionEnvelope(payload: unknown): unknown {
  let current = payload;
  while (
    isRecord(current) &&
    isRecord(current.data) &&
    !Array.isArray(current.data) &&
    ('data' in current.data || 'items' in current.data || 'entries' in current.data)
  ) {
    current = current.data;
  }
  return current;
}

function pickArray<T>(payload: unknown, keys: string[]): T[] {
  if (Array.isArray(payload)) {
    return payload as T[];
  }
  if (!isRecord(payload)) {
    return [];
  }
  for (const key of keys) {
    const value = payload[key];
    if (Array.isArray(value)) {
      return value as T[];
    }
  }
  return [];
}

export function normalizePaginatedResponse<T>(
  payload: unknown,
  itemKeys: string[] = ['data', 'items'],
): API.PaginatedResponse<T> {
  const source = unwrapCollectionEnvelope(payload);
  const data = pickArray<T>(source, itemKeys);
  const total = isRecord(source)
    ? toFiniteNumber(source.total, toFiniteNumber(source.count, data.length))
    : data.length;
  const page = isRecord(source) ? toFiniteNumber(source.page, 1) : 1;
  const pageSize = isRecord(source)
    ? toFiniteNumber(source.page_size, data.length || total)
    : data.length || total;
  const pages = isRecord(source)
    ? toFiniteNumber(
        source.pages,
        pageSize > 0 ? Math.ceil(total / pageSize) : total > 0 ? 1 : 0,
      )
    : pageSize > 0
      ? Math.ceil(total / pageSize)
      : total > 0
        ? 1
        : 0;

  return {
    data,
    total,
    page,
    page_size: pageSize,
    pages,
  };
}

export function normalizeCountedListResponse<T>(
  payload: unknown,
  itemKeys: string[] = ['items', 'data'],
): API.CountedListResponse<T> {
  const source = unwrapCollectionEnvelope(payload);
  const items = pickArray<T>(source, itemKeys);
  const count = isRecord(source)
    ? toFiniteNumber(source.count, toFiniteNumber(source.total, items.length))
    : items.length;

  return {
    count,
    items,
  };
}
