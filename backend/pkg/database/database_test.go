package database

import (
    "testing"
    "unicode/utf8"
)

func TestSanitizeUTF8_Empty(t *testing.T) {
    got := SanitizeUTF8("")
    if got != "" {
        t.Fatalf("expected empty string, got %q", got)
    }
}

func TestSanitizeUTF8_RemovesNullsAndInvalid(t *testing.T) {
    // build string with null byte and invalid utf8 sequence
    raw := string([]byte{'a', 0x00, 'b', 0xff})
    got := SanitizeUTF8(raw)

    // null byte should be removed
    if containsNull := (len(got) != 0 && got[1] == 0); containsNull {
        t.Fatalf("result contains null byte: %q", got)
    }

    // last rune should be the unicode replacement rune for invalid utf8
    r, _ := utf8.DecodeLastRuneInString(got)
    if r != utf8.RuneError {
        t.Fatalf("expected last rune to be RuneError, got %U (string %q)", r, got)
    }
}
