// Package gen provides shared types and errors for Scanspec code generators.
package gen

import "errors"

// ErrUnsupportedType is returned when a generator encounters a Scanspec type
// it cannot represent in the target language.
var ErrUnsupportedType = errors.New("unsupported type")
