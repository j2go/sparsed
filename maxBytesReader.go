package rio

import "io"

// MaxBytesReader reader which is capable of cpoying of no more than `maxBytes`
type MaxBytesReader struct {
	reader    io.Reader
	readBytes int64
	maxBytes  int64
}

// Count return current number of bytes read
func (r *MaxBytesReader) Count() int64 {
	return r.readBytes
}

// Max returns reader.maxBytes value
func (r *MaxBytesReader) Max() int64 {
	return r.maxBytes
}

//NewMaxBytesReader creates MaxBytesReader which is capable of cpoying of no more than `maxBytes`
func NewMaxBytesReader(innerReader io.Reader, maxBytes int64) *MaxBytesReader {
	return &MaxBytesReader{reader: innerReader, maxBytes: maxBytes, readBytes: 0}
}

// Read io.reader.Read implementation
func (r *MaxBytesReader) Read(p []byte) (n int, err error) {
	if r.readBytes >= r.maxBytes {
		return 0, io.EOF
	}
	if int64(len(p)) > (r.maxBytes - r.readBytes) {
		buf := make([]byte, r.maxBytes-r.readBytes)
		v, err := r.reader.Read(buf)
		r.readBytes += int64(v)
		for i := range buf {
			p[i] = buf[i]
		}
		return v, err
	}
	v, err := r.reader.Read(p)
	r.readBytes += int64(v)
	return v, err
}
