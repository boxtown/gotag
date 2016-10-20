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

// TestContext contains information necessary
// to run or skip tests
type TestContext struct {
	skip map[string]bool

	// Verbose will print information messages
	// if set to true
	Verbose bool

	// EditDistance is the maximum distance between a test flag
	// and a registered flag that will trigger a skip if Fuzzy
	// is true
	EditDistance int

	// If Fuzzy is true, will skip tests that are within
	// EditDistance of a registered skipped flag and output
	// to stdout why the skip occurred
	Fuzzy bool
}

// New constructs a new instance of TestContext
func New() *TestContext {
	return &TestContext{skip: make(map[string]bool), EditDistance: 2}
}

// Skip marks a test tag to be skipped when testing
// within the context of the TestContext instance
func (tc *TestContext) Skip(tag string) {
	tc.skip[tag] = true
}

// Test executes a test under the given tag with the given testing environment
// within the context of the TestContext instance
func (tc *TestContext) Test(tag string, t *testing.T, testFn func(t *testing.T)) {
	if _, ok := tc.skip[tag]; ok {
		t.SkipNow()
	} else {
		if tc.Fuzzy && tc.checkFuzzy(t, tag) {
			return
		}
		testFn(t)
	}
}

// Benchmark executes a benchmark under the given tag with the given benchmarking
// environment within the context of the TestFlags instance
func (tc *TestContext) Benchmark(tag string, b *testing.B, benchmarkFn func(b *testing.B)) {
	if _, ok := tc.skip[tag]; ok {
		b.SkipNow()
	} else {
		if tc.Fuzzy && tc.checkFuzzy(b, tag) {
			return
		}
		benchmarkFn(b)
	}
}

func (tc *TestContext) checkFuzzy(s skippable, tag string) bool {
	for k := range tc.skip {
		if levenshtein(k, tag) > tc.EditDistance {
			continue
		}

		if tc.Verbose {
			fmt.Printf(
				"Found registered skip tag '%s' within an edit distance of %d of tag '%s', skipping...\n",
				k, tc.EditDistance, tag)
		}
		s.SkipNow()
		return true
	}
	return false
}

type skippable interface {
	Skip(...interface{})
	SkipNow()
}

var tc *TestContext

func init() {
	tc = New()
}

// Skip marks a test tag to be skipped when running tests
// within the default context
func Skip(tag string) {
	tc.Skip(tag)
}

// Fuzzy sets fuzzy matching for the default context
func Fuzzy(fuzzy bool) {
	tc.Fuzzy = fuzzy
}

// Distance sets the fuzzy matching distance for the default context
func Distance(distance int) {
	tc.EditDistance = distance
}

// Test executes a test under the given tag with the given testing
// environment within the default context
func Test(tag string, t *testing.T, testFn func(t *testing.T)) {
	tc.Test(tag, t, testFn)
}

// Benchmark executes a benchmark under the given tag with the
// the given benchmarking environment withint the default context
func Benchmark(tag string, b *testing.B, benchmarkFn func(b *testing.B)) {
	tc.Benchmark(tag, b, benchmarkFn)
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
