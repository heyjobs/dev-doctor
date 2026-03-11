package diagnostics

import (
	"math/rand"
	"time"
)

// Mock diagnostic implementations
// These return randomized results to simulate a realistic diagnostic run

func init() {
	rand.Seed(time.Now().UnixNano())
}
