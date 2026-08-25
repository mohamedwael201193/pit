package session

import "testing"

func TestRefuseExportJSON(t *testing.T) {
	if err := RefuseExportJSON(AgentKey{}); err == nil {
		t.Fatal("export")
	}
}
