package script

import (
	"testing"
)

func TestValidateWhitelist(t *testing.T) {
	tests := []string{
		"Window.GetWidth()",
		"Sleep(1)",
		"new Image('x.png')",
		"SetRefreshFunction(function(){})",
	}
	for _, s := range tests {
		if err := Validate(s); err != nil {
			t.Errorf("Validate(%q) error: %v", s, err)
		}
	}
}

func TestValidateUndefined(t *testing.T) {
	if err := Validate("fooBar()"); err == nil {
		t.Fatal("expected error for undefined function")
	}
}

func TestValidateForbidden(t *testing.T) {
	if err := Validate("eval('1+1')"); err == nil {
		t.Fatal("expected error for forbidden token")
	}
}
