package compute

import (
	"math/big"
	"testing"
)

func TestEnoughForCommittee(t *testing.T) {
	low := AccountProbe{Acknowledged: true, BalanceWei: big.NewInt(0).SetUint64(1_010_430_000_000_000_000)}
	if low.EnoughForCommittee() {
		t.Fatal("1.01 OG must not pass the three-role floor")
	}
	ok := AccountProbe{Acknowledged: true, BalanceWei: new(big.Int).Set(CommitteeFloorWei)}
	if !ok.EnoughForCommittee() {
		t.Fatal("floor")
	}
	if (AccountProbe{Acknowledged: false, BalanceWei: new(big.Int).Set(CommitteeFloorWei)}).EnoughForCommittee() {
		t.Fatal("ack")
	}
	ether := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	reclaiming := AccountProbe{
		Acknowledged:     true,
		BalanceWei:       new(big.Int).Mul(big.NewInt(9), ether),
		PendingRefundWei: new(big.Int).Mul(big.NewInt(7), ether),
	}
	if reclaiming.EnoughForCommittee() {
		t.Fatal("reclaiming ledger must not count as locked committee credit")
	}
	if reclaiming.LockedOG() != "2" {
		t.Fatalf("locked %s", reclaiming.LockedOG())
	}
}

func TestSponsorQuotaIsolatesWorkspaces(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < sponsorQuotaPerDay; i++ {
		if err := ConsumeSponsorQuota(dir, "ws-a"); err != nil {
			t.Fatal(err)
		}
	}
	if err := ConsumeSponsorQuota(dir, "ws-a"); err == nil {
		t.Fatal("cap")
	}
	if err := ConsumeSponsorQuota(dir, "ws-b"); err != nil {
		t.Fatal(err)
	}
}

func TestBalanceOG(t *testing.T) {
	p := AccountProbe{BalanceWei: CommitteeFloorWei}
	if p.BalanceOG() != "3" {
		t.Fatal(p.BalanceOG())
	}
}

func TestFileKeyringSkipsSponsorDiscovery(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	t.Setenv("PIT_DIRECT_SPONSOR_FILE", "")
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	if DirectSponsorPath() != "" {
		t.Fatal("tests must not pick up a machine sponsor file")
	}
}
