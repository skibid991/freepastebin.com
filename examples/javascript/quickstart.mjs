import {FreePastebin} from '../../sdk/javascript/index.js';
const client = new FreePastebin();
const paste = await client.create('Hello from JavaScript', {expires_in: 300});
try { console.log(await client.raw(paste.id)); } finally { await client.delete(paste.id, paste.delete_token); }
