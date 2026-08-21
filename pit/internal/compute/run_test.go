package compute

import "testing"

func TestMustNativeSealer(t *testing.T) {
	if err := MustNativeSealer(""); err == nil {
		t.Fatal("empty")
	}
	if err := MustNativeSealer("committee.py"); err == nil {
		t.Fatal("python")
	}
	if err := MustNativeSealer("sealer.ts"); err == nil {
		t.Fatal("ts")
	}
	if err := MustNativeSealer("/usr/local/bin/pit-sealer"); err != nil {
		t.Fatal(err)
	}
}

func TestRunSealedAskStopsWithoutExecTheater(t *testing.T) {
	err := RunSealedAsk(DirectJob{
		Bin:           "/usr/local/bin/pit-sealer",
		AuthPath:      "/tmp/a",
		PromptPath:    "/tmp/p",
		OutPath:       "/tmp/o",
		Role:          Researcher,
		ProviderURL:   "https://compute-network-19.integratenetwork.work",
		OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err == nil || err.Error() != "sealer_exec_not_attached" {
		t.Fatalf("%v", err)
	}
}
