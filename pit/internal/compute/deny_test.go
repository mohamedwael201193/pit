package compute

import "testing"

func TestDenyRouter(t *testing.T) {
	if err := DenyRouter("https://router-api.0g.ai/v1"); err == nil {
		t.Fatal("router")
	}
	if err := DenyRouter("https://router-api-testnet.integratenetwork.work/v1"); err == nil {
		t.Fatal("testnet router")
	}
	if err := DenyRouter("https://compute-network-19.integratenetwork.work"); err != nil {
		t.Fatal(err)
	}
}

func TestValidateCommittee(t *testing.T) {
	t0 := Target{
		URL: "https://compute-network-19.integratenetwork.work", Model: "glm-5.3",
		Provider: "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D",
		TeeSigner: "0x041a09E5bEF30fd776D66Bb892d18B97637C7C7c",
		Verifiability: "TeeML", Role: "researcher",
	}
	if err := Validate(t0); err != nil {
		t.Fatal(err)
	}
	bad := t0
	bad.Verifiability = "TeeTLS"
	if err := Validate(bad); err == nil {
		t.Fatal("teetls")
	}
	roles := []Target{t0, t0, t0}
	roles[1].Role = "challenger"
	roles[2].Role = "risk"
	if Classify(roles) != EnvelopeOnly {
		t.Fatal("must not claim provider diversity")
	}
}
