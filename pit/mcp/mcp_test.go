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

func TestOpportunitiesDoNotTrade(t *testing.T) {
	r := Handle(Request{Tool: "opportunities"})
	if !r.OK {
		t.Fatal(r)
	}
	body, _ := r.Body.(map[string]any)
	if body["trade"] != false {
		t.Fatal(body)
	}
	if Handle(Request{Tool: "withdraw"}).OK {
		t.Fatal("withdraw")
	}
	if Handle(Request{Tool: "key"}).OK {
		t.Fatal("key")
	}
	watch := Handle(Request{Tool: "watch"})
	if !watch.OK {
		t.Fatal(watch)
	}
	wb, _ := watch.Body.(map[string]any)
	if wb["trade"] != false {
		t.Fatal(wb)
	}
	res := Handle(Request{Tool: "research"})
	if !res.OK {
		t.Fatal(res)
	}
	rb, _ := res.Body.(map[string]any)
	if rb["start"] != false || rb["trade"] != false {
		t.Fatal(rb)
	}
	if Handle(Request{Tool: "authorize"}).OK {
		t.Fatal("research must not unlock authorize")
	}
	cal := Handle(Request{Tool: "calibration"})
	if !cal.OK {
		t.Fatal(cal)
	}
	cb, _ := cal.Body.(map[string]any)
	if cb["copy"] != "NOT ENOUGH DATA" || cb["trade"] != false {
		t.Fatal(cb)
	}
}
