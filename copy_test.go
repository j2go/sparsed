package rio

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func copyFileTestBody(t *testing.T, testName string, size int) {
	file, err := ioutil.TempFile(os.TempDir(), "revel_test_copy")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	content := make([]byte, size)
	_, err = rand.Read(content)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	dstName := fmt.Sprintf("%v_dst", file.Name())
	defer os.Remove(dstName)
	err = copyFile(file.Name(), dstName)
	if err != nil {
		t.Fatal(err)
	}

	dstFile, err := os.Open(dstName)
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()
	content2, err := ioutil.ReadAll(dstFile)
	if err != nil {
		t.Fatal(err)
	}
	if md5.Sum(content) != md5.Sum(content2) {
		t.Fatalf("%v: content mismatch", testName)
	}
}

func copySparseFileTestBody(t *testing.T, testName string, allocatedSize int64) {
	if allocatedSize < 512 {
		t.Fatalf("%v is too small for sparse file", allocatedSize)
	}
	file, err := ioutil.TempFile(os.TempDir(), "revel_test_sparse")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	err = file.Truncate(allocatedSize)
	if err != nil {
		t.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
	srcFile, err := os.OpenFile(file.Name(), os.O_RDONLY, 0444)
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()
	srcContent, err := ioutil.ReadAll(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	dstName := fmt.Sprintf("%v_dst", file.Name())
	defer os.Remove(dstName)
	err = copyFileSizeSparse(file.Name(), dstName, allocatedSize)
	if err != nil {
		t.Fatal(err)
	}

	dstFile, err := os.Open(dstName)
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()
	dstContent, err := ioutil.ReadAll(dstFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(srcContent) != len(dstContent) {
		t.Fatalf("%v: length mismatch %v != %v", testName, len(srcContent), len(dstContent))
	}
	if md5.Sum(srcContent) != md5.Sum(dstContent) {
		t.Fatalf("%v: content mismatch", testName)
	}
}

// TestCopyBiggerFile checks whether file copy works as expected on 3MB file
func TestCopyBiggerFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestCopyBiggerFile in short mode.")
	}
	copyFileTestBody(t, "TestCopyBiggerFile", 3*1024*1024+24)
}

// TestCopyMediumFile checks whether file copy works as expected on 1MB file
func TestCopyMediumFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestCopyMediumFile in short mode.")
	}
	copyFileTestBody(t, "TestCopyMediumFile", 1*1024*1024-11)
}

// TestCopySmallFile checks whether file copy works as expected on 256KB file
func TestCopySmallFile(t *testing.T) {
	copyFileTestBody(t, "TestCopySmallFile", 256*1024-8)
}

// TestCopyRealSmallFile checks whether file copy works as expected on 3B file
func TestCopyRealSmallFile(t *testing.T) {
	copyFileTestBody(t, "TestCopyRealSmallFile", 3)
}

// TestCopySparseSmall checks small sparse files
func TestCopySparseSmall(t *testing.T) {
	copySparseFileTestBody(t, "TestCopySparseSmall", 512+48)
}

// TestCopySparseMedium checks small sparse files
func TestCopySparseMedium(t *testing.T) {
	copySparseFileTestBody(t, "TestCopySparseMedium", 8*1024-23)
}
