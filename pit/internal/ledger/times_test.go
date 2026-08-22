package ledger

import "testing"

func TestAfterUnknownExchange(t *testing.T) {
	if AfterUnknownExchange("signed", false) != "query_first" {
		t.Fatal("query")
	}
	if AfterUnknownExchange("previewed", true) != "previewed" {
		t.Fatal("keep")
	}
}
