package main

import "testing"

func TestRequirePubKeyMatchesOnchain(t *testing.T) {
	ok := "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"
	if err := requirePubKeyMatchesOnchain(ok, "0xa46ea4fc5889ad35a1487e1ed04dccfa872146b9"); err != nil {
		t.Fatal(err)
	}
	if err := requirePubKeyMatchesOnchain("0x0000000000000000000000000000000000000001", ok); err == nil {
		t.Fatal("wrong signer")
	}
	if err := requirePubKeyMatchesOnchain("", ok); err == nil {
		t.Fatal("empty pub")
	}
}
