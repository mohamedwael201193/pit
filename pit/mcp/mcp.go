package mcp

import (
	"encoding/json"
	"fmt"
)

var AllowedTools = []string{
	"market", "opportunities", "watch", "forecast", "status", "card", "verify", "preview",
	"receipts", "calibration", "policy", "security", "research", "experience", "memory",
}

var ForbiddenTools = []string{
	"authorize", "order", "cancel", "transfer", "withdraw", "export_session", "key", "sealer", "auth_file", "post",
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

// LiveOpportunities is set by cmd/mcp to read public Watch. Tests leave it nil.
var LiveOpportunities func() map[string]any

func Handle(req Request) Response {
	if !IsAllowed(req.Tool) {
		return Response{OK: false, Error: "mcp_read_only"}
	}
	switch req.Tool {
	case "status":
		return Response{OK: true, Body: map[string]any{"surface": "mcp", "authorize": false, "sign": false, "trade": false}}
	case "watch":
		req.Tool = "opportunities"
		return Handle(req)
	case "policy":
		return Response{OK: true, Body: map[string]any{
			"copy": "Policy is pinned on the desktop. MCP cannot raise clip, leverage, or permissions.",
			"sign": false, "trade": false, "withdraw": false,
		}}
	case "security":
		return Response{OK: true, Body: map[string]any{
			"order": false, "cancel": false, "withdraw": false, "transfer": false, "leverage": false,
			"sign": false, "trade": false, "authorize": false,
			"desk_session_order_cancel": true,
			"copy":                      "MCP cannot order or cancel. Order/cancel exists only on the desktop session.",
		}}
	case "research":
		return Response{OK: true, Body: map[string]any{
			"copy":  "Research runs on PIT Desktop. MCP cannot start sealed inference or hold a Direct token.",
			"start": false, "sign": false, "trade": false,
		}}
	case "opportunities":
		body := map[string]any{
			"count": 0,
			"copy":  "No opportunities match your policy.",
			"trade": false,
			"sign":  false,
		}
		if LiveOpportunities != nil {
			if live := LiveOpportunities(); live != nil {
				body = live
			}
		}
		body["trade"] = false
		body["sign"] = false
		return Response{OK: true, Body: body}
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
	case "preview":
		return PreparePreview()
	case "receipts":
		return Response{OK: true, Body: map[string]any{"count": 0, "copy": "No receipt until Hyperliquid accepts an order after AUTHORIZE.", "sign": false, "trade": false}}
	case "calibration":
		return Response{OK: true, Body: map[string]any{"copy": "NOT ENOUGH DATA", "n": 0, "sign": false, "trade": false}}
	case "experience":
		return Response{OK: true, Body: map[string]any{
			"copy": "Verified experience lives on PIT Desktop. MCP can only say NOT ENOUGH DATA. No memory key. No authorize.",
			"n":    0, "sign": false, "trade": false, "authorize": false,
		}}
	case "memory":
		req.Tool = "experience"
		return Handle(req)
	default:
		return Response{OK: true, Body: map[string]any{"tool": req.Tool, "note": "bind a workspace first"}}
	}
}

func Schema() string {
	return fmt.Sprintf("allowed=%v forbidden=%v", AllowedTools, ForbiddenTools)
}
