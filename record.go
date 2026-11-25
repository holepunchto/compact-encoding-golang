package compactencoding

import (
	"cmp"
	"maps"
	"slices"
)

// Record returns a codec that writes the number of entries in a map followed
// by each key and value, encoded with key and value. Entries are written in
// the order of their keys, as `orderedRecord` in JavaScript does, so that a
// map always encodes to the same bytes. Decoding keeps the last of repeated
// keys.
func Record[K cmp.Ordered, V any](key Codec[K], value Codec[V]) Codec[map[K]V] {
	return recordCodec[K, V]{key, value}
}

type recordCodec[K cmp.Ordered, V any] struct {
	key   Codec[K]
	value Codec[V]
}

func (c recordCodec[K, V]) Preencode(state *State, value map[K]V) {
	state.End += varintSize(uint64(len(value)))

	for _, k := range slices.Sorted(maps.Keys(value)) {
		c.key.Preencode(state, k)
		c.value.Preencode(state, value[k])
	}
}

func (c recordCodec[K, V]) Encode(state *State, value map[K]V) error {
	if err := encodeVarint(state, uint64(len(value))); err != nil {
		return err
	}

	for _, k := range slices.Sorted(maps.Keys(value)) {
		if err := c.key.Encode(state, k); err != nil {
			return err
		}

		if err := c.value.Encode(state, value[k]); err != nil {
			return err
		}
	}

	return nil
}

func (c recordCodec[K, V]) Decode(state *State) (map[K]V, error) {
	n, err := decodeVarint(state)
	if err != nil {
		return nil, err
	}

	// The map is left to grow, as a hostile count could otherwise force a
	// huge allocation, while each entry takes at least one byte to read.
	value := make(map[K]V)

	for range n {
		k, err := c.key.Decode(state)
		if err != nil {
			return nil, err
		}

		v, err := c.value.Decode(state)
		if err != nil {
			return nil, err
		}

		value[k] = v
	}

	return value, nil
}
