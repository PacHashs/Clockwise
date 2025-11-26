package checker

import "fmt"

// FormatDiagnostics returns a human-readable string for diagnostics.
func FormatDiagnostics(diags Diagnostics) string {
    if len(diags) == 0 {
        return "no diagnostics"
    }
    s := ""
    for i, d := range diags {
        s += fmt.Sprintf("%d: %s\n", i+1, d.Msg)
    }
    return s
}
