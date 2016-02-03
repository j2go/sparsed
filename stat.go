package sparsed

import "syscall"

// FileSize checks the real file size. Helpful for detecting sparse files.
func FileSize(path string) (allocated int64, taken int64, err error) {
	s := syscall.Stat_t{}
	err = syscall.Stat(path, &s)
	if err != nil {
		return 0, 0, err
	}
	blockSize, err := FSBlockSize(path)
	if err != nil {
		return 0, 0, err
	}
	return s.Size, s.Blocks * int64(blockSize), nil
}

// Owner returns uid and gid of a file
func Owner(path string) (uid uint32, gid uint32, err error) {
	s := syscall.Stat_t{}
	err = syscall.Stat(path, &s)
	if err != nil {
		return 0, 0, err
	}
	return s.Uid, s.Gid, nil
}

// FSBlockSize reports block size of particular file system
func FSBlockSize(path string) (blockSize int, err error) {
	s := syscall.Statfs_t{}
	err = syscall.Statfs(path, &s)
	if err != nil {
		return 0, err
	}
	return int(s.Bsize), nil
}
