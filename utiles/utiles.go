package utiles

import (
	"crypto/rand"
	"fmt"
)

func Max(x, y int) int {
	if x >= y {
		return x
	}

	return y
}

func Min(x, y int) int {
	if x >= y {
		return y
	}

	return x
}

func UniqueIdGeneretor() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
