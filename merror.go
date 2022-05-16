package merror

import (
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

func (b *Builder) Build() *Error {
	b.mu.Lock()
	defer b.mu.Unlock()
	errors := b.errors
	format := b.format
	b.errors = nil
	b.format = nil
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
