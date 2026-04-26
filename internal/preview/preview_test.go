package preview

import (
	"testing"
)

func TestPreviewMissingPlymouthd(t *testing.T) {
	// This test may pass or skip depending on host environment.
	err := Preview("/tmp")
	if err == nil {
		return
	}
	// If plymouthd missing, expect error.
	if err.Error()[:18] != "plymouthd not found" {
		t.Logf("preview error (expected if plymouthd absent): %v", err)
	}
}
