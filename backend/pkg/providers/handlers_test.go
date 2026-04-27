package providers

import (
    "context"
    "strings"
    "testing"
)

// Test that when the summarizer always returns a fixed short output, the
// chunked summarizer reduces multiple chunks down to that short output via
// iterative summarization.
func TestChunkedSummarizeIterativeReduction(t *testing.T) {
    var sb strings.Builder
    // build input of moderate size and use a small chunk size for the test
    for sb.Len() < 400 {
        sb.WriteString("a")
    }
    input := sb.String()

    calls := 0
    // summarizer: for chunk-sized inputs return a moderately large placeholder
    // so that combined partials exceed chunkSize and trigger iterative reduction.
    summarizer := func(ctx context.Context, txt string) (string, error) {
        calls++
        if len(txt) > 200 {
            return "REDUCED", nil
        }
        return strings.Repeat("A", 90), nil
    }

    // Use a small chunk size to ensure multiple chunks are created
    out, err := chunkedSummarize(context.Background(), input, 100, 3, summarizer)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out != "REDUCED" {
        t.Fatalf("expected REDUCED, got %q (calls=%d)", out, calls)
    }
}

// Test that a small input results in a single summarizer call and the same
// summarizer output is returned.
func TestChunkedSummarizeSingleCall(t *testing.T) {
    input := "short input"

    calls := 0
    summarizer := func(ctx context.Context, txt string) (string, error) {
        calls++
        return "SINGLE", nil
    }

    out, err := chunkedSummarize(context.Background(), input, msgSummarizerLimit, 3, summarizer)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out != "SINGLE" {
        t.Fatalf("expected SINGLE, got %q", out)
    }
    if calls != 1 {
        t.Fatalf("expected 1 call, got %d", calls)
    }
}
