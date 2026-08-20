package compute

import "testing"

func TestDirectRefusesRouterAndPython(t *testing.T) {
	_, err := PrepareDirect(DirectJob{
		Bin: "committee.py", Role: Researcher, ProviderURL: "https://compute-network-19.integratenetwork.work",
		AuthPath: "a.json", PromptPath: "p.txt", OutPath: "o.json", OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err == nil {
		t.Fatal("python")
	}
	_, err = PrepareDirect(DirectJob{
		Bin: "committee", Role: Researcher, ProviderURL: "https://router-api.0g.ai/v1",
		AuthPath: "a.json", PromptPath: "p.txt", OutPath: "o.json", OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err == nil {
		t.Fatal("router")
	}
	args, err := PrepareDirect(DirectJob{
		Bin: "committee", Role: Challenger, ProviderURL: "https://compute-network-19.integratenetwork.work",
		AuthPath: "a.json", PromptPath: "p.txt", OutPath: "o.json", OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err != nil || len(args) == 0 {
		t.Fatal(err)
	}
}
