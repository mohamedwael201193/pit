package mcp

func TransferNever() bool {
	r := Handle(Request{Tool: "transfer"})
	return !r.OK && r.Error == "mcp_read_only"
}
