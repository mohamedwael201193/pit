package policy

import "testing"

func TestCardsCoverLaw(t *testing.T) {
	cs := Cards(Default())
	if len(cs) != 14 {
		t.Fatalf("%d", len(cs))
	}
	seen := map[string]bool{}
	for _, c := range cs {
		if c.Title == "" || c.Value == "" || c.Law == "" {
			t.Fatalf("%+v", c)
		}
		if seen[c.Title] {
			t.Fatal(c.Title)
		}
		seen[c.Title] = true
	}
}
