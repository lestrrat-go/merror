package merror_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/lestrra-go/merror"
)

func ExampleMerror_NoErrors() {
	b := merror.NewBuilder()
	ctx := b.NewContext(context.Background())

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(ctx context.Context) (err error) {
			defer merror.AddToContext(ctx, &err)
			defer wg.Done()

			// No errors!
			return
		}(ctx)
	}

	wg.Wait()

	if errs := b.Build(); errs != nil {
		fmt.Printf("%s\n", errs)
	}

	// OUTPUT:
}
