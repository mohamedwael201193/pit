package compute

import (
	"math/big"
	"testing"
)

func TestAsBigInt(t *testing.T) {
	if asBigInt(nil).Sign() != 0 {
		t.Fatal("nil")
	}
	n := big.NewInt(42)
	if asBigInt(n).Cmp(n) != 0 {
		t.Fatal("ptr")
	}
}
