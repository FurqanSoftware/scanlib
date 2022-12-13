# Scanlib

Use Scanspec Validate input files and generate input code in different programming languages.

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

``` cpp
// Generated using Scanlib

#include <iostream>

using namespace std;

int main() {
	int A, B;
	cin >> A;
	cin >> B;
	
	return 0;
}
```

``` py
# Generated using Scanlib

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
	check re(G[i], "^[*.]+$")
	eol
end
eof
```

``` cpp
// Generated using Scanlib

#include <iostream>
#include <string>

using namespace std;

int main() {
	int R, C;
	cin >> R;
	cin >> C;
	string G[R];
	for (int i = 0; i < R; ++i) {
		cin >> G[i];
	}
	
	return 0;
}
```

``` py
# Generated using Scanlib

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
check eof eol for scan var
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

#### Check Statements

```
check n > 0, n < 1000
check e > 0, f < 5.0
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
#### EOL Statements

The following indicates end of file.

```
eof
```

#### Built-in Functions

```
len(a): Returns the length of array a.
re(s, x): Returns true if string s matches regular expression x.
pow(n, e): Returns n raised to the power of e. Result is int or int64 if both n and e are int or int64, otherwise float64.
```

## TODO

- [x] If Statements
- [ ] C Generator
- [x] Go Generator
- [ ] Graph Checks
- [ ] CLI Tool
- and more...