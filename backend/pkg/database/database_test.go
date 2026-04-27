package database

import (
	"strings"
	"testing"
)

func TestSanitizeUTF8_Empty(t *testing.T) {
	if got := SanitizeUTF8(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestSanitizeUTF8_RemovesNullsAndInvalid(t *testing.T) {
	input := "A\x00\xffB"

	got := SanitizeUTF8(input)

	if strings.ContainsRune(got, '\x00') {
		t.Fatalf("result contains null byte: %q", got)
	}
	if !strings.ContainsRune(got, '\ufffd') {
		t.Fatalf("expected replacement rune for invalid UTF-8, got %q", got)
	}
	if !strings.HasPrefix(got, "A") || !strings.HasSuffix(got, "B") {
		t.Fatalf("expected content to preserve valid surrounding bytes, got %q", got)
	}
}