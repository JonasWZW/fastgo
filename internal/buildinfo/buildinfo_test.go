package buildinfo

import (
	"strings"
	"testing"
)

func TestStringContainsEveryField(t *testing.T) {
	got := String()

	for _, want := range []string{"version=", "commit=", "built="} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %q, want it to contain %q", got, want)
		}
	}
}
