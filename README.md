# Free pastebin API for text and source code

[FreePastebin](https://freepastebin.com) shares UTF-8 text and code without an account or API key. Create a paste through GET or POST, open its HTML page, or retrieve the original text from a raw URL. Choose public, indexable pastes or private pastes protected by a separate read token.

## Create a paste with one GET request

```sh
curl --get 'https://freepastebin.com/api/pastes' \
  --data-urlencode 'content=Hello from a GET request' \
  --data-urlencode 'title=Example' --data-urlencode 'expires_in=300'
```

The 201 JSON response includes `id`, `url`, `raw_url`, `api_url`, `delete_token`, `visibility`, `created_at`, `expires_at`, `title` and `language`. Save the deletion token when creating the paste; it is not returned on later reads. GET creation is convenient for small, non-sensitive text. POST avoids putting the content in a URL and supports larger pastes.

## Upload a file or create a private paste

```sh
curl --data-binary @notes.txt -H 'Content-Type: text/plain; charset=utf-8' \
  'https://freepastebin.com/api/pastes'

curl 'https://freepastebin.com/api/pastes' -H 'Content-Type: application/json' \
  -d '{"content":"My private note","visibility":"private","expires_in":3600}'
```

Private creation also returns `read_token`. The returned browser share URL contains `#read_token=...` and unlocks the paste when opened. Anyone holding that full share URL can read the paste. For APIs, use the read token in the Authorization header:

```sh
curl 'https://freepastebin.com/raw/PASTE_ID' -H 'Authorization: Bearer READ_TOKEN'
curl 'https://freepastebin.com/api/pastes/PASTE_ID' -H 'Authorization: Bearer READ_TOKEN'
```

The read token is distinct from the deletion token. Private content routes require authorization and return `noindex` and `private, no-store`. A fragment is not sent to the server by a normal HTTP fetch, so API clients must explicitly provide the token. `?read_token=...` is supported for URL-only clients, but puts the token in the request URL. Private means access controlled, not end-to-end encrypted; the service operator can access submitted content.

## Read and delete

| Request | Behavior |
| --- | --- |
| `GET /{id}` | HTML paste page; public pages are eligible for search indexing |
| `GET /raw/{id}` | Raw UTF-8 content |
| `GET /api/pastes/{id}` | Content and metadata as JSON |
| `POST /api/pastes/{id}/delete` | JSON body `{"delete_token":"..."}`; success is 204 |
| `DELETE /api/pastes/{id}` | `Authorization: Bearer DELETE_TOKEN`; success is 204 |

Reads support HEAD and conditional `If-None-Match` requests. Unauthorized private reads return 401. A missing, deleted or expired paste returns 404. Deletion requires POST or DELETE; it cannot be performed through GET.

## Expiration and size

Maximum content size is 1 MiB of UTF-8 bytes. `expires_in` defaults to 2,592,000 seconds (30 days); use `0` for no expiry, or a whole number of seconds up to 315,360,000. Optional `title` is up to 200 characters and `language` up to 40. Public is the default visibility. Public pastes may be indexed when discovered; indexing is not guaranteed.

## Common questions

**Can I create a paste using only GET?** Yes. Pass `content` and optional fields as encoded query parameters. Browser and proxy URL-length limits still apply.

**Can a private paste appear in search results?** Its content requires a read token and its responses request no indexing. Keep the token and full share URL private.

## SDKs and runnable examples

This repository contains dependency-free JavaScript (with TypeScript declarations), Python and Go clients. They are source packages in this repository; no npm or PyPI release is implied. Clone the repository to run the examples:

```sh
git clone https://github.com/skibid991/freepastebin.com.git
cd freepastebin.com
```

### JavaScript and TypeScript

Requires Node.js 22+ or a browser with `fetch` and `AbortSignal.timeout`. Run `node examples/javascript/quickstart.mjs` from the repository root. You can also install the local package using `npm install ./sdk/javascript`; TypeScript declarations are included.

```js
import {FreePastebin} from './sdk/javascript/index.js';
const client = new FreePastebin();
const paste = await client.create('Hello from JavaScript', {expires_in: 300});
try { console.log(await client.raw(paste.id)); } finally { await client.delete(paste.id, paste.delete_token); }
```

### Python

Requires Python 3.10+, with no third-party packages. Run `PYTHONPATH=sdk/python python3 examples/python/quickstart.py`. To use the module in another project, copy `sdk/python/freepastebin.py` into your project or add that directory to your Python path.

```python
from freepastebin import FreePastebin

client = FreePastebin()
paste = client.create('Hello from Python', expires_in=300)
try:
    print(client.raw(paste['id']))
finally:
    client.delete(paste['id'], paste['delete_token'])
```

### Go

Requires Go 1.22+. Run `go -C examples/go run .`. The example module uses the local SDK. From another Go module, run `go get github.com/skibid991/freepastebin.com/sdk/go` and import that path.

See [the runnable Go example](examples/go/main.go) and [SDK method reference](docs/sdk.md). SDKs use a 45-second client timeout by default. Change it for your workload. Do not automatically retry creation requests: a successful request whose response is lost can still have created a resource.

## API reference and discovery

- [Interactive HTML API reference](https://freepastebin.com/docs)
- [OpenAPI 3.1 schema](https://freepastebin.com/openapi.json)
- [Plain-text usage guide](https://freepastebin.com/llms.txt)

The live API reference is authoritative. Documented GET routes work without an SDK. Encode query parameters with `URLSearchParams`, Python `urlencode`, Go `url.Values`, or cURL `--data-urlencode`.

## Development

Offline tests cover query encoding, Unicode, binary response preservation, request bodies, authorization headers, error handling and redirect behavior:

```sh
node --test sdk/javascript/client.test.mjs
PYTHONPATH=sdk/python python3 -m unittest discover -s sdk/python -p 'test_*.py'
go -C sdk/go test ./...
go -C examples/go build -o /tmp/freepastebin-example .
```

Examples call the public service. Paste/link quickstarts create temporary resources and delete them after reading. Contributions should focus on public API usage and client behavior. Licensed under [MIT](LICENSE).
