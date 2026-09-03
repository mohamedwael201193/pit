package deskid

import (
	"os"
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestSnapshotLiveAristotle(t *testing.T) {
	if os.Getenv("PIT_LIVE_IDENTITY") != "1" {
		t.Skip("set PIT_LIVE_IDENTITY=1")
	}
	ch := config.For(config.Mainnet)
	live, err := Snapshot(ch, MainnetTokenID, RecordedAgent)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(live.Owner, "0xbdfcee82bd42fefa58ee850b3709636a8b6b0034") {
		t.Fatalf("owner %s", live.Owner)
	}
	if !live.OwnerAuthorized {
		t.Fatal("owner must be authorized")
	}
	if !live.Supports7857 || !live.SupportsAuthorize || !live.Supports721 {
		t.Fatalf("ids %+v", live)
	}
}

func TestRequireOwnerMismatch(t *testing.T) {
	if os.Getenv("PIT_LIVE_IDENTITY") != "1" {
		t.Skip("set PIT_LIVE_IDENTITY=1")
	}
	ch := config.For(config.Mainnet)
	if err := RequireOwner(ch, MainnetTokenID, RecordedAgent); err == nil {
		t.Fatal("session agent is not owner")
	}
}
