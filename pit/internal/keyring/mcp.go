package keyring

func RefuseMCP() error {
	return RefuseWeb("mcp")
}
