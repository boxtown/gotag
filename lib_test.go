package gotag

import "testing"

func TestSkip(t *testing.T) {
	tc := New()
	tc.Skip("tagA", "tagB")

	mock := &mockT{}
	tc.Test("tagA", mock, func(t T) {})
	tc.Test("tagB", mock, func(t T) {})
	tc.Test("taga", mock, func(t T) {})
	if mock.skipped != 2 {
		t.Error("Wrong number of tests skipped")
		t.Fail()
	}
}

func TestSkipFuzzy(t *testing.T) {
	tc := New()
	tc.Skip("tagA")
	tc.Fuzzy = true

	mock := &mockT{}
	tc.Test("tagA", mock, func(t T) {})
	tc.Test("taga", mock, func(t T) {})
	if mock.skipped != 2 {
		t.Error("Wrong number of tests skipped")
		t.Fail()
	}
}

func TestRunOnly(t *testing.T) {
	tc := New()
	tc.RunOnly("tagA")

	mock := &mockT{}
	tc.Test("tagA", mock, func(t T) {})
	tc.Test("tagB", mock, func(t T) {})
	tc.Test("taga", mock, func(t T) {})
	if mock.skipped != 2 {
		t.Error("Wrong number of tests skipped")
		t.Fail()
	}
}

func TestRunOnlyFuzzy(t *testing.T) {
	tc := New()
	tc.RunOnly("tagA")
	tc.Fuzzy = true

	mock := &mockT{}
	tc.Test("tagA", mock, func(t T) {})
	tc.Test("taga", mock, func(t T) {})
	if mock.skipped != 0 {
		t.Error("Wrong number of tests skipped")
		t.Fail()
	}
}

func TestRunOverridesSkip(t *testing.T) {
	tc := New()
	tc.Skip("tagA")
	tc.RunOnly("tagA")

	mock := &mockT{}
	tc.Test("tagA", mock, func(t T) {})
	if mock.skipped != 0 {
		t.Error("Wrong number of tests skipped")
		t.Fail()
	}

	tc.Fuzzy = true
	tc.Test("taga", mock, func(t T) {})
	if mock.skipped != 0 {
		t.Error("Wrong number of tests skipped")
		t.Fail()
	}
}

type mockT struct {
	skipped int
}

func (t *mockT) Error(...interface{})              {}
func (t *mockT) Errorf(string, ...interface{})     {}
func (t *mockT) Fail()                             {}
func (t *mockT) FailNow()                          {}
func (t *mockT) Failed() bool                      { return false }
func (t *mockT) Fatal(...interface{})              {}
func (t *mockT) Fatalf(string, ...interface{})     {}
func (t *mockT) Log(...interface{})                {}
func (t *mockT) Logf(string, ...interface{})       {}
func (t *mockT) Parallel()                         {}
func (t *mockT) Run(string, func(*testing.T)) bool { return false }
func (t *mockT) Skip(...interface{})               { t.skipped++ }
func (t *mockT) SkipNow()                          { t.skipped++ }
func (t *mockT) Skipf(string, ...interface{})      { t.skipped++ }
func (t *mockT) Skipped() bool                     { return false }
