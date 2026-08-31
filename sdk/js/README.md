# pit-os

Read-only JavaScript helpers for [PIT](https://github.com/mohamedwael201193/pit).

This package **cannot** authorize a preview, sign a Hyperliquid order, export a session key, arm a Sleep Mission, or read a Direct auth file. Execution stays on PIT Desktop.

```ts
import {
  canSign,
  explorer,
  publicHealth,
  publicWatch,
  refuseAuthorize,
} from "pit-os";

canSign; // false
explorer("mainnet");
await publicHealth();
await publicWatch("mainnet");
// refuseAuthorize(); // throws authorize_denied
```

Public endpoints used:

- `GET https://pit-health.onrender.com/health`
- `GET https://pit-health.onrender.com/watch`
- `GET https://pit-health.onrender.com/release`

Optional loopback (same computer as PIT Desktop):

- `GET http://127.0.0.1:17373/health`
- `GET http://127.0.0.1:17373/local/status`

There is no `authorize` method.

Install:

```powershell
npm install pit-os
```
