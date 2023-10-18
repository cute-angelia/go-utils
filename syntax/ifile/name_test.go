package ifile

import (
	"log"
	"path"
	"testing"
)

/*
*
=== RUN   TestName

	name_test.go:8: a.jpg
	name_test.go:9: 1624955091917581000.jpg
	name_test.go:10: 361424000968063865.jpg
	name_test.go:13: test_a.jpg
	name_test.go:14: test_1624955091917614000.jpg
	name_test.go:15: test_361424000968129401.jpg

--- PASS: TestName (0.00s)
*/
func TestName(t *testing.T) {
	uri := "https://www.baidu.com/a.jpg?z=23"

	t.Log(NewFileName(uri).SetPrefix("good").SetSuffix("world").GetNameOrigin())
	t.Log(NewFileName(uri).GetNameTimeline())
	t.Log(NewFileName(uri).GetNameSnowFlow())

	prefix := "test_"
	t.Log(NewFileName(uri).SetPrefix(prefix).GetNameOrigin())
	t.Log(NewFileName(uri).SetPrefix(prefix).GetNameTimeline())
	t.Log(NewFileName(uri).SetPrefix(prefix).GetNameSnowFlow())

	// test 3
	iurl := "https://pbs.twimg.com/media/Eq3Ilp6XYAItcSw.jpg?format=jpg&name=orig"
	name := NewFileName(iurl).SetPrefix(prefix).GetNameTimeline()
	t.Log(name)

	ext := path.Ext(iurl)
	t.Log(ext)

	t.Log(NewFileName(iurl).GetNameTimelineReverse(true))
	t.Log(NewFileName(iurl).GetNameTimelineReverse(true))
	t.Log(NewFileName(iurl).GetNameTimelineReverse(true))
	t.Log(NewFileName(iurl).GetNameTimelineReverse(true))
	t.Log(NewFileName(iurl).GetNameTimelineReverse(true))
	t.Log(NewFileName(iurl).GetNameTimelineReverse(true))

	t.Log("GetNameSnowFlow")
	t.Log(NewFileName(iurl).GetNameSnowFlow())
	t.Log(NewFileName(iurl).GetNameSnowFlow())
	t.Log(NewFileName(iurl).GetNameSnowFlow())
	t.Log(NewFileName(iurl).GetNameSnowFlow())
	t.Log(NewFileName(iurl).GetNameSnowFlow())
	t.Log(NewFileName(iurl).GetNameSnowFlow())
}

func TestF(t *testing.T) {
	iurl := "https://pbs.twimg.com/media/Eq3Ilp6XYAItcSw.jpg?format=jpg&name=orig"
	log.Println(NewFileName(iurl).SetExt(".mp4").GetNameTimelineReverse(true))
	log.Println(NewFileName(iurl).GetExt())
}
