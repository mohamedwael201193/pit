package mcp

func OrderNever() bool {
	r := Handle(Request{Tool: "order"})
	return !r.OK && r.Error == "mcp_read_only"
}
