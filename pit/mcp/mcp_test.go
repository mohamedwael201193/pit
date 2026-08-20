package mcp

import "testing"

func TestReadOnly(t *testing.T) {
	if Handle(Request{Tool: "authorize"}).OK {
		t.Fatal("authorize")
	}
	if Handle(Request{Tool: "order"}).OK {
		t.Fatal("order")
	}
	if Handle(Request{Tool: "export_session"}).OK {
		t.Fatal("export")
	}
	if !Handle(Request{Tool: "status"}).OK {
		t.Fatal("status")
	}
	if IsAllowed("cancel") {
		t.Fatal("cancel is not an MCP tool")
	}
}
