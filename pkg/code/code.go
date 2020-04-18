package code

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/peteraba/roadmapper/pkg/helper"
)

const maxCode64 = 0xfffffffffff

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

// NewCode64 returns a random Code64
func NewCode64() Code64 {
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

	for idx, runeValue := range helper.Reverse(s) {
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

const allowed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_~"

var allowedLock = &sync.Mutex{}
var allowedMap map[rune]int

// getAllowedMap returns a map between runes and their numeric values
func getAllowedMap() map[rune]int {
	allowedLock.Lock()
	defer allowedLock.Unlock()

	if len(allowedMap) > 0 {
		return allowedMap
	}

	allowedMap = map[rune]int{}
	for idx, runeValue := range allowed {
		allowedMap[runeValue] = idx
	}

	return allowedMap
}
