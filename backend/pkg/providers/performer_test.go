package providers

import (
	"errors"
	"testing"
)

// TestShouldRepairToolCallArgs verifies the non-repairable error classification.
// Previously "failed to store tool result in long-term memory" was non-repairable,
// which caused the entire agent chain to abort when vector store writes failed
// (e.g. embedding batch size exceeded). It must now be repairable so the chain
// continues and only a warning is logged.
func TestShouldRepairToolCallArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		wantRepair  bool
	}{
		{
			name:       "nil error returns false (nothing to repair)",
			err:        nil,
			wantRepair: false,
		},
		{
			name:       "generic error is repairable",
			err:        errors.New("unexpected token in args"),
			wantRepair: true,
		},
		{
			name:       "update toolcall result error is non-repairable",
			err:        errors.New("failed to update toolcall result: connection refused"),
			wantRepair: false,
		},
		// The following two were previously non-repairable and caused task aborts
		// when the embedding model rejected oversized inputs (>512 tokens).
		// They must now be treated as repairable (warn + continue).
		{
			name:       "long-term memory store error is now repairable",
			err:        errors.New("failed to store tool result in long-term memory: failed to create openai embeddings: API returned unexpected status code: 500: input (720 tokens) is too large to process"),
			wantRepair: true,
		},
		{
			name:       "store tool result error is now repairable",
			err:        errors.New("failed to store tool result: error embedding batch"),
			wantRepair: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := shouldRepairToolCallArgs(tt.err)
			if got != tt.wantRepair {
				t.Errorf("shouldRepairToolCallArgs(%v) = %v, want %v", tt.err, got, tt.wantRepair)
			}
		})
	}
}
