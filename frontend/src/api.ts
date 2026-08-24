import { client } from './api/generated/client.gen';
import type { ErrorModel } from './api/generated/types.gen';

export const apiBaseUrl =
  import.meta.env.VITE_API_BASE_URL ?? (import.meta.env.DEV ? 'http://localhost:8080' : window.location.origin);

client.setConfig({
  baseUrl: apiBaseUrl,
  credentials: 'include',
});

export function readApiError(error: unknown): string {
  if (error && typeof error === 'object') {
    const model = error as ErrorModel;
    return model.detail || model.title || 'Request failed';
  }
  if (error instanceof Error) {
    return error.message;
  }
  return 'Request failed';
}

export function unwrap<T>(result: { data?: T; error?: unknown }): T {
  if (result.error) {
    throw new Error(readApiError(result.error));
  }
  if (result.data === undefined) {
    throw new Error('Empty API response');
  }
  return result.data;
}
