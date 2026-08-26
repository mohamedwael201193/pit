# PIT SDK (browser)

Typed helpers for web apps. `canSign` is always false. This package cannot authorize, place orders, or read a Direct auth file.

```ts
import { attention, canSign, explorer } from "@pit/sdk";

attention(0); // No opportunities match your policy.
explorer("mainnet");
canSign; // false
```
