package eval

// Option configures the evaluator.
type Option interface {
	apply(*evaluator)
}

type optionFunc func(*evaluator)

func (f optionFunc) apply(e *evaluator) {
	f(e)
}

// ScannerBuffer returns an Option that sets the initial buffer and maximum
// buffer size for the input scanner.
func ScannerBuffer(buf []byte, max int) Option {
	return optionFunc(func(e *evaluator) {
		e.Input.sc.Buffer(buf, max)
	})
}
