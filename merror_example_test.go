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
