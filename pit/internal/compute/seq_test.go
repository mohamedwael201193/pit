package compute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCommitteeIncomplete(t *testing.T) {
	err := RunCommittee("/opt/pit/sealer", nil)
	if err == nil || err.Error() != "committee_incomplete" {
		t.Fatalf("%v", err)
	}
}

func TestRunCommitteeDuplicateRole(t *testing.T) {
	j := DirectJob{
		Bin:           "/opt/pit/sealer",
		AuthPath:      "a",
		PromptPath:    "p",
		OutPath:       "o",
		Role:          Researcher,
		ProviderURL:   "https://compute-network-19.integratenetwork.work",
		OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	}
	err := RunCommittee("/opt/pit/sealer", []DirectJob{j, j, j})
	if err == nil || err.Error() != "duplicate_role" {
		t.Fatalf("%v", err)
	}
}

func TestRunCommitteeMissingBinary(t *testing.T) {
	dir := t.TempDir()
	jobs := []DirectJob{}
	for _, role := range CommitteeRoles() {
		jobs = append(jobs, DirectJob{
			Bin:           filepath.Join(dir, "missing"),
			AuthPath:      filepath.Join(dir, "a"),
			PromptPath:    filepath.Join(dir, "p"),
			OutPath:       filepath.Join(dir, string(role)),
			Role:          role,
			ProviderURL:   "https://compute-network-19.integratenetwork.work",
			OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
		})
	}
	err := RunCommittee(filepath.Join(dir, "missing"), jobs)
	if err == nil || err.Error() != "sealer_not_wired" {
		t.Fatalf("%v", err)
	}
}

func TestCommitteeRolesOrder(t *testing.T) {
	r := CommitteeRoles()
	if r[0] != Researcher || r[1] != Challenger || r[2] != Risk {
		t.Fatalf("%v", r)
	}
	_ = os.DevNull
}
