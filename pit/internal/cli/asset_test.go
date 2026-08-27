package cli

import "testing"

func TestLiveAssetNeedsCoin(t *testing.T) {
	if _, err := LiveAsset("not-a-net", "ETH"); err == nil {
		t.Fatal("network")
	}
}
