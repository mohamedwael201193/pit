package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestRefuseSlippage(t *testing.T) {
	p := policy.Default()
	if err := RefuseSlippage(p, 10); err != nil {
		t.Fatal(err)
	}
	if err := RefuseSlippage(p, p.MaxSlippageBps+1); err == nil {
		t.Fatal("slip")
	}
}
