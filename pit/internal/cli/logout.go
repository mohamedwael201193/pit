package cli

import (
	"os"

	"github.com/mohamedwael201193/pit/internal/keyring"
)

func Logout(dir string, forget bool) error {
	st, err := Load(dir)
	if err != nil && !forget {
		return err
	}
	sf, _ := LoadSession(dir)
	if sf.ID != "" && st.WorkspaceID != "" {
		if ring, err := keyring.OpenProduct(KeyringDir(dir)); err == nil {
			_ = ring.Delete(st.WorkspaceID+"/session", sf.ID)
		}
	}
	_ = os.Remove(SessionPath(dir))
	if forget {
		_ = os.Remove(Path(dir))
	}
	return nil
}
