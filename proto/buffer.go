package proto

import (
	"sync"
)

var bufferp = sync.Pool{
	New: func() interface{} {
		return new(Buffer)
	},
}

// AllocBuffer gets buffer from pool
func AllocBuffer() *Buffer {
	buf := bufferp.Get().(*Buffer)
	buf.Reset()
	return buf
}

// FreeBuffer puts buffer to pool if cap of buffer less than 64k
func FreeBuffer(b *Buffer) {
	if b.Cap() < 1<<16 {
		bufferp.Put(b)
	}
}

type Buffer struct {
	buf []byte
	off int
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func (b *Buffer) Len() int {
	return len(b.buf) - b.off
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.off:]
}

func (b *Buffer) Reset() {
	b.off = 0
	b.buf = b.buf[:0]
}

func (b *Buffer) Reserve(n int) {
	if cap(b.buf) < n {
		buf := make([]byte, len(b.buf)-b.off, n)
		copy(buf, b.buf[b.off:])
		b.buf = buf
		b.off = 0
	}
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *Buffer) Unmarshal(m Message) error {
	return Unmarshal(b.Bytes(), m)
}

func (b *Buffer) Encode(m Message) error {
	b.Reset()
	var err error
	b.buf, err = EncodeAppend(b.buf, m)
	return err
}
