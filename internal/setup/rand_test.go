package setup

import "testing"

func TestSeededRandDeterministicIntn(t *testing.T) {
	seed := int64(42)
	a := NewSeededRand(seed)
	b := NewSeededRand(seed)
	for i := 0; i < 100; i++ {
		ai := a.Intn(1000)
		bi := b.Intn(1000)
		if ai != bi {
			t.Fatalf("Intn mismatch at %d: %d vs %d", i, ai, bi)
		}
	}
}

func TestSeededRandDifferentSeedsDiverge(t *testing.T) {
	a := NewSeededRand(1)
	b := NewSeededRand(2)
	// Very unlikely to be identical for 10 draws
	same := true
	for i := 0; i < 10; i++ {
		if a.Intn(1000) != b.Intn(1000) {
			same = false
			break
		}
	}
	if same {
		t.Fatalf("different seeds unexpectedly gave same 10 Intn sequence")
	}
}

func TestSeededRandShuffleDeterministic(t *testing.T) {
	seed := int64(12345)
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	b := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r1 := NewSeededRand(seed)
	r2 := NewSeededRand(seed)
	r1.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	r2.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("Shuffle mismatch at %d: %v vs %v", i, a, b)
		}
	}
	// Ensure shuffle actually changed order (not identity)
	unchanged := true
	for i, v := range a {
		if v != i {
			unchanged = false
			break
		}
	}
	if unchanged {
		t.Fatalf("Shuffle with seed %d did not change order: %v", seed, a)
	}
}

func TestSeededRandInt63Deterministic(t *testing.T) {
	a := NewSeededRand(99)
	b := NewSeededRand(99)
	for i := 0; i < 50; i++ {
		if a.Int63() != b.Int63() {
			t.Fatalf("Int63 mismatch at %d", i)
		}
	}
}
