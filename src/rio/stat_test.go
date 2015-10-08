package rio

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

// TestCopyBigFile checks whether file copy works as expected on 17MB file
func TestFileSizeSparse(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	file, err := ioutil.TempFile(os.TempDir(), "revel_test_sparse")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	size := int64(1 * 1024 * 1024)
	file.Truncate(size)
	err = file.Sync()
	if err != nil {
		t.Fatal(err)
	}
	statSize, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if statSize.Size() != size {
		t.Fatalf("allocated size mismatch: %v != %v", statSize.Size(), size)
	}
	allocated, taken, err := FileSize(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if allocated != size {
		t.Fatalf("allocated size mismatch: %v != %v", allocated, size)
	}
	if taken != 0 {
		t.Fatalf("taken size mismatch: %v != %v", taken, 0)
	}
}
