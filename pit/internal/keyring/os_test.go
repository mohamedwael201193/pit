package keyring

import (
	"testing"

	osk "github.com/zalando/go-keyring"
)

func TestOSStoreRoundTrip(t *testing.T) {
	osk.MockInit()
	s := OSStore{}
	if err := s.Put("ws-a", "session", []byte("secret-a")); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("ws-b", "session"); err == nil {
		t.Fatal("B must not read A")
	}
	got, err := s.Get("ws-a", "session")
	if err != nil || string(got) != "secret-a" {
		t.Fatalf("%s %v", got, err)
	}
	if err := s.Delete("ws-a", "session"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("ws-a", "session"); err == nil {
		t.Fatal("deleted")
	}
}

func TestOpenProductFileInTests(t *testing.T) {
	s, err := OpenProduct(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.(*FileStore); !ok {
		t.Fatalf("%T", s)
	}
}
