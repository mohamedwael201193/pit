package compute

import "testing"

func TestHonestLabelSameProvider(t *testing.T) {
	got := Classify([]Target{
		{Provider: "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D", Model: "glm-5.2"},
		{Provider: "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D", Model: "glm-5.2"},
		{Provider: "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D", Model: "glm-5.2"},
	})
	if HonestLabel(got) != "role separation and envelope separation on the same provider" {
		t.Fatal(HonestLabel(got))
	}
}

func TestRefuseProviderSpoof(t *testing.T) {
	if err := RefuseProviderSpoof("0xa", "0xb"); err == nil {
		t.Fatal("spoof")
	}
}
