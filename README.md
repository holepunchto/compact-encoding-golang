# compact-encoding (Go)

> 🚧 **Work in Progress** - A Go port of [compact-encoding](https://github.com/holepunchto/compact-encoding)

A compact binary encoding library for Bare/Pear interop.

## Overview

This is a Go implementation of the compact-encoding format, providing efficient binary serialization with minimal overhead. The library uses a two-pass encoding strategy: a preencode pass to calculate buffer size, followed by the actual encoding.

## Installation
```bash
go get github.com/holepunchto/compact-encoding-golang
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
    enc "github.com/holepunchto/compact-encoding-golang"
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

This implementation is wire-compatible with the JavaScript [compact-encoding](https://github.com/holepunchto/compact-encoding) library, last verified against **v3.5.0**.

`testdata/vectors.json` holds byte vectors produced by the JS library for every codec ported here. `go test` replays them in both directions (Go encode must match JS bytes, Go decode must read JS bytes back). Regenerate them from a checkout of the JS library with its dependencies installed:

```bash
node testdata/gen-vectors.js /path/to/compact-encoding > testdata/vectors.json
```

Behaviour shared with the JS library:

- `buffer` follows the JS 3.x semantics: a `nil` slice encodes as length `0` and an empty payload decodes to an empty (non-nil) slice, never `nil`. The JS `optionalBuffer` codec, which decodes an empty payload as `null`, is not ported.
- `array` refuses to decode more than `MaxArrayLength` (`0x100000`) elements, returning `EncodingErrorArrayTooBig`, exactly as the JS decoder throws `Array is too big`.
- Truncated or oversized input returns `EncodingErrorOutOfBounds` rather than panicking.

Differences to be aware of:

- **Integer range.** JS numbers are IEEE doubles, so the JS library rejects `uint`/`uint64` values above `2^53 - 1` and `int`/`int64` values outside `±2^52` (it throws on both encode and decode). Go encodes the full 64-bit range. Values outside the JS-safe range will round-trip between Go peers but fail to decode in JS; use the JS `biguint64`/`bigint64`/`biguint`/`bigint` codecs on that side, which share the fixed 8-byte wire format with Go `uint64`/`int64`.
- **Unported codecs.** The JS library also ships `uint24`, `uint40`, `uint48`, `uint56`, `uint32be`, `uint64be` and their `int` counterparts, `biguint`/`bigint`, `lexint`, `float32`/`float64`, `optionalBuffer`, `binary`, `arraybuffer`, `bitarray`, typed arrays, `ascii`/`hex`/`base64`/`utf16le` strings, `fixed(n)` and `fixed8..64`, `frame`, `date`, `json`/`ndjson`, `none`, `any`, `record`/`stringRecord`, the `ip`/`port`/`*Address` network codecs, and the `raw` (unframed) variants. None of these exist in the Go port yet.

## License

Apache-2.0

## Related

- [compact-encoding](https://github.com/holepunchto/compact-encoding) - Original JavaScript implementation
