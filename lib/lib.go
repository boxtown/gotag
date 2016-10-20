package tflag

import (
	"strings"
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

	Verbose bool
}

// New constructs a new instance of TestFlags
func New() *TestFlags {
	return &TestFlags{skip: make(map[string]bool)}
}

// Skip marks a test flag to be skipped when testing
// within the context of the TestFlags instance
func (tf *TestFlags) Skip(flag string) {
	tf.skip[strings.ToLower(flag)] = true
}

// Test executes a test under the given flag with the given testing environment
// within the context of the TestFlags instance
func (tf *TestFlags) Test(flag string, t *testing.T, testFn func(t *testing.T)) {
	if _, ok := tf.skip[strings.ToLower(flag)]; ok {
		t.SkipNow()
	} else {
		testFn(t)
	}
}

// Benchmark executes a benchmark under the given flag with the given benchmarking
// environment within the context of the TestFlags instance
func (tf *TestFlags) Benchmark(flag string, b *testing.B, benchmarkFn func(b *testing.B)) {
	if _, ok := tf.skip[strings.ToLower(flag)]; ok {
		b.SkipNow()
	} else {
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
