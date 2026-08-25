package mcp

func ExportNever() bool {
	r := Handle(Request{Tool: "export_session"})
	return !r.OK && r.Error == "mcp_read_only"
}
