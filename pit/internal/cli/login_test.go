package cli

import "testing"

func TestLoginCopy(t *testing.T) {
	if err := RefuseLoginSecret("private_key=abc"); err == nil {
		t.Fatal("secret")
	}
	if LoginCopy() == "" {
		t.Fatal("copy")
	}
}
