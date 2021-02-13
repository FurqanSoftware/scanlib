# Scanlib

Use Scanspec Validate input files and generate input code in different programming languages.

## Scanspec

Learn more about Scanspec [here](https://help.toph.co/drafts/scanspec).

### Specification

#### Keywords

```
check eof eol for scan var
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

#### For Statements

```
for i 0 n
	scan G[i]
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
```

## TODO

- [ ] If Statements
- [ ] C Generator
- [ ] Go Generator
- [ ] Graph Checks
- [ ] CLI Tool
- and more...