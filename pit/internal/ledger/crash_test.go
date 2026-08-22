package ledger

import "testing"

func TestAfterCrash(t *testing.T) {
	if err := AfterCrash(Record{Status: StatusSigned}, false, ""); err == nil {
		t.Fatal("signed")
	}
	if err := AfterCrash(Record{Status: StatusReceipt}, true, "oid-1"); err != nil {
		t.Fatal(err)
	}
	if err := AfterCrash(Record{Status: StatusPreviewed}, false, ""); err != nil {
		t.Fatal(err)
	}
	if err := AfterCrash(Record{Status: StatusTimeout}, false, ""); err == nil {
		t.Fatal("timeout")
	}
}
