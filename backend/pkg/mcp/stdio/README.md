STDIO connector
================

This package provides a small helper to run a lightweight connectivity test against a process that
implements the MCP stdio adapter protocol (newline-delimited JSON ping/pong). It is intentionally
minimal and used by the backend `testMcpServer` resolver.

Usage
-----

Call `Test(ctx, cmdName, args, envMap, timeout)` to spawn the process, send a single JSON ping
and wait for one JSON response line. The process is killed after the response.

Note: The helper executes the provided executable directly; do not pass shell-interpolated
strings as `cmdName` (use `bash -c` only for testing).
