package mcp

func KeyNever() bool {
	r := Handle(Request{Tool: "key"})
	return !r.OK && r.Error == "mcp_read_only"
}
