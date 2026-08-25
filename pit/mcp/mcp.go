package mcp

import (
	"encoding/json"
	"fmt"
)

var AllowedTools = []string{
	"market", "opportunities", "forecast", "status", "card", "verify",
}

var ForbiddenTools = []string{
	"authorize", "order", "cancel", "transfer", "withdraw", "export_session", "key", "sealer", "auth_file",
}

func IsAllowed(name string) bool {
	for _, t := range AllowedTools {
		if t == name {
			return true
		}
	}
	return false
}

type Request struct {
	Tool string          `json:"tool"`
	Args json.RawMessage `json:"args"`
}

type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Body  any    `json:"body,omitempty"`
}

func Handle(req Request) Response {
	if !IsAllowed(req.Tool) {
		return Response{OK: false, Error: "mcp_read_only"}
	}
	switch req.Tool {
	case "status":
		return Response{OK: true, Body: map[string]any{"surface": "mcp", "authorize": false, "sign": false}}
	case "opportunities":
		return Response{OK: true, Body: map[string]any{
			"count": 0,
			"copy":  "No opportunities match your policy.",
			"trade": false,
		}}
	case "market":
		return Response{OK: true, Body: map[string]any{
			"note":   "bind a workspace and a live venue quote with a timestamp",
			"source": "hyperliquid",
		}}
	case "forecast":
		return Response{OK: true, Body: map[string]any{"note": "forecasts are host-scored; model size is ignored"}}
	case "card":
		return EmptyCard()
	case "verify":
		return Response{OK: true, Body: VerifyHint()}
	default:
		return Response{OK: true, Body: map[string]any{"tool": req.Tool, "note": "bind a workspace first"}}
	}
}

func Schema() string {
	return fmt.Sprintf("allowed=%v forbidden=%v", AllowedTools, ForbiddenTools)
}
