package merror_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/lestrra-go/merror"
	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	b := merror.NewBuilder()
	errors := []error{
		fmt.Errorf(`1`),
		fmt.Errorf(`2`),
		fmt.Errorf(`3`),
	}
	for _, err := range errors {
		b.Error(err)
	}

	err := b.Build()
	require.Error(t, err, `b.Build() should create a prorper error`)

	msg := err.Error()
	require.Regexp(t, regexp.MustCompile(`^errors found:`), msg, `message matches`)
	require.Contains(t, msg, "\n  ✔ 1\n  ✔ 2\n  ✔ 3")
}

func TestFormatter(t *testing.T) {
	b := merror.NewBuilder()
	errors := []error{
		fmt.Errorf(`1`),
		fmt.Errorf(`2`),
		fmt.Errorf(`3`),
	}
	for _, err := range errors {
		b.Error(err)
	}

	b.Formatter(merror.FormatFunc(func(err *merror.Error) string {
		var b strings.Builder
		b.WriteString(`{"errors":[`)
		for i, e := range err.Errors() {
			if i > 0 {
				b.WriteRune(',')
			}
			fmt.Fprintf(&b, "%q", e.Error())
		}
		b.WriteString(`]}`)
		return b.String()
	}))

	err := b.Build()
	require.Error(t, err, `b.Build() should create a prorper error`)

	msg := err.Error()
	require.Equal(t, msg, `{"errors":["1","2","3"]}`)
}

func TestGoroutine(t *testing.T) {
	b := merror.NewBuilder()
	m := make(map[string]struct{})
	var wg sync.WaitGroup
	wg.Add(10)
	const max = 10
	for i := 0; i < max; i++ {
		i := i
		m[fmt.Sprintf(`%d`, i)] = struct{}{}
		go func(ctx context.Context) (err error) {
			defer merror.AddToContext(ctx, &err)
			defer wg.Done()

			return fmt.Errorf(`%d`, i)
		}(b.NewContext(context.Background()))
	}

	wg.Wait()

	merr := b.Build()
	for _, err := range merr.Errors() {
		t.Logf("m = %#v", m)
		t.Logf("err = %s", err)
		delete(m, err.Error())
	}

	require.Len(t, m, 0, `expected m to contain zero entries`)
}
