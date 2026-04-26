package stdio

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"
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
    stderr, _ := cmd.StderrPipe()

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start command: %w", err)
    }

    // write a single ping JSON line
    ping := map[string]string{"mcp": "ping"}
    pingBytes, _ := json.Marshal(ping)
    if _, err := stdin.Write(append(pingBytes, '\n')); err != nil {
        // attempt to kill process on write error
        _ = cmd.Process.Kill()
        return fmt.Errorf("failed to write ping: %w", err)
    }
    stdin.Close()

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
        // ensure process is killed
        _ = cmd.Process.Kill()
        return fmt.Errorf("timeout waiting for stdio response: %w", testCtx.Err())
    case err := <-errCh:
        _ = cmd.Process.Kill()
        return fmt.Errorf("failed to read response: %w", err)
    case line := <-respCh:
        // try to decode JSON (be permissive)
        var obj map[string]interface{}
        if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &obj); err != nil {
            // still consider success if non-empty
            if strings.TrimSpace(line) == "" {
                _ = cmd.Process.Kill()
                return fmt.Errorf("empty response from stdio connector")
            }
            // best-effort success
        }
        // clean up process if it didn't exit
        _ = cmd.Process.Kill()
        return nil
    }
}
