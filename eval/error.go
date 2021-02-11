package eval

type ErrUndefined struct {
	Name string
}

func (e ErrUndefined) Error() string {
	return "undefined: " + e.Name
}

type ErrCantScanType struct{}

func (e ErrCantScanType) Error() string {
	return "can't scan type"
}

type ErrInvalidArgument struct{}

func (e ErrInvalidArgument) Error() string {
	return "invalid argument"
}

type ErrInvalidOperation struct{}

func (e ErrInvalidOperation) Error() string {
	return "invalid operation"
}

type ErrNonIntegerIndex struct{}

func (e ErrNonIntegerIndex) Error() string {
	return "non-integer index"
}

type ErrCheckError struct {
}

func (e ErrCheckError) Error() string {
	return "check error"
}
