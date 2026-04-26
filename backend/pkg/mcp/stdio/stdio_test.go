package stdio

import (
    "context"
    "testing"
    "time"
)

func TestTest_SucceedsWithBashEcho(t *testing.T) {
    ctx := context.Background()

    // Use bash -c to run a small inline script that reads a line and prints a JSON response
    cmdName := "bash"
    args := []string{"-c", `read line; echo '{"pong":"ok"}'`}

    // should complete quickly
    if err := Test(ctx, cmdName, args, map[string]string{}, 2*time.Second); err != nil {
        t.Fatalf("expected no error, got: %v", err)
    }
}
