package runtimelib

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

// HMACSHA256 returns hex-encoded HMAC-SHA256 of data using key.
func HMACSHA256(key, data string) string {
    mac := hmac.New(sha256.New, []byte(key))
    mac.Write([]byte(data))
    sum := mac.Sum(nil)
    return hex.EncodeToString(sum)
}
