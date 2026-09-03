package main

const agentCardJSON = `{
  "type": "https://eips.ethereum.org/EIPS/eip-8004#registration-v1",
  "name": "PIT-4bbee556",
  "description": "Private intelligence desk on 0G. Direct TeeML research, host VerifyE2EE, policy-bounded Hyperliquid execution on this computer.",
  "image": "https://pit0g.vercel.app/mark.svg",
  "services": [
    {"name": "web", "endpoint": "https://pit0g.vercel.app/"},
    {"name": "MCP", "endpoint": "https://www.npmjs.com/package/pit-mcp", "version": "0.9.12"}
  ],
  "x402Support": false,
  "active": true,
  "registrations": [
    {
      "agentId": 3489333,
      "agentRegistry": "eip155:16661:0x8004A169FB4a3325136EB29fA0ceB6D2e539a432"
    }
  ],
  "supportedTrust": ["reputation", "tee"]
}
`