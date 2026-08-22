package obs

import "testing"

func TestRefusePrivateBook(t *testing.T) {
	if err := RefusePrivateBook([]string{"requestId", "phase"}); err != nil {
		t.Fatal(err)
	}
	if err := RefusePrivateBook([]string{"book"}); err == nil {
		t.Fatal("book")
	}
}
