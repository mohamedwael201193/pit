package memory

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/storage"
)

func ListDenied(viewer, owner, kind, name string) error {
	key, err := storage.ObjectKey(config.Mainnet, owner, kind, name)
	if err != nil {
		return err
	}
	if err := storage.AssertWorkspace(key, viewer); err == nil {
		return fmt.Errorf("memory_list_cross_access")
	}
	return nil
}
