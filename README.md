# Gotag

[![GoDoc](https://godoc.org/github.com/boxtown/gotag?status.svg)](https://godoc.org/github.com/boxtown/gotag)


**Gotag** is a testing utility tool that makes it easy for you to selectively skip/run tests in Go. If you ever needed to mark a suite
of integration tests to be skipped, then **Gotag** is the tool for the job. 

# Contents
[Usage](#usage)
[Selectively running tests](#selectively-running-tests)
[Tags](#tags)
[Roadmap](#roadmap)

## Usage

Simply run
```
go get github.com/boxtown/gotag
```
to install **Gotag**.  
  
To use **Gotag**, configure the test context in either `init` or `TestMain` and then wrap your tests inside  
`Test` or `Benchmark` like so:  

```Go
import (
  "fmt"
  "testing"
   "github.com/boxtown/gotag"
)

func TestMain(m *testing.Main) {
  gotag.Skip(gotag.Integration)
  os.exit(m.Run())
}

// This test will not run
func TestSomethingIntegrated(t *testing.T) {
  gotag.Test(gotag.Integration, t, func(t gotag.T) {
    t.FailNow()
  })
}

// This test will
func TestSomethingElse(t *testing.T) {
  gotag.Test("something else", t, func(t gotag.T) {
    fmt.Println("I'm running inside the Gotag context!")
  })
}

// This test will also run
func TestSomethingBasic(t *testing.T) {
  fmt.Println("Gotag has no knowledge of me!")
}
```

## Selectively running tests

You can also choose to run only certain tags. Note that by calling RunOnly skip is ignored

```Go
import (
  "fmt"
  "testing"
  "github.com/boxtown/gotag"
)

func TestMain(m *testing.M) {
  gotag.RunOnly("tagA", "tagB")
  os.Exit(m.Run())
}

// Does not get run because tag is not marked by RunOnly
func TestSomethingIntegrated(t *testing.T) {
  gotag.Test(gotag.Integrated, t, func(t gotag.T) {
    t.FailNow()
  })
}

// This will run
func TestTagA(t *testing.T) {
  gotag.Test("tagA", t, func(t gotag.T) {
    fmt.Println("I'm tagA!")
  })
}

// This will also run
func BenchmarkTagB(b *testing.B) {
  gotag.Benchmark("tagB", b, func(b gotag.B) {
    fmt.Println("I'm tagB!")
  })
}
```

## Tags

Tags are just simple strings. By default, **Gotag** does a strict match when checking for skip/run tags.
**Gotag** can be configured however, to do a fuzzy matching on tags

```Go
import (
  "fmt"
  "testing"
  "github.com/boxtown/gotag"
)

func TestMain(m *testing.M) {
  gotag.Skip("tagA")
  gotag.Fuzzy(true)
  os.Exit(m.Run())
}

// Skipped
func TestTagA(t *testing.T) {
  gotag.Test("tagA", t, func(t gotag.T) {
    t.FailNow()
  })
}

// Also skipped because of fuzzy matching
func TestTaga(t *testing.T) {
  gotag.Test("taga", t, func(t gotag.T) {
    t.FailNow()
  })
}
```

In the above example, the second test runs because it is within an edit distance of 2 (the default) of the registered tag.
The edit distance can be configured as well

```Go
gotag.Distance(5)
```

## Roadmap

- Load configuration from configuration files
- Hooks for Before/After test logic
