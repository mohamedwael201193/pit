package ledger

import "testing"

func TestRecoverNeverBlindRepost(t *testing.T) {
	ok, err := Recover(Record{Status: StatusPreviewed}, "", false)
	if err != nil || !ok {
		t.Fatal(err)
	}
	ok, err = Recover(Record{Status: StatusSigned}, "123", true)
	if err != nil || ok {
		t.Fatal("signed with oid must not repost")
	}
	ok, err = Recover(Record{Status: StatusSigned}, "", false)
	if err == nil || ok {
		t.Fatal("signed without exchange view must query, not repost")
	}
	ok, err = Recover(Record{Status: StatusTimeout}, "", false)
	if err == nil || ok {
		t.Fatal("timeout")
	}
	ok, err = Recover(Record{Status: StatusAuthorized}, "", false)
	if err == nil || ok {
		t.Fatal("authorized without exchange view must not repost")
	}
}
