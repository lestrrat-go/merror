package merror_test

import (
	"fmt"
	"regexp"
	"strings"
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
