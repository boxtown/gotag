package gotag

import (
	"fmt"
	"testing"
)

const (
	// Integration is a flag for integration tests
	Integration = "integration"

	// EndToEnd is a flag for end to end tests
	EndToEnd = "end-to-end"
)

// T is an interface that matches testing.T. This allows
// gotag to actually be testable
type T interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Log(...interface{})
	Logf(string, ...interface{})
	Parallel()
	Run(string, func(*testing.T)) bool
	Skip(...interface{})
	SkipNow()
	Skipf(string, ...interface{})
	Skipped() bool
}

// B is an interface that matches testing.B. This allows
// gotag to actually be testable
type B interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Log(...interface{})
	Logf(string, ...interface{})
	ReportAllocs()
	ResetTimer()
	Run(string, func(*testing.T)) bool
	RunParallel(func(*testing.PB))
	SetBytes(int64)
	SetParallelism(int)
	Skip(...interface{})
	SkipNow()
	Skipf(string, ...interface{})
	Skipped() bool
	StartTimer()
	StopTimer()
}

// TestContext contains information necessary
// to run or skip tests
type TestContext struct {
	skip    map[string]bool
	runOnly map[string]bool

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
	return &TestContext{
		skip:         make(map[string]bool),
		runOnly:      make(map[string]bool),
		EditDistance: 2,
	}
}

// Skip marks test tags to be skipped when testing
// within the context of the TestContext instance
func (tc *TestContext) Skip(tags ...string) {
	for _, tag := range tags {
		tc.skip[tag] = true
	}
}

// RunOnly marks specific tests to be run. If this method is called
// with a non-empty argument, then only the given tests will run.
// Marking tags as run only will by default make the context ignore skipped tags.
func (tc *TestContext) RunOnly(tags ...string) {
	for _, tag := range tags {
		tc.runOnly[tag] = true
	}
}

// Test executes a test under the given tag with the given testing environment
// within the context of the TestContext instance
func (tc *TestContext) Test(tag string, t T, testFn func(t T)) {
	tc.run(tag, t, func(s skippable) {
		testFn(s.(T))
	})
}

// Benchmark executes a benchmark under the given tag with the given benchmarking
// environment within the context of the TestFlags instance
func (tc *TestContext) Benchmark(tag string, b B, benchmarkFn func(b B)) {
	tc.run(tag, b, func(s skippable) {
		benchmarkFn(s.(B))
	})
}

func (tc *TestContext) run(tag string, s skippable, fn func(s skippable)) {
	match, reason := tc.shouldSkip(tag)
	switch reason {
	case foundInSkip, notInRunOnly:
		s.SkipNow()
	case fuzzyMatchSkip:
		if tc.Verbose {
			fmt.Printf(
				"Found registered skip tag '%s' within an edit distance of %d of tag '%s', skipping...\n",
				match, tc.EditDistance, tag)
		}
		s.SkipNow()
	case doNotSkipFuzzy:
		if tc.Verbose {
			fmt.Printf(
				"Found registered run tag '%s' within an edit distance of %d of tag '%s', running...\n",
				match, tc.EditDistance, tag)
		}
		fn(s)
	default:
		fn(s)
	}
}

func (tc *TestContext) shouldSkip(tag string) (string, skipReason) {
	if len(tc.runOnly) > 0 {
		run := tc.runOnly[tag]
		if run {
			return "", doNotSkip
		}
		if !tc.Fuzzy {
			return "", notInRunOnly
		}

		match, runFuzzy := tc.checkFuzzy(tag, tc.runOnly)
		if !runFuzzy {
			return "", notInRunOnly
		}
		return match, doNotSkipFuzzy
	}

	skip := tc.skip[tag]
	if skip {
		return "", foundInSkip
	}
	if !tc.Fuzzy {
		return "", doNotSkip
	}

	match, skipFuzzy := tc.checkFuzzy(tag, tc.skip)
	if !skipFuzzy {
		return "", doNotSkip
	}
	return match, fuzzyMatchSkip
}

func (tc *TestContext) checkFuzzy(tag string, collection map[string]bool) (string, bool) {
	for k := range collection {
		if levenshtein(k, tag) > tc.EditDistance {
			continue
		}
		return k, true
	}
	return "", false
}

// Skip marks test tags to be skipped when running tests
// within the default context
func Skip(tags ...string) {
	tc.Skip(tags...)
}

// RunOnly marks specific tests to be run within the default context.
// If this method is called with a non-empty argument, then only the
// given tests will run. Marking tags as run only will by default make
// the context ignore skipped tags.
func RunOnly(tags ...string) {
	tc.RunOnly(tags...)
}

// Verbose sets the verbosity of the default context
func Verbose(verbose bool) {
	tc.Verbose = verbose
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
func Test(tag string, t T, testFn func(t T)) {
	tc.Test(tag, t, testFn)
}

// Benchmark executes a benchmark under the given tag with the
// the given benchmarking environment within the default context
func Benchmark(tag string, b B, benchmarkFn func(b B)) {
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

type skippable interface {
	Skip(...interface{})
	SkipNow()
}

type skipReason int

const (
	doNotSkip skipReason = iota
	doNotSkipFuzzy
	foundInSkip
	fuzzyMatchSkip
	notInRunOnly
)

var tc *TestContext

func init() {
	tc = New()
}
