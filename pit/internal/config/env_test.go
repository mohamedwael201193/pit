package config

import (
	"os"
	"testing"
)

func TestRefuseSessionEnv(t *testing.T) {
	t.Setenv("PIT_PRODUCT_MODE", "true")
	t.Setenv("PIT_SESSION_KEY", "x")
	if err := RefuseSessionEnv(); err == nil {
		t.Fatal("session")
	}
	os.Unsetenv("PIT_SESSION_KEY")
	if err := RefuseSessionEnv(); err != nil {
		t.Fatal(err)
	}
}
