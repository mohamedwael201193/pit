package storage

func RefuseNodeClient(path string) error {
	return RejectUnofficialClient(path)
}
