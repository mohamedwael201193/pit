package verify

import "testing"

func TestReceiptBind(t *testing.T) {
	h := HashPreview([]byte(`{"sz":"0"}`))
	if err := Check(Receipt{PreviewHash: h, StorageRoot: h, Network: "mainnet", Workspace: "ws"}); err != nil {
		t.Fatal(err)
	}
	if err := Check(Receipt{PreviewHash: "abc", StorageRoot: h, Network: "mainnet", Workspace: "ws"}); err == nil {
		t.Fatal("hash")
	}
}
