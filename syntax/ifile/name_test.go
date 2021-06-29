package ifile

import "testing"

/**
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

	t.Log(NewFileName(uri).GetNameOrigin(""))
	t.Log(NewFileName(uri).GetNameTimeline(""))
	t.Log(NewFileName(uri).GetNameSnowFlow(""))

	prefix := "test_"
	t.Log(NewFileName(uri).GetNameOrigin(prefix))
	t.Log(NewFileName(uri).GetNameTimeline(prefix))
	t.Log(NewFileName(uri).GetNameSnowFlow(prefix))
}
