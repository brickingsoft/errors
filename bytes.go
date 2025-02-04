package errors

import (
	"bytes"
	"sync"
)

var bytebufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func acquireByteBuffer() *bytes.Buffer {
	return bytebufferPool.Get().(*bytes.Buffer)
}

func releaseByteBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bytebufferPool.Put(buf)
}
