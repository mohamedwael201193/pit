package memory

import "testing"

func TestRefuseGlobalKey(t *testing.T) {
	t.Setenv("PIT_PRODUCT_MODE", "true")
	t.Setenv("PIT_MEMORY_KEY", "0xab")
	if err := RefuseGlobalKey(); err == nil {
		t.Fatal("global")
	}
	t.Setenv("PIT_MEMORY_KEY", "")
	if err := RefuseGlobalKey(); err != nil {
		t.Fatal(err)
	}
}
