package redteam

import "testing"

func TestTwoUsersDoNotShareASession(t *testing.T) {
	a := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	b := "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	if err := TwoUsers(a, a, b); err != nil {
		t.Fatal(err)
	}
}
