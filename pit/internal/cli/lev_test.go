package cli

import "testing"

func TestRefuseLeverageChange(t *testing.T) {
	if err := RefuseLeverageChange("updateLeverage"); err == nil {
		t.Fatal("lev")
	}
	if err := RefuseLeverageChange("order"); err != nil {
		t.Fatal(err)
	}
}
