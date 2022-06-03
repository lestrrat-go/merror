// Package merror provides an error type that is suitable for representing
// multiple errors.
package merror

import (
	"context"
	"strings"
	"sync"
)

// DefaultFormat is the Formatter that will be used to stringify
// all merror.Error instances.
var DefaultFormat Formatter = NewFormatBuilder().Marker("✔ ").Indent("  ").MustBuild()

// Error is the main error object that holds multiple errors
type Error struct {
	format Formatter
	errors []error
}

func (err Error) Error() string {
	format := err.format
	if format == nil {
		format = DefaultFormat
	}
	return format.Format(&err)
}

func (err Error) Errors() []error {
	return err.errors
}

// Builder is used to create an Error instance.
type Builder struct {
	mu     sync.Mutex
	errors []error
	format Formatter
}

func NewBuilder() *Builder {
	return &Builder{}
}

// Build returns the error that it has accumulated. If no errors have been
// specified using `Error()` method, then this method return `nil`, indicating
// that there weren't any errors.
func (b *Builder) Build() *Error {
	b.mu.Lock()
	defer b.mu.Unlock()
	errors := b.errors
	format := b.format
	b.errors = nil
	b.format = nil

	// If there are no errors, there is no error!
	if len(errors) == 0 {
		return nil
	}

	return &Error{
		errors: errors,
		format: format,
	}
}

func (b *Builder) Error(err error) *Builder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.errors = append(b.errors, err)
	return b
}

func (b *Builder) Formatter(f Formatter) *Builder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.format = f
	return b
}

type identMerrorContext struct{}

// NewContext creates a new context.Context that has the `Builder`
// b associated with it. You can use `merror.AddToContext` to
// add more errors through the context
func (b *Builder) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, identMerrorContext{}, b)
}

// If the `context.Context` object does not already contain a
// `merror.Builder` in it, this function is a no-op.
//
// The second argument must be a pointer, as this function is
// intended to be used in a defer statement.
//
//   func Foo(ctx context.Context) (err error) {
//     defer merror.AddToContext(ctx, err) // bound to the value of err NOW (i.e. nil)
//     return fmt.Errorf(`foo`)
//   }
//
//   func Foo(ctx context.Context) (err error) {
//     defer merror.AddToContext(ctx, &err) // bound to the pointer, so can detect assignments to it
//     return fmt.Errorf(`foo`)
//   }
func AddToContext(ctx context.Context, ptr *error) {
	if ptr == nil {
		return
	}
	err := *ptr
	if err == nil {
		return
	}

	v := ctx.Value(identMerrorContext{})
	if v == nil {
		return
	}
	b, ok := v.(*Builder)
	if !ok {
		return
	}
	b.Error(err)
}

// FormatFunc is a Formatter that is represented as a function
type FormatFunc func(*Error) string

func (fn FormatFunc) Format(err *Error) string {
	return fn(err)
}

// Formatter is an object that is responsible for stringifying
// the *Error object. You may create your custom Formatter object,
// or you can use the existing FormatBuilder to configure the
// default Formatter
type Formatter interface {
	Format(*Error) string
}

// FormatBuilder is the object that builds an instance of the
// default Formatter.
type FormatBuilder struct {
	mu      sync.RWMutex
	marker  string
	message string
	indent  string
}

const (
	defaultMarker  = "- "
	defaultIndent  = "  "
	defaultMessage = "errors found:"
)

// Creates a new `FormatBuilder` to build `Format` objects
func NewFormatBuilder() *FormatBuilder {
	return &FormatBuilder{
		marker:  defaultMarker,
		indent:  defaultIndent,
		message: defaultMessage,
	}
}

func (b *FormatBuilder) MustBuild() Formatter {
	f, err := b.Build()
	if err != nil {
		panic(err)
	}
	return f
}

func (b *FormatBuilder) Build() (Formatter, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	indent := b.indent
	marker := b.marker
	message := b.message
	b.indent = defaultIndent
	b.marker = defaultMarker
	b.message = defaultMessage
	return &formatter{
		indent:  indent,
		marker:  marker,
		message: message,
	}, nil
}

func (b *FormatBuilder) Marker(s string) *FormatBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.marker = s
	return b
}

func (b *FormatBuilder) Message(s string) *FormatBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.message = s
	return b
}

func (b *FormatBuilder) Indent(s string) *FormatBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.indent = s
	return b
}

type formatter struct {
	marker  string
	indent  string
	message string
}

func (f *formatter) Format(err *Error) string {
	var b strings.Builder
	b.WriteString(f.message)
	for _, suberr := range err.Errors() {
		b.WriteRune('\n')
		b.WriteString(f.indent)
		b.WriteString(f.marker)
		b.WriteString(suberr.Error())
	}
	return b.String()
}
