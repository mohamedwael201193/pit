# pit-mcp

Read-only [Model Context Protocol](https://modelcontextprotocol.io) server for [PIT](https://github.com/mohamedwael201193/pit).

This process can read public health, public watch, and — if PIT Desktop is running on the same computer — loopback status, activity, positions, research status, and proof index.

It **cannot** authorize a preview, place an order, export a session, or arm a Sleep Mission.

```json
{
  "mcpServers": {
    "pit": {
      "command": "npx",
      "args": ["-y", "pit-mcp"]
    }
  }
}
```

The older Go binary `go run ./cmd/mcp` is a custom NDJSON loop, not this protocol. Use this package in Cursor. Use the Go binary if you already script against `{"tool":"..."}` lines.

Install:

```powershell
npm install -g pit-mcp
```

Package version **0.9.12**. `npx -y pit-mcp` starts the stdio server.
