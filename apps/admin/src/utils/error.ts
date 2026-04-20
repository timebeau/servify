export const getErrorMessage = (error: unknown, fallback: string): string => {
  if (typeof error === 'object' && error !== null && 'message' in error) {
    const messageValue = (error as { message?: unknown }).message;
    if (typeof messageValue === 'string' && messageValue.trim()) {
      return messageValue;
    }
  }
  return fallback;
};

export const isFormValidationError = (error: unknown): error is { errorFields: unknown[] } =>
  typeof error === 'object' && error !== null && 'errorFields' in error;
