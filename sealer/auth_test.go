package main

import "testing"

func TestValidateAuthRefusesRouter(t *testing.T) {
	err := validateAuth(authFile{
		Provider: "0xabc", URL: "https://router-api.0g.ai/v1", Model: "glm-5.2",
		TeeSigner: "0xdef", Verifiability: "TeeML", Authorization: "Bearer app-sk-test",
	})
	if err == nil || err.Error() != "ROUTER_DOWNGRADE_DENIED" {
		t.Fatalf("%v", err)
	}
}

func TestValidateAuthRefusesRouterKey(t *testing.T) {
	err := validateAuth(authFile{
		Provider: "0xabc", URL: "https://compute-network-19.integratenetwork.work", Model: "glm-5.2",
		TeeSigner: "0xdef", Verifiability: "TeeML", Authorization: "Bearer sk-dashboard",
	})
	if err == nil || err.Error() != "router_api_key_denied" {
		t.Fatalf("%v", err)
	}
}

func TestValidateAuthRequiresTeeML(t *testing.T) {
	err := validateAuth(authFile{
		Provider: "0xabc", URL: "https://compute-network-19.integratenetwork.work", Model: "glm-5.2",
		TeeSigner: "0xdef", Verifiability: "TeeTLS", Authorization: "Bearer app-sk-test",
	})
	if err == nil || err.Error() != "NOT_TEEML" {
		t.Fatalf("%v", err)
	}
}

func TestRequireRole(t *testing.T) {
	if err := requireRole("researcher"); err != nil {
		t.Fatal(err)
	}
	if err := requireRole("challenger"); err != nil {
		t.Fatal(err)
	}
	if err := requireRole("risk"); err != nil {
		t.Fatal(err)
	}
	if err := requireRole("orchestrator"); err == nil {
		t.Fatal("role")
	}
}
