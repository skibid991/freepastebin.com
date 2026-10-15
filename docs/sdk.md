# freepastebin.com SDK reference

All clients use only the public HTTP API. No API key is needed.

| JavaScript | Python | Go | Behavior |
| --- | --- | --- | --- |
| `create(content, options)` | `create(content, **options)` | `Create(ctx, content, params, method)` | POST by default. Set method GET for query creation. Options include visibility, title, language and expires_in. |
| `get(id, readToken)` | `get(id, read_token=None)` | `Get(ctx, id, readToken)` | Read JSON metadata/content. For private pastes, pass the read token. |
| `raw(id, readToken)` | `raw(id, read_token=None)` | `Raw(ctx, id, readToken)` | Read original UTF-8 text. |
| `delete(id, deleteToken)` | `delete(id, delete_token)` | `Delete(ctx, id, deleteToken)` | POST deletion with the separate deletion token; no response body. |

## Transport, errors and binary data

JavaScript uses native `Response`, including `.text()`, `.json()`, `.arrayBuffer()`, `.status` and `.headers`. Python returns `Response` with `.body` bytes, `.text()`, `.json()`, `.status` and `.headers`. Go returns `*Response` with `Body []byte`, `Text()`, `JSON(&value)`, `StatusCode` and `Header`.

Low-level request methods accept an API path, query parameters, method, headers and body. JavaScript/Python encode object bodies as JSON; Go uses `Options.Body` for JSON and `Options.RawBody` for bytes. Redirects remain visible by default. Non-2xx statuses raise `APIError` unless disabled through JavaScript/Python `check: false` / `check=False` or Go `AllowHTTPError: true`. Methods intended to return arbitrary HTTP statuses disable that check automatically. Network errors and timeouts remain ordinary runtime errors.

Clients do not perform automatic retries. A timeout does not prove a creation request failed. Binary responses must be saved as bytes. These convenience clients buffer responses; use your language's native HTTP streaming client for long-lived SSE or incremental reads.

## Custom settings

JavaScript: `new FreePastebin({timeout: 60000})`; inject a fetch implementation with `fetch: myFetch`.

Python: `FreePastebin(timeout=60)`.

Go: `client := freepastebin.New()`; configure `client.HTTPClient.Timeout` and pass a cancellable context to each method. Timeout units are milliseconds in JavaScript, seconds in Python, and `time.Duration` in Go.

For testing, override JavaScript `baseUrl`, Python `base_url`, or Go `BaseURL`. Production defaults to `https://freepastebin.com`.
