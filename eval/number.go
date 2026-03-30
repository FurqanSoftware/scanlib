package eval

import "math"

func init() {
	Functions["numbers.prime"] = numberPrime
	Functions["numbers.gcd"] = numberGCD
	Functions["numbers.lcm"] = numberLCM
	Functions["numbers.coprime"] = numberCoprime
}

// numberPrime returns true if n is a prime number.
func numberPrime(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	n, ok := toInt(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	if n < 2 {
		return false, nil
	}
	if n == 2 {
		return true, nil
	}
	if n%2 == 0 {
		return false, nil
	}
	for i := 3; i <= int(math.Sqrt(float64(n))); i += 2 {
		if n%i == 0 {
			return false, nil
		}
	}
	return true, nil
}

// numberGCD returns the greatest common divisor of a and b.
func numberGCD(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	a, ok := toInt(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	b, ok := toInt(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return gcd(a, b), nil
}

// numberLCM returns the least common multiple of a and b.
func numberLCM(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	a, ok := toInt(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	b, ok := toInt(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return a / gcd(a, b) * b, nil
}

// numberCoprime returns true if a and b are coprime.
func numberCoprime(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	a, ok := toInt(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	b, ok := toInt(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return gcd(a, b) == 1, nil
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
