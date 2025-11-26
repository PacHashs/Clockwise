package runtimelib

import (
    "crypto/rand"
    "fmt"
)

// UUIDv4 returns a random UUIDv4 string (simple implementation)
func UUIDv4() string {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        return ""
    }
    // Set version (4) and variant bits per RFC 4122
    b[6] = (b[6] & 0x0f) | 0x40
    b[8] = (b[8] & 0x3f) | 0x80
    return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
        uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]),
        uint16(b[4])<<8|uint16(b[5]),
        uint16(b[6])<<8|uint16(b[7]),
        uint16(b[8])<<8|uint16(b[9]),
        b[10:16],
    )
}
