package theme

import (
	"strings"
	"testing"
)

func TestGeneratePlymouth(t *testing.T) {
	got, err := GeneratePlymouth("mytheme")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"[Plymouth Theme]",
		"Name=mytheme",
		"ModuleName=script",
		"ImageDir=assets",
		"ScriptFile=mytheme.script",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in generated .plymouth", want)
		}
	}
}

func TestValidatePlymouthMissingKey(t *testing.T) {
	bad := "[Plymouth Theme]\nName=x\n"
	if err := ValidatePlymouth(bad); err == nil {
		t.Fatal("expected error for missing key")
	}
}
