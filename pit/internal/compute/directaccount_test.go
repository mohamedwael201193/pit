package compute

import (
	"math/big"
	"os"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
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

func TestApplyDecodedAccountSlice(t *testing.T) {
	var out AccountProbe
	applyDecodedAccount([]any{
		"user", "prov", big.NewInt(1), big.NewInt(5_000_000_000_000_000_000), big.NewInt(0), nil, "", true, big.NewInt(0), uint64(1), big.NewInt(0),
	}, &out)
	if !out.Present || !out.Acknowledged || out.BalanceOG() != "5" {
		t.Fatalf("%+v %s", out, out.BalanceOG())
	}
}

func TestLiveDirectAccount(t *testing.T) {
	if os.Getenv("PIT_LIVE_LEDGER") != "1" {
		t.Skip("set PIT_LIVE_LEDGER=1")
	}
	p := ProbeDirectAccount(config.MainnetChain(), "0xbdfcee82bd42fefa58ee850b3709636a8b6b0034", "0x7DCFe6AEa70350C2090041524c9B4A9262DCe87D")
	t.Logf("present=%v ack=%v og=%s err=%s gen=%d", p.Present, p.Acknowledged, p.BalanceOG(), p.Err, p.Generation)
	if !p.Present {
		t.Fatalf("unread %s", p.Err)
	}
}
