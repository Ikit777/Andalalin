package utils

import (
    "crypto/rand"
    "io"
)

var number = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func Encode(s int) string {
	b := make([]byte, s)
    n, err := io.ReadAtLeast(rand.Reader, b, s)
    if n != s {
        panic(err)
    }
    for i := 0; i < len(b); i++ {
        b[i] = number[int(b[i])%len(number)]
    }
    return string(b)
}