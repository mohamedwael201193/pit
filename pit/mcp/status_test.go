package mcp

import "testing"

func TestStatusNeverSigns(t *testing.T) {
	r := StatusNeverSigns()
	body, _ := r.Body.(map[string]any)
	if body["sign"] != false || body["authorize"] != false {
		t.Fatal(body)
	}
}
