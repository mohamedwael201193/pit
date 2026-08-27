package cli

import (
	"path/filepath"

	"github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/ledger"
)

func LedgerDir(dir string) string {
	return filepath.Join(dir, "ledger")
}

func RememberAuthorized(dir, network, workspace, cloid, previewHash string) error {
	st, err := ledger.Open(LedgerDir(dir), network, workspace)
	if err != nil {
		return err
	}
	defer st.Close()
	ok, err := st.Apply(ledger.Record{
		Workspace: workspace,
		Cloid:     cloid,
		Preview:   previewHash,
		Status:    ledger.StatusAuthorized,
	})
	return exec.DuplicateClick(ok, err)
}

func RememberPosted(dir, network, workspace, cloid, oid string) error {
	st, err := ledger.Open(LedgerDir(dir), network, workspace)
	if err != nil {
		return err
	}
	defer st.Close()
	return st.Mark(workspace, cloid, ledger.StatusReceipt, oid)
}

func LookupAction(dir, network, workspace, cloid string) (ledger.Record, error) {
	st, err := ledger.Open(LedgerDir(dir), network, workspace)
	if err != nil {
		return ledger.Record{}, err
	}
	defer st.Close()
	return st.Get(workspace, cloid)
}
