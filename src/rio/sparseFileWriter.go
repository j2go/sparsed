package rio

import (
	"io"
	"os"
)

const dataBlockDefaultCapacity int = 4

// Flusher io.Writer with Flush method. Helpful for bufferedWriter underneath.
type Flusher interface {
	Flush() error
}

// SparseFileWriter smarter writer that can detect sparse blocks.
type SparseFileWriter struct {
	writer io.WriteSeeker
}

// NewSparseFilesWriter constructor for SparseFileWriter. Wraps another writer (file most likely).
func NewSparseFilesWriter(innerWriter io.WriteSeeker) *SparseFileWriter {
	return &SparseFileWriter{writer: innerWriter}
}

// Flush calls Flush() on underlying writer.
func (w *SparseFileWriter) Flush() error {
	if flusher, ok := w.writer.(Flusher); ok {
		return flusher.Flush()
	}
	return nil
}

// Write copies bytes to underlying writer. Detects sparse blocks and calls truncate when needed.
func (w *SparseFileWriter) Write(p []byte) (n int, err error) {
	n = 0
	for _, block := range detectBlocks(p) {
		switch block.empty {
		case true:
			_, err = w.writer.Seek(int64(block.count), os.SEEK_CUR)
			if err != nil {
				return n, err
			}
			n += block.count
		case false:
			wn, err := w.writer.Write(p[block.first : block.first+block.count])
			n += wn
			if err != nil {
				return n, err
			}
		}
	}
	return len(p), err
}

type dataBlock struct {
	first int
	count int
	empty bool
}

// detectBlocks splits incoming bytes into zero or data blocks
func detectBlocks(data []byte) []*dataBlock {
	blocks := make([]*dataBlock, dataBlockDefaultCapacity, dataBlockDefaultCapacity)
	blocksCount := 0
	fastAppend := func(item *dataBlock) {
		if blocksCount >= len(blocks) {
			blocks = append(blocks, make([]*dataBlock, dataBlockDefaultCapacity, dataBlockDefaultCapacity)...)
		}
		blocks[blocksCount] = item
		blocksCount++
	}

	sparseBegin := -1
	dataBegin := -1

	for i, v := range data {
		if v == 0x00 {
			if sparseBegin == -1 {
				sparseBegin = i
			}
			if dataBegin != -1 {
				fastAppend(&dataBlock{first: dataBegin, empty: false, count: i - dataBegin})
				dataBegin = -1
			}
		} else {
			if dataBegin == -1 {
				dataBegin = i
			}
			if sparseBegin != -1 {
				fastAppend(&dataBlock{first: sparseBegin, empty: true, count: i - sparseBegin})
				sparseBegin = -1
			}
		}
	}
	if sparseBegin != -1 {
		fastAppend(&dataBlock{first: sparseBegin, empty: true, count: len(data) - sparseBegin})
	}
	if dataBegin != -1 {
		fastAppend(&dataBlock{first: dataBegin, empty: false, count: len(data) - dataBegin})
	}
	return blocks[:blocksCount]
}
