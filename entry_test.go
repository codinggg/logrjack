package logrjack

import (
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestAddCallstack(t *testing.T) {
	// arrange
	want := `stacktrace="github.com/judwhite/logrjack/entry_test.go:TestAddCallstack:17"`
	e := NewEntry()

	// act
	e.AddError(errors.WithStack(io.EOF))

	// assert
	actual := e.String()
	if !strings.Contains(actual, want) {
		t.Fatalf("could not find expected callstack string: '%s'\ngot:\n\t%s", want, actual)
	}
}
