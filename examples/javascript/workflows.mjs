import {FreePastebin} from '../../sdk/javascript/index.js';
const client = new FreePastebin();
const paste = await client.create('A token-protected example', {visibility: 'private', expires_in: 300});
try {
  console.log(await client.raw(paste.id, paste.read_token));
  // Store paste.read_token and paste.delete_token securely.
} finally { await client.delete(paste.id, paste.delete_token); }
