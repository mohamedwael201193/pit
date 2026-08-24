package mcp

func StatusNeverSigns() Response {
	r := Handle(Request{Tool: "status"})
	return r
}
