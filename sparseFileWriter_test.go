package sparsed

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var randomData []byte

func init() {
	randomData = make([]byte, 65*1024*1024)
	rand.Read(randomData)
}

// TestDetectOneZeroBlock checks if one zero block is correctly detected by underlying mechanism of SparseFileWriter.
func TestDetectOneZeroBlock(t *testing.T) {
	data := make([]byte, 1024+22)
	blocks := detectBlocks(data)
	if len(blocks) != 1 {
		t.Fatalf("TestDetectZeroBlock: expected 1 zero block, got %v", len(blocks))
	}
}

// TestDetectString checks if a simple string is correctly detected by underlying mechanism of SparseFileWriter.
func TestDetectString(t *testing.T) {
	data := []byte("Hello brave new world")

	detectTestBody(t, "TestDetectString", data)
}

// TestDetectWrappedString checks if a zero wrapped string is correctly detected by underlying mechanism of SparseFileWriter.
func TestDetectWrappedString(t *testing.T) {
	var data []byte
	data = append(data, make([]byte, 128-3)...)
	data = append(data, []byte("Hello brave new")...)
	data = append(data, make([]byte, 512+4)...)
	data = append(data, []byte("world")...)
	data = append(data, make([]byte, 64-11)...)

	detectTestBody(t, "TestDetectWrappedString", data)
}

// TestDetectWrappingString checks if a wrapping string is correctly detected by underlying mechanism of SparseFileWriter.
func TestDetectWrappingString(t *testing.T) {
	var data []byte
	data = append(data, []byte("Hello")...)
	data = append(data, make([]byte, 128-3)...)
	data = append(data, []byte("brave")...)
	data = append(data, make([]byte, 512+4)...)
	data = append(data, []byte("new")...)
	data = append(data, make([]byte, 64-11)...)
	data = append(data, []byte("world")...)

	detectTestBody(t, "TestDetectWrappingString", data)
}

// TestDetectRandomBytes checks if various random bytes are correctly detected by underlying mechanism of SparseFileWriter.
func TestDetectRandomBytes(t *testing.T) {
	testBody := func(bytesCount int) {
		data := randomData[0:bytesCount]

		detectTestBody(t, fmt.Sprintf("TestDetectRandomBytes(%v)", bytesCount), data)
	}
	testSamples := []int{256, 512, 1024, 2048, 4096, 256 + 8, 512 - 13, 1024 + 2, 2048 - 1, 4096 + 31}
	for _, v := range testSamples {
		testBody(v)
	}
}

func detectBlocksAndConvert(data []byte) []byte {
	result := make([]byte, len(data))
	for _, block := range detectBlocks(data) {
		switch block.empty {
		case true:
			for i := 0; i < block.count; i++ {
				result[i+block.first] = 0x00
			}
		case false:
			for i := 0; i < block.count; i++ {
				result[i+block.first] = data[i+block.first]
			}
		}
	}
	return result
}

// detectTestBody test helper method
func detectTestBody(t testing.TB, testName string, data []byte) {
	result := detectBlocksAndConvert(data)
	if string(data) != string(result) {
		t.Fatalf("%v: %v does not match to %v", testName, string(data), string(result))
	}
}

// BenchmarkDetectBigRandomBytes benchmark random bytes on big samples.
func BenchmarkDetectBigRandomBytes(t *testing.B) {
	testBody := func(bytesCount int) {
		data := randomData[0:bytesCount]

		_ = detectBlocksAndConvert(data)
	}
	testSamples := []int{64 * 1024 * 1024, 32 * 1024 * 1024, 8*1024*1024 - 5}
	for _, v := range testSamples {
		testBody(v)
	}
}

// BenchmarkDetectSmallRandomBytes benchmark random bytes on small samples.
func BenchmarkDetectSmallRandomBytes(t *testing.B) {
	testBody := func(bytesCount int) {
		data := randomData[0:bytesCount]
		_ = detectBlocksAndConvert(data)
	}
	testSamples := []int{1024*1024 - 13, 512 - 11, 128 + 1}
	for _, v := range testSamples {
		testBody(v)
	}
}

// BenchmarkSparseFileWriter benchmark file write.
func BenchmarkSparseFileWriter(t *testing.B) {
	testBody := func(bytesCount int64) {
		file, err := ioutil.TempFile(os.TempDir(), "revel_benchmark_sparse")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()

		err = file.Truncate(bytesCount)
		if err != nil {
			t.Fatal(err)
		}

		dstName := fmt.Sprintf("%v_dst", file.Name())
		defer os.Remove(dstName)
		err = copyFileSizeSparse(file.Name(), dstName, bytesCount)
		if err != nil {
			t.Fatal(err)
		}
	}
	testSamples := []int64{32 * 1024 * 1024, 12 * 1024 * 1024, 4 * 1024 * 1024}
	for _, v := range testSamples {
		testBody(v)
	}
}
