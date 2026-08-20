package storage

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/google/uuid"
)

func TestObjectKeyCrossWorkspace(t *testing.T) {
	a := uuid.NewString()
	b := uuid.NewString()
	ka, err := ObjectKey(config.Mainnet, a, "observation", "book")
	if err != nil {
		t.Fatal(err)
	}
	if err := AssertWorkspace(ka, b); err == nil {
		t.Fatal("A key accepted as B")
	}
	if err := AssertWorkspace(ka, a); err != nil {
		t.Fatal(err)
	}
}
