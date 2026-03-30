// Package gen provides shared types and errors for Scanspec code generators.
package gen

import (
	"errors"
	"fmt"
)

// ErrUnsupported is the base error for unsupported features in a generator.
var ErrUnsupported = errors.New("unsupported")

// ErrUnsupportedType is returned when a generator encounters a Scanspec type
// it cannot represent in the target language.
var ErrUnsupportedType = fmt.Errorf("%w type", ErrUnsupported)

// ErrUnsupportedFunction is returned when a generator encounters a function
// call it cannot represent in the target language.
var ErrUnsupportedFunction = fmt.Errorf("%w function", ErrUnsupported)
