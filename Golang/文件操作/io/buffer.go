package io

import (
	"os"
)

type BufferFileWriter struct {
	fout 	*os.File
	buffer 	[]byte
	endOfIndex int
}

func NewBufferFileWriter(fout *os.File,buffersize int) *BufferFileWriter {
	return &BufferFileWriter{
		fout: fout,
		buffer: make([]byte, buffersize),
		endOfIndex: 0,
	}
}

func (w *BufferFileWriter) Flush() {
	w.fout.Write(w.buffer[0:w.endOfIndex])
	w.endOfIndex = 0
}


func (w *BufferFileWriter) Write(cont []byte) {
	if len(cont) > len(w.buffer) {
		w.Flush()
		w.fout.Write(cont)
	}else {
		if len(cont) + w.endOfIndex > len(w.buffer) {
			w.Flush()
		}
		copy(w.buffer[w.endOfIndex:], cont)
		w.endOfIndex += len(cont)
	}
}

func (w *BufferFileWriter) WriteString(s string) {
	w.Write([]byte(s))
}

