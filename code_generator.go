package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"unicode/utf8"
)

const maxCode64 = 0xfffffffffff

// CodeBuilder is a factory for codes
type CodeBuilder struct {
}

// NewFromString creates a new code from a string representation (typically a URL path)
func (cb CodeBuilder) NewFromString(s string) (Code, error) {
	return newCode64FromString(s)
}

// NewFromID creates a new code from a number
func (cb CodeBuilder) NewFromID(id uint64) (Code, error) {
	return Code64(id), nil
}

// New creates a new random code
func (cb CodeBuilder) New() Code {
	return newCode64()
}

// NewCodeBuilder creates a new CodeBuilder instance
func NewCodeBuilder() CodeBuilder {
	return CodeBuilder{}
}

// Code represents a positive integer that also has a human-readable, url-friendly, rather minimalistic representation
type Code interface {
	String() string
	ID() uint64
}

// Code64 is a 64-bit Code
type Code64 uint64

// String returns a human-readable, url-friendly representation of a 64-bit Code
func (c Code64) String() string {
	if c > maxCode64 {
		panic("code out of bound")
	}

	return toCode64(uint64(c))
}

// ID returns a numeric representation of a 64-bit Code
func (c Code64) ID() uint64 {
	return uint64(c)
}

// newCode64 returns a random Code64
func newCode64() Code64 {
	return Code64(rand.Int63n(maxCode64))
}

// newCode64FromString creates a new Code64 from a string representation
func newCode64FromString(s string) (Code64, error) {
	var (
		n   uint64
		m   = getAllowedMap()
		num int
		ok  bool
	)

	for idx, runeValue := range reverse(s) {
		num, ok = m[runeValue]
		if !ok {
			return Code64(0), fmt.Errorf("invalid character '%c' in code: %s", runeValue, s)
		}

		n += uint64(num) << (idx * 6)
	}

	if n > maxCode64 {
		return Code64(0), fmt.Errorf("code is out of bound: %d", n)
	}

	return Code64(n), nil
}

// toCode64 converts a number into a string representation
func toCode64(n uint64) string {
	// mask the first two bits, we'll only use 30
	n &= maxCode64

	var runes []rune
	for i := 0; i < 10; i++ {
		shift := 54 - i*6
		n := (n & (0x3f << shift)) >> shift
		runes = append(runes, rune(allowed[n]))
	}

	return strings.TrimLeft(string(runes), "0")
}

// reverse reverses a string
func reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

const allowed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_~"

var allowedLock = &sync.Mutex{}
var allowedMap map[rune]int

// getAllowedMap returns a map between runes and their numeric values
func getAllowedMap() map[rune]int {
	if len(allowedMap) > 0 {
		return allowedMap
	}

	allowedLock.Lock()
	defer allowedLock.Unlock()

	// re-check to see if concurrent run has already pre-filled the map
	if len(allowedMap) > 0 {
		return allowedMap
	}

	allowedMap = map[rune]int{}
	for idx, runeValue := range allowed {
		allowedMap[runeValue] = idx
	}

	return allowedMap
}
