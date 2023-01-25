package eval

type Option interface {
	apply(*evaluator)
}

type optionFunc func(*evaluator)

func (f optionFunc) apply(e *evaluator) {
	f(e)
}

func ScannerBuffer(buf []byte, max int) Option {
	return optionFunc(func(e *evaluator) {
		e.Input.sc.Buffer(buf, max)
	})
}
