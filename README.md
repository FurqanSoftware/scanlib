# Scanlib

Use Scanspec to validate input files and generate input parsing code in C99, C++14, Go, and Python 3.

## Scanspec

Learn more about Scanspec [here](https://help.toph.co/drafts/scanspec).

### Examples

#### Two Integers

```
3 2
```

```
var A, B int
scan A, B
check A >= 0, A < 10, B >= 0, B < 20
eol
eof
```

``` c
#include <stdio.h>

int main() {
	int A, B;
	scanf("%d %d", &A, &B);

	return 0;
}
```

``` cpp
#include <iostream>

using namespace std;

int main() {
	int A, B;
	cin >> A >> B;

	return 0;
}
```

``` go
package main

import "fmt"

func main() {
	var A, B int
	fmt.Scan(&A, &B)

}
```

``` py
A, B = map(int, input().split())
```

#### R, C, and Grid

```
3 5
**...
..*..
....*
```

```
var R, C int
scan R, C
check R >= 1, R < 25, C >= 1, C < 25
eol
var G [R]string
for i := 0 ... R
	scan G[i]
	check len(G[i]) == C
	check regexp.match(G[i], "^[*.]+$")
	eol
end
eof
```

``` cpp
#include <iostream>
#include <string>

using namespace std;

int main() {
	int R, C;
	cin >> R >> C;
	string G[R];
	for (int i = 0; i < R; ++i) {
		cin >> G[i];
	}

	return 0;
}
```

``` py
R, C = map(int, input().split())
G = [""] * R
for i in range(0, R):
	G[i] = input()
```

### Specification

#### Comments

A comment begins with the # character, and ends at the end of the line. A comment cannot begin within a string literal.

```
# This is a comment
```

#### Keywords

```
check else end eof eol for if scan scanln var
```

#### Types

```
bool
int
int64
float32
float64
string
[]T
```

#### Operators

```
+  -  *  /  **
==  !=  <  >  <=  >=
&&  ||
```

The `**` operator performs exponentiation.

#### Check Statements

```
check n > 0, n < 1000
check e > 0, f < 5.0
check n ** 2 < 1000000
```

#### Variable Declarations

```
var n int
var e, f float64
var a, b string
var G [R]string
```

#### Scan Statements

```
scan n
scan a, e
scan e, f, n
scan G[2]
```

#### Scanln Statements

```
scanln S
scanln A, B
```

#### If Statements

```
if q == 1
	scan i, G[i]
else if q == 2
	scan a
else if q == 3
	scan l, h
else
	scan G[q]
end
```

#### Assignment Statements

```
s = s + A
G[i] = G[i] + 1
```

#### For Range Statements

```
for i := 0 ... n
	scan G[i]
end
```

#### For Scan/Scanln Statements

```
var s string
var i int
for scan s, i
	# Scans s and i repeatedly until EOF
end
```

```
var s string
for scanln s
	# Scans line to s repeatedly until EOF
end
```

#### EOL Statements

The following indicates end of line.

```
eol
```
#### EOF Statements

The following indicates end of file.

```
eof
```

#### Built-in Functions

```
len(a): Returns the length of array a.
toInt64(s, b=10): Parses string s in base b and returns in int64.
```

#### Modules

Built-in modules provide additional functions to make validation of input data easier.

##### arrays

Array validation functions.

```
arrays.sorted(A): Returns true if the array is sorted in non-decreasing order.
arrays.distinct(A): Returns true if all elements in the array are unique.
arrays.permutation(N, A): Returns true if A is a permutation of 1..N.
arrays.range(A, lo, hi): Returns true if all elements are in [lo, hi].
```

##### math

Mathematical functions.

```
math.abs(n): Returns the absolute value of n.
math.min(a, b): Returns the minimum of a and b.
math.max(a, b): Returns the maximum of a and b.
math.sqrt(n): Returns the square root of n.
math.log2(n): Returns the base-2 logarithm of n.
math.ceil(n): Returns the smallest integer not less than n.
math.floor(n): Returns the largest integer not greater than n.
math.pow(n, e): Returns n raised to the power of e.
math.sum(a...): Returns the sum of the arguments.
```

##### graphs

Graph check functions take node count N (1-indexed) and two parallel edge arrays U and V.

```
graphs.simple(N, U, V): Returns true if the graph has no self-loops or duplicate edges.
graphs.connected(N, U, V): Returns true if all N nodes are connected.
graphs.acyclic(N, U, V): Returns true if the graph has no cycles.
graphs.tree(N, U, V): Returns true if the edges form a tree on N nodes.
```

Example:

```
var N int
scan N
check N >= 1, N <= 100000
eol
var U [N-1]int
var V [N-1]int
for i := 0 ... N-1
	scan U[i], V[i]
	check U[i] >= 1, U[i] <= N
	check V[i] >= 1, V[i] <= N
	eol
end
check graphs.tree(N, U, V)
eof
```

##### numbers

Number theory functions for validating numeric properties.

```
numbers.prime(n): Returns true if n is a prime number.
numbers.gcd(a, b): Returns the greatest common divisor of a and b.
numbers.lcm(a, b): Returns the least common multiple of a and b.
numbers.coprime(a, b): Returns true if a and b are coprime.
```

##### strings

String validation functions.

```
strings.distinct(s): Returns true if all characters in s are unique.
strings.sorted(s): Returns true if s is sorted in non-decreasing order.
strings.palindrome(s): Returns true if s is a palindrome.
strings.lowercase(s): Returns true if s contains only lowercase letters.
strings.uppercase(s): Returns true if s contains only uppercase letters.
strings.alpha(s): Returns true if s contains only letters.
strings.digit(s): Returns true if s contains only digits.
strings.alphanumeric(s): Returns true if s contains only letters and digits.
strings.binary(s): Returns true if s contains only '0' and '1'.
strings.contains(s, sub): Returns true if s contains the substring sub.
```

##### regexp

Regular expression functions.

```
regexp.match(s, pattern): Returns true if s matches the pattern.
regexp.fullmatch(s, pattern): Returns true if the entire string s matches the pattern.
regexp.count(s, pattern): Returns the number of non-overlapping matches of pattern in s.
```

##### matrices

Matrix/grid validation functions.

```
matrices.dimensions(G, R, C): Returns true if string array G has R elements, each of length C.
```

## CLI

```
go install git.furqansoftware.net/toph/scanlib/cmd/scanlib@latest
```

Evaluate a Scanspec file against an input file:

```
scanlib scanspec input.txt
```

Or read input from stdin:

```
echo "3 2" | scanlib scanspec
```

## Library Usage

The evaluator supports timeout and cancellation via Go's `context.Context`:

``` go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
values, err := eval.Evaluate(source, input, eval.WithContext(ctx))
```
