package file

import "os"

// OpenCreateFile open an existed file or create a file if not exists
func OpenCreateFile(path string, flag int, perm os.FileMode) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		return f
	}
	// open file in read-write mode
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		panic(err)
	}
	return f
}
