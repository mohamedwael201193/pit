package mcp

func PostNever() bool {
	r := Handle(Request{Tool: "post"})
	return !r.OK && r.Error == "mcp_read_only"
}
