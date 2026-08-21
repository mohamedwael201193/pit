package exec

import "testing"

func TestQueryBeforeRetryNeverBlind(t *testing.T) {
	if err := QueryBeforeRetry(false, "1", ""); err == nil || err.Error() != "query_exchange_first" {
		t.Fatalf("%v", err)
	}
	if err := QueryBeforeRetry(true, "oid-1", "oid-1"); err == nil || err.Error() != "already_posted" {
		t.Fatalf("%v", err)
	}
	if err := QueryBeforeRetry(true, "oid-1", ""); err == nil {
		t.Fatal("repost")
	}
}
