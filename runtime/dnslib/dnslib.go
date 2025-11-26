package runtimelib

import (
    "net"
    "strings"
)

// LookupHost returns a comma-separated list of IP addresses for the given host
// or an empty string on error.
func LookupHost(host string) string {
    ips, err := net.LookupHost(host)
    if err != nil {
        return ""
    }
    return strings.Join(ips, ",")
}
