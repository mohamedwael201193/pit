package cli

import "testing"

func TestConfirmDeskAuthorize(t *testing.T) {
	if err := ConfirmDeskAuthorize("yes", true); err == nil {
		t.Fatal("yes")
	}
	if err := ConfirmDeskAuthorize("AUTHORIZE", false); err == nil {
		t.Fatal("expired")
	}
	if err := ConfirmDeskAuthorize("AUTHORIZE", true); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteDeskOrderNeedsPreview(t *testing.T) {
	got := ExecuteDeskOrder(t.TempDir(), "AUTHORIZE", "")
	if got.OK || got.Error != "unbound" {
		t.Fatalf("%+v", got)
	}
}

func TestReconcileLastOrderDoesNotInventFill(t *testing.T) {
	dir := t.TempDir()
	saveLastOrder(dir, map[string]any{
		"ok": true, "posted": true, "oid": "1", "status": "submitted", "lifecycle": "submitted",
		"sign": false, "trade": false,
	})
	ReconcileLastOrder(dir)
	last := LoadLastOrder(dir)
	if last["lifecycle"] == "filled" || last["status"] == "filled" {
		t.Fatalf("invented fill %+v", last)
	}
}
