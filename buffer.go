package bytebuffer

import (
	"reflect"
	"sync"
	"unsafe"
)

const (
	// BaseLength is default length for byte array
	BaseLength = 0
	// BaseCapacity is default capacity for byte array
	BaseCapacity = 1024
)

// Pool is struct for buffer pool
type Pool struct {
	p *sync.Pool
}

// Buffer is main type for bytebuffer
type Buffer struct {
	BS   []byte
	len  int
	pool Pool
}

// Config ...
type Config struct {
	Cap int
}

// NewPool get new pool for buffer
func NewPool(c *Config) Pool {
	var conf *Config
	if c == nil {
		conf = &Config{
			Cap: BaseCapacity,
		}
	} else {
		cap := BaseCapacity

		if c.Cap > 1 {
			cap = c.Cap
		}

		conf = &Config{
			Cap: cap,
		}
	}

	return Pool{p: &sync.Pool{
		New: func() interface{} {
			return &Buffer{BS: make([]byte, BaseLength, conf.Cap)}
		},
	}}
}

// Get new buffer from pool
func (p Pool) Get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

// put buffer back to pool
func (p Pool) put(buf *Buffer) {
	p.p.Put(buf)
}

// Reset bytes buffer
func (b *Buffer) Reset() {
	b.BS = b.BS[:0]
	b.len = 0
}

// Grow will increase length of byte array
func (b *Buffer) Grow(n int) {
	b.len += n
	b.BS = b.BS[:b.len]
}

// Write bytes to buffer
func (b *Buffer) Write(bs []byte) (int, error) {
	m := b.len
	n := m + len(bs)
	b.Grow(len(bs))
	return copy(b.BS[m:n], bs), nil
}

// WriteByte will write a Byte
func (b *Buffer) WriteByte(v byte) error {
	b.BS = append(b.BS, v)
	return nil
}

// WriteString will string string as []byte
func (b *Buffer) WriteString(s string) error {
	b.Write([]byte(s))
	return nil
}

// Len will return length of byte array
func (b *Buffer) Len() int {
	return b.len
}

// Cap will return capacity of byte array
func (b *Buffer) Cap() int {
	return cap(b.BS)
}

// Bytes return current bytes
func (b *Buffer) Bytes() []byte {
	return b.BS
}

func (b *Buffer) String() string {
	return byteToString(b.BS)
}

// Free will Release Buffer to pool
func (b *Buffer) Free() {
	b.pool.put(b)
}

// byteToString unsafe cast from array of byte to string
func byteToString(bytes []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}
