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
		"https://scontent-lax3-2.cdninstagram.com/v/t51.2885-15/369770987_718587280037819_6034790256291537838_n.jpg?stp=dst-jpg_e35_p1080x1080&_nc_ht=scontent-lax3-2.cdninstagram.com&_nc_cat=111&_nc_ohc=1w35qNRTa2sAX-wkt6d&edm=ACWDqb8BAAAA&ccb=7-5&ig_cache_key=MzE3NjM3NzQ2NzIwMzMyMDE4OQ%3D%3D.2-ccb7-5&oh=00_AfBFrXIrSe2ORWISyJnAliyXOrBKFMuE_7Nw6YJIb8mCRg&oe=64F5D3B6&_nc_sid=ee9879",
	}

	for _, path := range testPaths {
		//log.Println(path)
		//log.Println(Clean(path))
		//log.Println(filepath.Ext(path))
		log.Println(GetFileName(path))
		// log.Println(GetFileExt(path))
	}
}
