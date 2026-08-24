package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestRefuseLeverage(t *testing.T) {
	p := policy.Default()
	if err := RefuseLeverage(p, 2); err == nil {
		t.Fatal("lev")
	}
	if err := RefuseLeverage(p, 1); err != nil {
		t.Fatal(err)
	}
}
