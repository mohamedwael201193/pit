package memory

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/storage"
)

func Isolate(net config.Network, workspaceA, workspaceB, kind, name string) error {
	if workspaceA == workspaceB {
		return fmt.Errorf("same_workspace")
	}
	a, err := storage.ObjectKey(net, workspaceA, kind, name)
	if err != nil {
		return err
	}
	b, err := storage.ObjectKey(net, workspaceB, kind, name)
	if err != nil {
		return err
	}
	if a == b {
		return fmt.Errorf("memory_cross_access")
	}
	if err := storage.AssertWorkspace(a, workspaceB); err == nil {
		return fmt.Errorf("memory_cross_access")
	}
	return storage.AssertWorkspace(a, workspaceA)
}
