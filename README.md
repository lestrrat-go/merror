merrors
=======

Simple multi-error `error` type for Go.

tl;dr:

* Sometimes you want multiple errors to be bundled into a single error.
* Whereas other libraries directly act on the error object, this package uses a builder to create errors. This significanly reduces the complexity of the error object itself.

# DESCRIPTION

Simple usage:

<!-- INCLUDE(merror_example_test.go) -->
```go
package merror_test

import (
  "fmt"

  "github.com/lestrra-go/merror"
)

func ExampleMerror() {
  err1 := fmt.Errorf(`first error`)
  err2 := fmt.Errorf(`second error`)
  err3 := fmt.Errorf(`third error`)

  err := merror.NewBuilder().
    Error(err1).
    Error(err2).
    Error(err3).
    Build()

  for _, suberr := range err.Errors() {
    fmt.Printf("%s\n", suberr.Error())
  }

  // OUTPUT:
  // first error
  // second error
  // third error
}
```
source: [merror_example_test.go](https://github.com/lestrrat-go/merror/blob/main/merror_example_test.go)
<!-- END INCLUDE -->

Use with multiple goroutines:

<!-- INCLUDE(merror_goroutine_example_test.go) -->
```go
package merror_test

import (
  "context"
  "fmt"
  "sort"
  "sync"

  "github.com/lestrra-go/merror"
)

func ExampleMerror_MultipleGoroutines() {
  b := merror.NewBuilder()
  ctx := b.NewContext(context.Background())

  var wg sync.WaitGroup
  wg.Add(10)

  for i := 0; i < 10; i++ {
    i := i
    go func(ctx context.Context) (err error) {
      defer merror.AddToContext(ctx, &err)
      defer wg.Done()

      return fmt.Errorf(`%d`, i)
    }(ctx)
  }

  wg.Wait()

  // note: in order to make the the output deterministic,
  // we're having to sort the errors
  errs := b.Build().Errors()
  sort.Slice(errs, func(i, j int) bool {
    return errs[i].Error() < errs[j].Error()
  })
  for _, err := range errs {
    fmt.Printf("%s\n", err)
  }

  // OUTPUT:
  // 0
  // 1
  // 2
  // 3
  // 4
  // 5
  // 6
  // 7
  // 8
  // 9
}
```
source: [merror_goroutine_example_test.go](https://github.com/lestrrat-go/merror/blob/main/merror_goroutine_example_test.go)
<!-- END INCLUDE -->

