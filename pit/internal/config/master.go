package config

func RefuseMasterAsUser(addr string) error {
	return RejectGlobalUser(addr)
}
