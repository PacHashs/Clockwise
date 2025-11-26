package runtimelib

import (
    "time"
)

func NowISO() string {
    return time.Now().Format(time.RFC3339)
}

func SleepMs(ms int) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
}
