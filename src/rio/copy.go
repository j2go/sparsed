package rio

import (
	"bufio"
	"io"
	"os"
)

const ioBufferSize int = 4 * 1024
const ioReadBufferSize int = 100 * ioBufferSize

type createWriterFunc func(file *os.File) (writer io.Writer)
type createReaderFunc func(file *os.File) (writer io.Reader)

func doCopy(source string, dest string, createReader createReaderFunc, createWriter createWriterFunc) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	reader := createReader(sf)
	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	var writer io.Writer
	defer func() {
		if writer != nil {
			if flusher, ok := writer.(Flusher); ok {
				flusher.Flush()
			}
		}
		df.Close()
	}()
	writer = createWriter(df)
	buffer := make([]byte, ioBufferSize)
	_, err = io.CopyBuffer(writer, reader, buffer)
	if err != nil && err != io.EOF {
		return err
	}
	err = df.Sync()
	if err != nil {
		return err
	}
	return nil
}

// copyFile copies file source to destination dest.
func copyFile(source string, dest string) (err error) {
	passWriter := func(file1 *os.File) (writer io.Writer) {
		return file1
	}
	passReader := func(file2 *os.File) (writer io.Reader) {
		return file2
	}
	return doCopy(source, dest, passReader, passWriter)
}

// copyFileSizeSparse copies file source to destination dest.
func copyFileSizeSparse(source string, dest string, allocatedSize int64) (err error) {
	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = os.Truncate(dest, allocatedSize)
	if err != nil {
		return err
	}
	err = df.Close()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	createWriterTruncator := func(file *os.File) (writer io.Writer) {
		return NewSparseFilesWriter(NewSeekBufferedWriter(file, ioBufferSize))
	}
	createBufReader := func(file2 *os.File) (writer io.Reader) {
		return bufio.NewReaderSize(file2, ioReadBufferSize)
	}
	return doCopy(source, dest, createBufReader, createWriterTruncator)
}

// CopyFile copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	allocated, taken, err := FileSize(source)
	if err != nil {
		return err
	}
	// sparse file
	if allocated > taken {
		err = copyFileSizeSparse(source, dest, allocated)
	} else {
		err = copyFile(source, dest)
	}
	if err != nil {
		return err
	}

	si, err := os.Stat(source)
	if err != nil {
		return err
	}
	err = os.Chmod(dest, si.Mode())
	if err != nil {
		return err
	}
	uid, gid, err := Owner(source)
	if err != nil {
		return err
	}
	return os.Chown(dest, int(uid), int(gid))
}
