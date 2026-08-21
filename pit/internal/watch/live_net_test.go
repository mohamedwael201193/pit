package watch

import (
	"os"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestLiveVenueBooksIfReachable(t *testing.T) {
	if os.Getenv("PIT_LIVE_MARKET") != "1" {
		t.Skip("set PIT_LIVE_MARKET=1 to hit the venue")
	}
	cands, err := Live(hl.New(config.MainnetChain()), policy.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) == 0 {
		t.Log(Attention(0))
	}
	if err := MayPlaceOrder(true); err == nil {
		t.Fatal("trade")
	}
}
