package exec

import "testing"

func TestRefuseTransfer(t *testing.T) {
	if err := RefuseTransfer("sendAsset"); err == nil {
		t.Fatal("xfer")
	}
	if err := RefuseTransfer("usdClassTransfer"); err == nil {
		t.Fatal("class")
	}
	if err := RefuseTransfer("cancel"); err != nil {
		t.Fatal(err)
	}
}
