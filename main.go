package main

import (
	"math/rand"
	"time"
)

const inputFile = "roadmap.txt"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	serve(0, "", "")
}
