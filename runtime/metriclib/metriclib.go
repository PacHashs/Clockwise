package runtimelib

import "time"

// TimestampMs returns current unix time in milliseconds.
func TimestampMs() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
