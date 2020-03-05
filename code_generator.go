package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"unicode/utf8"
)

const maxCode64 = 0xfffffffffff

type CodeBuilder struct {
}

func (cb CodeBuilder) NewFromString(s string) (Code, error) {
	return newCode64FromString(s)
}

func (cb CodeBuilder) NewFromID(id int64) (Code, error) {
	return Code64(id), nil
}

func (cb CodeBuilder) New() Code {
	return newCode64()
}

func NewCodeBuilder() CodeBuilder {
	return CodeBuilder{}
}

type Code interface {
	String() string
	ID() int64
}

type Code64 int64

func (c Code64) String() string {
	if c > maxCode64 || c < 0 {
		panic("code out of bound")
	}

	return toCode64(int64(c))
}
func (c Code64) ID() int64 {
	return int64(c)
}

func newCode64() Code64 {
	return Code64(rand.Int63n(maxCode64))
}

func newCode64FromString(s string) (Code64, error) {
	var (
		n   int64
		m   = getAllowedMap()
		num int
		ok  bool
	)

	for idx, runeValue := range reverse(s) {
		num, ok = m[runeValue]
		if !ok {
			return Code64(0), fmt.Errorf("invalid character '%c' in code: %s", runeValue, s)
		}

		n += int64(num) << (idx * 6)
	}

	if n > maxCode64 {
		return Code64(0), fmt.Errorf("code is out of bound: %d", n)
	}

	return Code64(n), nil
}

func toCode64(n int64) string {
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
