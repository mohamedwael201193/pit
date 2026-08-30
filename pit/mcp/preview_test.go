package mcp

import "testing"

func TestPreviewCannotAuthorize(t *testing.T) {
	r := Handle(Request{Tool: "preview"})
	if !r.OK {
		t.Fatal(r)
	}
	body, _ := r.Body.(map[string]any)
	if body["authorize"] != false || body["sign"] != false || body["trade"] != false || body["prepare"] != false {
		t.Fatal(body)
	}
	if !IsAllowed("preview") {
		t.Fatal("prepare")
	}
	if Handle(Request{Tool: "authorize"}).OK {
		t.Fatal("authorize")
	}
}

func TestLiveOpportunitiesCannotTrade(t *testing.T) {
	LiveOpportunities = func() map[string]any {
		return map[string]any{"count": 2, "trade": true, "sign": true}
	}
	defer func() { LiveOpportunities = nil }()
	r := Handle(Request{Tool: "opportunities"})
	body, _ := r.Body.(map[string]any)
	if body["trade"] != false || body["sign"] != false {
		t.Fatal(body)
	}
	if body["count"] != 2 {
		t.Fatal(body)
	}
}
