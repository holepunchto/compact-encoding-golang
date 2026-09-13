// Generates wire vectors from the JavaScript compact-encoding library so the
// Go port can be checked byte-for-byte against it.
//
// Usage (from a checkout of https://github.com/holepunchto/compact-encoding
// with its dependencies installed):
//
//   node testdata/gen-vectors.js /path/to/compact-encoding > testdata/vectors.json

const path = require('path')
const c = require(path.resolve(process.argv[2] || 'compact-encoding'))
const b4a = require(path.join(path.resolve(process.argv[2] || 'compact-encoding'), 'node_modules/b4a'))

const MAX_SAFE_UINT = Number.MAX_SAFE_INTEGER // 2^53 - 1
const MAX_SAFE_INT = 2 ** 52 - 1
const MIN_SAFE_INT = -(2 ** 52)

const out = { version: require(path.resolve(process.argv[2] || 'compact-encoding', 'package.json')).version, vectors: [] }

function add (codec, enc, value, jsonValue = value) {
  out.vectors.push({ codec, value: jsonValue, hex: b4a.toString(c.encode(enc, value), 'hex') })
}

for (const n of [0, 1, 0xfc, 0xfd, 0xfe, 0xff, 0x100, 0xffff, 0x10000, 0xffffffff, 0x100000000, MAX_SAFE_UINT]) add('uint', c.uint, n)
for (const n of [0, 1, -1, 2, -2, 126, -126, 127, -127, 128, -128, 0x7e, -0x7f, 0x7fff, -0x8000, 0x7fffffff, -0x80000000, MAX_SAFE_INT, MIN_SAFE_INT]) add('int', c.int, n)

for (const n of [0, 1, 0xfc, 0xfd, 0xff]) add('uint8', c.uint8, n)
for (const n of [0, 1, 0xff, 0x100, 0xabcd, 0xffff]) add('uint16', c.uint16, n)
for (const n of [0, 1, 0xffff, 0x10000, 0xdeadbeef, 0xffffffff]) add('uint32', c.uint32, n)
for (const n of [0, 1, 0xffffffff, 0x100000000, 0x0001020304050607, MAX_SAFE_UINT]) add('uint64', c.uint64, n)

for (const n of [0, 1, -1, 63, -64, 127, -128]) add('int8', c.int8, n)
for (const n of [0, 1, -1, 127, -128, 128, -129, 32767, -32768]) add('int16', c.int16, n)
for (const n of [0, 1, -1, 32767, -32768, 32768, -32769, 2147483647, -2147483648]) add('int32', c.int32, n)
for (const n of [0, 1, -1, 2147483647, -2147483648, 2147483648, -2147483649, MAX_SAFE_INT, MIN_SAFE_INT]) add('int64', c.int64, n)

for (const s of ['', 'hello', '🌾', 'høsten er fin', 'x'.repeat(0xfc), 'y'.repeat(0xfd), 'z'.repeat(300), 'w'.repeat(0x10000)]) add('string', c.string, s)

for (const b of [[], [0], [1, 2, 3], Array.from({ length: 300 }, (_, i) => i & 0xff)]) {
  const buf = b4a.from(b)
  add('buffer', c.buffer, buf, b4a.toString(buf, 'hex'))
}

add('optionalBuffer', c.optionalBuffer, null)
for (const b of [[0], [1, 2, 3], Array.from({ length: 300 }, (_, i) => i & 0xff)]) {
  const buf = b4a.from(b)
  add('optionalBuffer', c.optionalBuffer, buf, b4a.toString(buf, 'hex'))
}

add('bool', c.bool, true)
add('bool', c.bool, false)

for (const a of [[], [1, 2, 3], [0, 0xfc, 0xfd, 0xffff, 0x10000, 0xffffffff, 0x100000000]]) add('array:uint', c.array(c.uint), a)
for (const a of [[], [0, -1, 1, MAX_SAFE_INT, MIN_SAFE_INT]]) add('array:int', c.array(c.int), a)
for (const a of [[], ['hello', 'world'], ['', '🌾']]) add('array:string', c.array(c.string), a)
add('array:bool', c.array(c.bool), [true, false, true])
add('array:array:uint', c.array(c.array(c.uint)), [[], [1], [2, 3]])

// A composite message, laid out the way the Go Marshal helper lays out a
// struct: fields in declaration order with no framing between them.
const message = {
  name: 'pear',
  age: 42,
  admin: true,
  data: b4a.from([0xde, 0xad, 0xbe, 0xef]),
  scores: [1, -2, 3],
  id: 0x0102030405060708n
}
const state = c.state()
c.string.preencode(state, message.name)
c.uint.preencode(state, message.age)
c.bool.preencode(state, message.admin)
c.buffer.preencode(state, message.data)
c.array(c.int).preencode(state, message.scores)
c.biguint64.preencode(state, message.id)
state.buffer = b4a.allocUnsafe(state.end)
c.string.encode(state, message.name)
c.uint.encode(state, message.age)
c.bool.encode(state, message.admin)
c.buffer.encode(state, message.data)
c.array(c.int).encode(state, message.scores)
c.biguint64.encode(state, message.id)
out.vectors.push({ codec: 'struct', value: null, hex: b4a.toString(state.buffer, 'hex') })

process.stdout.write(JSON.stringify(out, null, 1) + '\n')
