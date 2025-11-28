package runtimelib

import (
    "encoding/hex"
    "hash/crc32"
)

// CRC32Hex returns the hex-encoded IEEE CRC32 checksum of s.
func CRC32Hex(s string) string {
    sum := crc32.ChecksumIEEE([]byte(s))
    buf := make([]byte, 4)
    buf[0] = byte(sum >> 24)
    buf[1] = byte(sum >> 16)
    buf[2] = byte(sum >> 8)
    buf[3] = byte(sum)
    return hex.EncodeToString(buf)
}
