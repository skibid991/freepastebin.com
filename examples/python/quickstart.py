from freepastebin import FreePastebin

client = FreePastebin()
paste = client.create('Hello from Python', expires_in=300)
try:
    print(client.raw(paste['id']))
finally:
    client.delete(paste['id'], paste['delete_token'])
