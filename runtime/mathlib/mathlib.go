package runtimelib

import "math"

// Math helpers (int)
func Abs(i int) int {
    if i < 0 {
        return -i
    }
    return i
}

func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Floating point helpers
func FloatPow(x, y float64) float64 {
    return math.Pow(x, y)
}

func FloatSqrt(x float64) float64 {
    return math.Sqrt(x)
}

func Sin(x float64) float64 {
    return math.Sin(x)
}

func Cos(x float64) float64 {
    return math.Cos(x)
}

func Floor(x float64) float64 {
    return math.Floor(x)
}

func Ceil(x float64) float64 {
    return math.Ceil(x)
}
