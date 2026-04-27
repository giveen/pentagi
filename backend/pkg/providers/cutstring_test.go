package providers

import (
    "strings"
    "testing"
)

func TestCutString_NoTruncate(t *testing.T) {
    s := "short"
    got := cutString(s, 10)
    if got != s {
        t.Fatalf("expected %q, got %q", s, got)
    }
}

func TestCutString_Truncate(t *testing.T) {
    s := strings.Repeat("a", 50)
    got := cutString(s, 10)
    if len(got) <= 10 {
        t.Fatalf("expected truncated suffix, got %q", got)
    }
    if !strings.Contains(got, "[truncated full size is") {
        t.Fatalf("expected truncation note, got %q", got)
    }
}
