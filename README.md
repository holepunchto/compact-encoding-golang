# compact-encoding (Go)

> 🚧 **Work in Progress** - A Go port of [compact-encoding](https://github.com/compact-encoding/compact-encoding)

A compact binary encoding library for Bare/Pear interop.

## Overview

This is a Go implementation of the compact-encoding format, providing efficient binary serialization with minimal overhead. The library uses a two-pass encoding strategy: a preencode pass to calculate buffer size, followed by the actual encoding.

## Installation
```bash
go get github.com/yourusername/compact-encoding
```

## Features

- ✅ Variable-length integer encoding (uint, int)
- ✅ Fixed-size integer types (uint8, uint16, uint32, uint64, int8, int16, int32, int64)
- ✅ String encoding
- ✅ Buffer/byte slice encoding
- ✅ Boolean encoding
- ✅ Array encoding with generic element types
- ✅ Generic `Decode`/`Encode` helpers

## Usage

### Basic Example
```go
package main

import (
    "fmt"
    enc "github.com/yourusername/compact-encoding"
)

func main() {
    // Using the helper functions
    encoder := enc.NewInt()
    
    // Encode
    buf, err := enc.Encode(encoder, 42)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %v\n", buf)
    
    // Decode
    value, err := enc.Decode(encoder, buf)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %d\n", value)
}
```

### Manually
```go
state := enc.NewState()
encoder := enc.NewString()

// Preencode to calculate size
encoder.Preencode(state, "hello")
state.Allocate() // Create the buffer

// Encode
err := encoder.Encode(state, "hello")
if err != nil {
    panic(err)
}

// Decode
state.Rewind()
value, err := encoder.Decode(state)
```

### Array Encoding
```go
// Array of strings
encoder := enc.NewArray(enc.NewString())
buf, err := enc.Encode(encoder, []string{"hello", "world"})

// Array of uints
uintEncoder := enc.NewArray(enc.NewUint())
buf, err := enc.Encode(uintEncoder, []uint{1, 2, 3})
```

## API

### Encoders

- `NewUint()` - Variable-length unsigned integer
- `NewInt()` - Variable-length signed integer (zigzag encoding)
- `NewUint8()`, `NewUint16()`, `NewUint32()`, `NewUint64()` - Fixed-size unsigned integers
- `NewInt8()`, `NewInt16()`, `NewInt32()`, `NewInt64()` - Fixed-size signed integers
- `NewString()` - UTF-8 string encoding
- `NewBuffer()` - Byte slice encoding
- `NewBool()` - Boolean encoding
- `NewArray(elementEncoder)` - Generic array encoding

### Helper Functions
```go
func Encode[T any](encoder Encoder[T], value T) ([]byte, error)
func Decode[T any](encoder Encoder[T], data []byte) (T, error)
```

## Development

Run tests:
```bash
go test .
```

## Compatibility

This implementation aims to be wire-compatible with the JavaScript [compact-encoding](https://github.com/compact-encoding/compact-encoding) library.

## License

Apache-2.0

## Related

- [compact-encoding](https://github.com/compact-encoding/compact-encoding) - Original JavaScript implementation
