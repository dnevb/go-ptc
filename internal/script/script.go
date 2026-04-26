package script

import (
	"fmt"
	"regexp"
	"strings"
)

var callRe = regexp.MustCompile(`[\w\.]+\s*\(`)

var whitelist = map[string]bool{
	"SetImage":              true,
	"SetImageScale":         true,
	"SetImageOpacity":       true,
	"SetImagePosition":      true,
	"SetText":               true,
	"SetTextOpacity":        true,
	"SetTextPosition":       true,
	"GetImageWidth":         true,
	"GetImageHeight":        true,
	"AddProgressAnimation":  true,
	"Refresh":               true,
	"Sleep":                 true,
	"FadeIn":                true,
	"FadeOut":               true,
	"new":                   true,
	"Image":                 true,
	"Window":                true,
	"SetRefreshFunction":    true,
	"GetWidth":              true,
	"GetHeight":             true,
	"SetPosition":           true,
	"Draw":                  true,
	"DrawBox":               true,
	"SetBackgroundColor":    true,
	"SetMessage":            true,
	"SetProgress":           true,
	"function":              true,
}

var forbidden = []string{
	"eval",
	"XMLHttpRequest",
	"fetch",
	"require",
	"import",
}

func Validate(content string) error {
	for _, tok := range forbidden {
		if strings.Contains(content, tok) {
			return fmt.Errorf("forbidden token: %s", tok)
		}
	}

	calls := callRe.FindAllString(content, -1)
	for _, c := range calls {
		c = strings.TrimSpace(strings.TrimSuffix(c, "("))
		base := c
		if idx := strings.LastIndex(c, "."); idx != -1 {
			base = c[idx+1:]
		}
		if !whitelist[base] {
			return fmt.Errorf("undefined function call: %s", c)
		}
	}
	return nil
}
