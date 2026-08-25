package config

import "testing"

func TestRefuseMasterAsUser(t *testing.T) {
	t.Setenv("PIT_PRODUCT_MODE", "true")
	t.Setenv("PIT_MASTER_ADDRESS", "0xabc")
	if err := RefuseMasterAsUser("0xabc"); err == nil {
		t.Fatal("master")
	}
	if err := RefuseMasterAsUser("0xdef"); err != nil {
		t.Fatal(err)
	}
}
