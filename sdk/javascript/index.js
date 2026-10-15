/** A dependency-free HTTP client for freepastebin.com. */
export class APIError extends Error {
  constructor(response, body) {
    super(`HTTP ${response.status}: ${body.slice(0, 500)}`);
    this.name = 'APIError'; this.status = response.status; this.body = body;
  }
}
export class FreePastebin {
  constructor({baseUrl = 'https://freepastebin.com', timeout = 45000, fetch: fetchImpl = globalThis.fetch} = {}) {
    this.baseUrl = baseUrl.replace(/\/$/, ''); this.timeout = timeout; this.fetch = fetchImpl;
  }
  url(path, query = {}) {
    if (!path.startsWith('/') || path.startsWith('//')) throw new Error('Use an API path beginning with one slash.');
    const url = new URL(this.baseUrl + path);
    for (const [key, value] of Object.entries(query)) {
      if (value !== undefined && value !== null) url.searchParams.set(key, typeof value === 'object' ? JSON.stringify(value) : String(value));
    }
    return url.toString();
  }
  /** Returns a standard Response. Set check=false to inspect non-2xx responses. */
  async request(path, {query = {}, method = 'GET', body, headers = {}, check = true, redirect = 'manual'} = {}) {
    const h = new Headers(headers);
    if (body !== undefined && body !== null && typeof body === 'object' && !(body instanceof Uint8Array) && !(body instanceof ArrayBuffer) && !(body instanceof Blob) && !(body instanceof FormData) && !(body instanceof URLSearchParams)) {
      body = JSON.stringify(body); if (!h.has('content-type')) h.set('content-type', 'application/json');
    }
    const response = await this.fetch(this.url(path, query), {method, body, headers: h, redirect, signal: AbortSignal.timeout(this.timeout)});
    if (check && !response.ok) throw new APIError(response, await response.text());
    return response;
  }
  async create(content, {method = 'POST', ...options} = {}) {
    const data = {...options, content}; const r = await this.request('/api/pastes', method === 'GET' ? {query: data} : {method, body: data}); return r.json();
  }
  async get(id, readToken) { return (await this.request('/api/pastes/' + encodeURIComponent(id), {headers: readToken ? {Authorization: 'Bearer ' + readToken} : {}})).json(); }
  async raw(id, readToken) { return (await this.request('/raw/' + encodeURIComponent(id), {headers: readToken ? {Authorization: 'Bearer ' + readToken} : {}})).text(); }
  async delete(id, deleteToken) { await this.request('/api/pastes/' + encodeURIComponent(id) + '/delete', {method: 'POST', body: {delete_token: deleteToken}}); }
}
