package session

func RefuseExportJSON(k AgentKey) error {
	return k.ExportJSON()
}
