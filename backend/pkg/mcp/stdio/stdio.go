package stdio

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"
    "sync"
    "time"
)

// Test sends a lightweight JSON ping to a stdio-backed connector process and
// waits for a single-line JSON response within the provided timeout.
// cmdName is the executable path, args are its args, env is a map of env vars.
func Test(ctx context.Context, cmdName string, args []string, env map[string]string, timeout time.Duration) error {
    if cmdName == "" {
        return fmt.Errorf("command required")
    }

    // create a timeout-aware context
    testCtx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    cmd := exec.CommandContext(testCtx, cmdName, args...)

    // merge environment
    baseEnv := os.Environ()
    for k, v := range env {
        baseEnv = append(baseEnv, fmt.Sprintf("%s=%s", k, v))
    }
    cmd.Env = baseEnv

    stdin, err := cmd.StdinPipe()
    if err != nil {
        return fmt.Errorf("failed to get stdin pipe: %w", err)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("failed to get stdout pipe: %w", err)
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        stderr = io.NopCloser(strings.NewReader(""))
    }

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start command: %w", err)
    }

    // Always reap the process after Start() to avoid zombies on Unix.
    waitCh := make(chan error, 1)
    go func() {
        waitCh <- cmd.Wait()
    }()

    var reapOnce sync.Once
    reap := func(kill bool) {
        reapOnce.Do(func() {
            _ = stdin.Close()
            if kill && cmd.Process != nil {
                _ = cmd.Process.Kill()
            }
            <-waitCh
        })
    }

    // write a single ping JSON line
    ping := map[string]string{"mcp": "ping"}
    pingBytes, _ := json.Marshal(ping)
    if _, err := stdin.Write(append(pingBytes, '\n')); err != nil {
        // attempt to kill process on write error, then reap
        reap(true)
        return fmt.Errorf("failed to write ping: %w", err)
    }
    _ = stdin.Close()

    // read a single line response (newline-terminated) within context
    reader := bufio.NewReader(stdout)
    respCh := make(chan string, 1)
    errCh := make(chan error, 1)

    go func() {
        line, err := reader.ReadString('\n')
        if err != nil {
            // try to capture stderr for better diagnostics
            buf := make([]byte, 4096)
            n, _ := stderr.Read(buf)
            if n > 0 {
                errCh <- fmt.Errorf("read error: %v, stderr: %s", err, strings.TrimSpace(string(buf[:n])))
                return
            }
            errCh <- err
            return
        }
        respCh <- line
    }()

    select {
    case <-testCtx.Done():
        // ensure process is killed and reaped
        reap(true)
        return fmt.Errorf("timeout waiting for stdio response: %w", testCtx.Err())
    case err := <-errCh:
        reap(true)
        return fmt.Errorf("failed to read response: %w", err)
    case line := <-respCh:
        // try to decode JSON (be permissive)
        var obj map[string]interface{}
        if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &obj); err != nil {
            // still consider success if non-empty
            if strings.TrimSpace(line) == "" {
                reap(true)
                return fmt.Errorf("empty response from stdio connector")
            }
            // best-effort success
        }
        // clean up process if it didn't exit, and always reap
        reap(true)
        return nil
    }
}
