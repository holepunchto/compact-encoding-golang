# compact-encoding-golang

Compact encoding schemes for Go. The Go counterpart of [compact-encoding](https://github.com/holepunchto/compact-encoding), and wire compatible with it.

```
go get github.com/holepunchto/compact-encoding-golang/v2
```

## Usage

Encode and decode a single value with `Encode()` and `Decode()`:

```go
package main

import (
	"fmt"

	c "github.com/holepunchto/compact-encoding-golang/v2"
)

func main() {
	data, err := c.Encode(c.Uint, 42)
	if err != nil {
		panic(err)
	}

	value, err := c.Decode(c.Uint, data)
	if err != nil {
		panic(err)
	}

	fmt.Println(value) // 42
}
```

To encode several values into one buffer, use a `State` directly. `Preencode()` sums the size of each value, `Allocate()` allocates the buffer, and `Encode()` writes each value. Decoding reads them back from the start:

```go
var state c.State

c.Uint.Preencode(&state, 42)
c.String.Preencode(&state, "hi")

state.Allocate()

c.Uint.Encode(&state, 42)
c.String.Encode(&state, "hi")

state.Start = 0

n, _ := c.Uint.Decode(&state) // 42
s, _ := c.String.Decode(&state) // "hi"
```

Decoding fails with `ErrOutOfBounds` for input that ends early and leaves any bytes after the value for the next decode.

## Codecs

Every codec is a `Codec[T]`.

| Codec                                              | Go value                                         |
| -------------------------------------------------- | ------------------------------------------------ |
| `Uint`, `Uint8`, `Uint16`, `Uint32`, `Uint64`      | unsigned integers, `Uint` as a varint            |
| `Int`, `Int8`, `Int16`, `Int32`, `Int64`           | signed integers, zigzagged, `Int` as a varint    |
| `Float32`, `Float64`                               | floating point                                   |
| `Bool`                                             | `bool`                                           |
| `String`                                           | `string`, as UTF-8                               |
| `Buffer`, `OptionalBuffer`                         | `[]byte`, where `OptionalBuffer` decodes empty as `nil` |
| `Array(codec)`                                     | `[]T`                                            |
| `Record(key, value)`                               | `map[K]V`, encoded in the order of its keys      |
| `Frame(codec)`                                     | `T`, prefixed with its encoded length            |
| `JSON[T]()`                                        | `T`, as a `String` of its JSON encoding          |

Unlike in JavaScript, integers carry their full 64-bit range. Values above 2^53 - 1 decode in Go but fail to decode in JavaScript, unless decoded there with the `biguint` and `bigint` codecs.

## Custom codecs

Implement `Codec[T]` to encode your own types. `Preencode()` sizes a value, `Encode()` writes it, and `Decode()` reads it back:

```go
type Point struct {
	X, Y uint
}

type pointCodec struct{}

func (pointCodec) Preencode(state *c.State, p Point) {
	c.Uint.Preencode(state, p.X)
	c.Uint.Preencode(state, p.Y)
}

func (pointCodec) Encode(state *c.State, p Point) error {
	if err := c.Uint.Encode(state, p.X); err != nil {
		return err
	}

	return c.Uint.Encode(state, p.Y)
}

func (pointCodec) Decode(state *c.State) (p Point, err error) {
	if p.X, err = c.Uint.Decode(state); err != nil {
		return p, err
	}

	p.Y, err = c.Uint.Decode(state)

	return p, err
}
```

## Reflection

`Marshal()` and `Unmarshal()` encode a value by its type, which saves writing a codec for plain structs. A struct encodes its exported fields in order, a `[]byte` as `Buffer`, any other slice as `Array`, and a scalar as the codec named after its type:

```go
type Point struct {
	X, Y uint
}

data, err := c.Marshal(Point{X: 1, Y: 2})

var p Point

err = c.Unmarshal(data, &p)
```

Tag a field with `compact:"-"` to skip it, and implement `Marshaler` and `Unmarshaler` for a type that needs to encode itself.

## License

Apache-2.0
