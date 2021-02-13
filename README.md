# Scanlib

Validate input files and generate input code in different programming languages.

## Specification

### Keywords

```
check eof eol for scan var
```

### Check Statements

```
check n > 0, n < 1000
check e > 0, f < 5.0
```

### Variable Declarations

```
var n int
var e, f float64
var a, b string
var G [R]string
```

### Scan Statements

```
scan n
scan a, e
scan e, f, n
scan G[2]
```

### For Statements

```
for i 0 n
	scan G[i]
end
```

### EOL Statements

The following indicates end of line.

```
eol
```
### EOL Statements

The following indicates end of file.

```
eof
```

## TODO

- [ ] If Statements
- [ ] C Generator
- [ ] Go Generator
- [ ] CLI Tool
- and more...