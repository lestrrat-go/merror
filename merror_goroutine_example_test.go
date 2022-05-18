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
			defer wg.Done()
			defer merror.AddToContext(ctx, &err)

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
