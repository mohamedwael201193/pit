package phase

func WebMaySign(surface string) bool {
	switch surface {
	case "web", "browser", "vercel":
		return false
	default:
		return surface == "desktop" || surface == "cli"
	}
}
