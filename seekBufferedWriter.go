package rio

import (
	"bufio"
	"io"
	"os"
)

var emptyBuffer []byte

func init() {
	emptyBuffer = make([]byte, ioBufferSize)
}

// SeekBufferedWriter bufio.Writer with seek method
type SeekBufferedWriter struct {
	writer         io.WriteSeeker
	bufferedWriter *bufio.Writer
	bufferSize     int
	emptyBuffer    []byte
}

// NewSeekBufferedWriter creates SeekBufferedWriter
func NewSeekBufferedWriter(writer io.WriteSeeker, bufferSize int) *SeekBufferedWriter {
	bufWriter := bufio.NewWriterSize(writer, bufferSize)
	return &SeekBufferedWriter{writer: writer, bufferSize: bufferSize, bufferedWriter: bufWriter, emptyBuffer: emptyBuffer}
}

// Seek flush, seek and replace bufferedWriter
func (w *SeekBufferedWriter) Seek(offset int64, whence int) (int64, error) {
	onSeek := func() (int64, error) {
		err := w.bufferedWriter.Flush()
		if err != nil {
			return 0, err
		}
		result, err := w.writer.Seek(offset, whence)
		if err != nil {
			return result, err
		}
		w.bufferedWriter = bufio.NewWriterSize(w.writer, w.bufferSize)
		return result, err
	}

	switch whence {
	case os.SEEK_CUR:
		if offset < int64(w.bufferedWriter.Available()) {
			w.bufferedWriter.Write(emptyBuffer[0:offset])
			return offset, nil
		}
		return onSeek()
	default:
		return onSeek()
	}
}

// Flush pushes the changes to underlying writer.
func (w *SeekBufferedWriter) Flush() error {
	return w.bufferedWriter.Flush()
}

// Write write p bytes to bufferedWriter.
func (w *SeekBufferedWriter) Write(p []byte) (int, error) {
	return w.bufferedWriter.Write(p)
}
