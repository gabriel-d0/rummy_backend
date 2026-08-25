// Package setup — Deterministic random source (Day 17).
// Injectable/seedable randomness for reproducible shuffles and tests per
// docs/rules-decisions.md:2 MVP shuffle (Fisher–Yates with seed).
package setup

import "math/rand"

// Rand is an injectable random source. It wraps math/rand.Rand and is
// used by Shuffle (Day 18) so tests can pass a seeded instance.
type Rand struct {
	r *rand.Rand
}

// NewSeededRand returns a deterministic Rand for the given seed.
// Two Rands with the same seed produce identical Intn/Shuffle sequences.
func NewSeededRand(seed int64) *Rand {
	return &Rand{r: rand.New(rand.NewSource(seed))}
}

// Intn returns a uniform int in [0,n). Panics if n <=0, matching math/rand.
func (r *Rand) Intn(n int) int {
	return r.r.Intn(n)
}

// Int63 returns a non-negative random 63-bit integer.
func (r *Rand) Int63() int64 {
	return r.r.Int63()
}

// Shuffle randomizes the order of n elements via swap, using Fisher–Yates
// via the underlying Rand. It is a thin wrapper around rand.Shuffle.
func (r *Rand) Shuffle(n int, swap func(i, j int)) {
	r.r.Shuffle(n, swap)
}
