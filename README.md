merrors
=======

Simple multi-error `error` type for Go.

tl;dr:

* Sometimes you want multiple errors to be bundled into a single error.
* Whereas other libraries directly act on the error object, this package uses a builder to create errors. This significanly reduces the complexity of the error object itself.

# DESCRIPTION

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
