package runtimelib

import (
    "os/exec"
)

// RunCommand runs a command with arguments (space-separated) and returns
// combined stdout/stderr as string and an exit code (0 on success, -1 on error).
func RunCommand(cmd string, args string) (string, int) {
    var argv []string
    if args != "" {
        // naive split by space; keep simple for runtime helper
        // callers should avoid complex quoting for now
        argv = append(argv, args)
    }
    c := exec.Command(cmd, argv...)
    out, err := c.CombinedOutput()
    if err != nil {
        // try to extract exit code if possible
        if ee, ok := err.(*exec.ExitError); ok {
            return string(out), ee.ExitCode()
        }
        return string(out), -1
    }
    return string(out), 0
}
