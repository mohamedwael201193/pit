package compute

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestClassifyCatalogNeverPrivate(t *testing.T) {
	e := classifyCatalog("claude-opus-5", "", "glm-5.2")
	if e.PrivateBook || e.Path != "catalog-listing" {
		t.Fatalf("%+v", e)
	}
	if !strings.Contains(strings.ToLower(e.Note), "router") && !strings.Contains(strings.ToLower(e.Note), "not") {
		t.Fatal(e.Note)
	}
	tls := classifyCatalog("glm-5.3", "TeeTLS", "glm-5.2")
	if tls.PrivateBook || tls.UsableFor == "" {
		t.Fatalf("%+v", tls)
	}
}

func TestCatalogUsableForChat(t *testing.T) {
	ok, why := CatalogUsableForChat("host-parsed", config.Mainnet)
	if !ok || why != "host-parsed" {
		t.Fatalf("%v %s", ok, why)
	}
	ok, _ = CatalogUsableForChat("glm-5.2", config.Mainnet)
	if !ok {
		t.Fatal("direct sku")
	}
	ok, _ = CatalogUsableForChat("glm-5.3", config.Mainnet)
	if !ok {
		t.Fatal("direct sku alias")
	}
	ok, why = CatalogUsableForChat("claude-opus-5", config.Mainnet)
	if ok || why != "not_direct_on_this_workspace" {
		t.Fatalf("%v %s", ok, why)
	}
}
