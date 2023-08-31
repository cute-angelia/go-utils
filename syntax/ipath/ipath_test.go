package ipath

import (
	"log"
	"testing"
)

func TestPath(t *testing.T) {
	testPaths := []string{
		"/Users/vanilla/Downloads/tuli/x.jpg",
		"/x.jpg",
		"x.jpg",
		"/ab/bc/de/x",
		"/ab/bc/de/x.jpg",
		"/Users/vanilla/Downloads/tuli/x.jpg?quiery=dfsa&dfas=2323",
	}

	for _, path := range testPaths {
		//log.Println(path)
		//log.Println(Clean(path))
		//log.Println(filepath.Ext(path))
		log.Println(GetFileName(path))
		// log.Println(GetFileExt(path))
	}
}
