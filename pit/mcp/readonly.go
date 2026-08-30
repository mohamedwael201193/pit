package mcp

func TradeDenied(tool string) bool {
	for _, t := range ForbiddenTools {
		if t == tool {
			return true
		}
	}
	switch tool {
	case "pin", "kill", "session", "direct", "mission", "flatten", "execute", "approveAgent", "guarded":
		return true
	default:
		return false
	}
}
