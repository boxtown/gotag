package gotag

import (
	"fmt"
	"testing"
)

const (
	// Integration is a flag for integration tests
	Integration = "integration"

	// EndToEnd is a flag for end to end tests
	EndToEnd = "endtoend"
)

// TestFlags contains information necessary
// to run or skip tests
type TestFlags struct {
	skip map[string]bool

	// EditDistance is the maximum distance between a test flag
	// and a registered flag that will trigger a skip if Strict
	// is true
	EditDistance int

	// If Strict is true, will skip tests that are within
	// EditDistance of a registered skipped flag and output
	// to stdout why the skip occurred
	Strict bool
}

// New constructs a new instance of TestFlags
func New() *TestFlags {
	return &TestFlags{skip: make(map[string]bool), EditDistance: 2}
}

// Skip marks a test flag to be skipped when testing
// within the context of the TestFlags instance
func (tf *TestFlags) Skip(flag string) {
	tf.skip[flag] = true
}

// Test executes a test under the given flag with the given testing environment
// within the context of the TestFlags instance
func (tf *TestFlags) Test(flag string, t *testing.T, testFn func(t *testing.T)) {
	if _, ok := tf.skip[flag]; ok {
		t.SkipNow()
	} else {
		if tf.Strict {
			for k := range tf.skip {
				if levenshtein(k, flag) < tf.EditDistance {
					fmt.Printf(
						"Found registered skip flag %s within %d edit distance of flag %s, skipping\n",
						k, tf.EditDistance, flag)
					t.SkipNow()
					return
				}
			}
		}
		testFn(t)
	}
}

// Benchmark executes a benchmark under the given flag with the given benchmarking
// environment within the context of the TestFlags instance
func (tf *TestFlags) Benchmark(flag string, b *testing.B, benchmarkFn func(b *testing.B)) {
	if _, ok := tf.skip[flag]; ok {
		b.SkipNow()
	} else {
		if tf.Strict {
			for k := range tf.skip {
				if levenshtein(k, flag) < tf.EditDistance {
					fmt.Printf(
						"Found registered skip flag %s within %d edit distance of flag %s, skipping\n",
						k, tf.EditDistance, flag)
					b.SkipNow()
					return
				}
			}
		}
		benchmarkFn(b)
	}
}

var tf *TestFlags

func init() {
	tf = New()
}

// Skip marks a test flag to be skipped when running tests
// within the default context
func Skip(flag string) {
	tf.Skip(flag)
}

// Test executes a test under the given flag with the given testing
// environment within the default context
func Test(flag string, t *testing.T, testFn func(t *testing.T)) {
	tf.Test(flag, t, testFn)
}

// Benchmark executes a benchmark under the given flag with the
// the given benchmarking environment withint the default context
func Benchmark(flag string, b *testing.B, benchmarkFn func(b *testing.B)) {
	tf.Benchmark(flag, b, benchmarkFn)
}

// iterative implementation of levenshtein distance algorithm
// between 2 strings.
//
// Sourced from https://en.wikipedia.org/wiki/Levenshtein_distance
func levenshtein(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	n1 := len(s1)
	n2 := len(s2)
	if n1 == 0 {
		return n2
	}
	if n2 == 0 {
		return n1
	}

	v0 := make([]int, n2+1)
	v1 := make([]int, n2+1)
	for i := 0; i < n2+1; i++ {
		v0[i] = i
	}
	for i := 0; i < n1; i++ {
		v1[0] = i + 1
		for j := 0; j < n2; j++ {
			if s1[i] == s2[j] {
				v1[j+1] = min(v1[j]+1, v0[j+1]+1, v0[j])
			} else {
				v1[j+1] = min(v1[j]+1, v0[j+1]+1, v0[j]+1)
			}
		}
		copy(v0, v1)
	}
	return v1[n2]
}

// Returns the minimum of all passed in values.
// Returns 0 if no values are passed in.
func min(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}

	min := vals[0]
	for i := 1; i < len(vals); i++ {
		if vals[i] < min {
			min = vals[i]
		}
	}
	return min
}
