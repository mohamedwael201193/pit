package cli

import "testing"

func TestPreviewCopy(t *testing.T) {
	if PreviewCopy == "" || MutationInvalidates() == "" {
		t.Fatal("copy")
	}
}
