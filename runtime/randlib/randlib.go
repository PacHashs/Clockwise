package runtimelib

import (
    "math/rand"
    "time"
)

func RandInt(n int) int {
    rand.Seed(time.Now().UnixNano())
    if n <= 0 {
        return 0
    }
    return rand.Intn(n)
}
