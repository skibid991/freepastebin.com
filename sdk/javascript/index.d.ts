export interface Paste {id: string; content: string; visibility: 'public' | 'private'; title: string; language: string; created_at: string; expires_at: string | null}
export interface CreatedPaste extends Omit<Paste, 'content'> {url: string; raw_url: string; api_url: string; delete_token: string; read_token?: string}
export interface RequestOptions {query?: Record<string, unknown>; method?: string; body?: BodyInit | Record<string, unknown>; headers?: HeadersInit; check?: boolean; redirect?: RequestRedirect}
export class APIError extends Error {status: number; body: string}
export class FreePastebin {
  constructor(options?: {baseUrl?: string; timeout?: number; fetch?: typeof globalThis.fetch});
  url(path: string, query?: Record<string, unknown>): string;
  request(path: string, options?: RequestOptions): Promise<Response>;
  create(content: string, options?: {title?: string; language?: string; expires_in?: number; visibility?: 'public' | 'private'; method?: 'GET' | 'POST'}): Promise<CreatedPaste>;
  get(id: string, readToken?: string): Promise<Paste>;
  raw(id: string, readToken?: string): Promise<string>;
  delete(id: string, deleteToken: string): Promise<void>;
}
