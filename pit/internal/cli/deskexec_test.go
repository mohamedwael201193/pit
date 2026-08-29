package cli

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

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

func TestClipForAccountRefusesBelowFloor(t *testing.T) {
	_, err := ClipForAccount(policy.Default(), 9.38)
	if err == nil || err.Error() != "insufficient_margin" {
		t.Fatalf("%v", err)
	}
	got, err := ClipForAccount(policy.Default(), 40)
	if err != nil || got != 10 {
		t.Fatalf("%v %v", got, err)
	}
}

func TestExecuteGuardedNeedsEnable(t *testing.T) {
	got := ExecuteGuardedDeskOrder(t.TempDir(), "")
	if got.OK || got.Error != "need_guarded_enable" {
		t.Fatalf("%+v", got)
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
