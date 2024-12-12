# TDC
TDC, abbreviated to **T**errible, **D**igit **C**ompressor, is a
numerical digit compressor. As the name suggests, the algorithm
specializes in compressing digits to the smallest possible size,
to save state.

In the case of `dist2`, without TDC, saving $n$ digits would use
$n$ bytes, which is obviously inefficient and takes up lots of
storage for say, a billion digits would take up one gigabyte of
storage.

That may not seem like a lot, but that compounds quickly.

## Why?
To save space.

# Usage/Functions

## TDCEncodeString
Encodes a string with the TDC algorithm, returning the encoded
value as a byte array.
```go
func TDCEncodeString(data string) []byte
```

## TDCDecodeString
Returns the string of original bytes from a byte array.
```go
func TDCDecodeString(encode []byte) string
```