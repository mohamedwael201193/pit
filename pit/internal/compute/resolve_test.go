package compute

import "testing"

func TestPickDirectModelPrefersTeeML53(t *testing.T) {
	live := LiveService{Model: "glm-5.2", Verifiability: "TeeML", Present: true, TeeAck: true}
	got, err := PickDirectModel(live, []ProviderModel{
		{ID: "glm-5.2", Verifiability: "TeeML"},
		{ID: "glm-5.3", Verifiability: "TeeML"},
	})
	if err != nil || got != "glm-5.3" {
		t.Fatalf("%s %v", got, err)
	}
}

func TestPickDirectModelRefusesRouterTeeTLS53(t *testing.T) {
	live := LiveService{Model: "glm-5.2", Verifiability: "TeeML", Present: true, TeeAck: true}
	got, err := PickDirectModel(live, []ProviderModel{
		{ID: "glm-5.2", Verifiability: "TeeML"},
		{ID: "glm-5.3", Verifiability: "TeeTLS"},
	})
	if err != nil || got != "glm-5.2" {
		t.Fatalf("teeTLS 5.3 must not win: %s %v", got, err)
	}
}

func TestPickDirectModelStopsWhenTeeMLGone(t *testing.T) {
	live := LiveService{Model: "glm-5.2", Verifiability: "TeeML", Present: true, TeeAck: true}
	_, err := PickDirectModel(live, []ProviderModel{
		{ID: "glm-5.3", Verifiability: "TeeTLS"},
	})
	if err == nil {
		t.Fatal("must stop when this provider no longer lists a TeeML chat SKU")
	}
}

func TestPickDirectModelUsesOnChainWhenModelsEmpty(t *testing.T) {
	live := LiveService{Model: "glm-5.2", Verifiability: "TeeML", Present: true, TeeAck: true}
	got, err := PickDirectModel(live, nil)
	if err != nil || got != "glm-5.2" {
		t.Fatalf("%s %v", got, err)
	}
}
