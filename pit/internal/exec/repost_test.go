package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/ledger"
)

func TestAfterTimeoutNeverRepostsSigned(t *testing.T) {
	err := AfterTimeout(ledger.Record{Status: ledger.StatusSigned, Cloid: "c1"}, "", false)
	if err == nil {
		t.Fatal("signed")
	}
	err = AfterTimeout(ledger.Record{Status: ledger.StatusTimeout}, "oid-1", true)
	if err == nil {
		t.Fatal("known oid")
	}
	if err := AfterTimeout(ledger.Record{Status: ledger.StatusPreviewed}, "", false); err != nil {
		t.Fatal(err)
	}
}
