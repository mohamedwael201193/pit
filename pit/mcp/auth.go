package mcp

func AuthFileNever() bool {
	r := Handle(Request{Tool: "auth_file"})
	return !r.OK && r.Error == "mcp_read_only"
}

func SealerNever() bool {
	r := Handle(Request{Tool: "sealer"})
	return !r.OK && r.Error == "mcp_read_only"
}
