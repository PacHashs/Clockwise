package runtimelib

import (
    "encoding/csv"
    "strings"
)

// CSVToLines parses CSV content and returns each record joined by comma.
// Returns an empty slice on parse error.
func CSVToLines(content string) []string {
    r := csv.NewReader(strings.NewReader(content))
    recs, err := r.ReadAll()
    if err != nil {
        return []string{}
    }
    out := make([]string, 0, len(recs))
    for _, rec := range recs {
        out = append(out, strings.Join(rec, ","))
    }
    return out
}
