package policy

import "testing"

func TestDefaultHashStable(t *testing.T) {
	h, err := Default().Hash()
	if err != nil {
		t.Fatal(err)
	}
	const want = "384bfd62f42d9a3b167ab043aae468ece453cb328d27d825d7747ecb007ccab2"
	if h != want {
		t.Fatalf("default pin identity changed: %s", h)
	}
}
