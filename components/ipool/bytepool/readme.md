## byte 池

降低频繁GC引起的开销

### usage

```
// github.com/dustin/go-humanize
// https://github.com/minio/minio/blob/30471212551b2507c44d0afe7d629a5ff460834b/pkg/bpool/bpool.go

maxCount := 32
blockSizeV1 := 10 * humanize.MiByte

bp := bpool.NewBytePoolCap(maxCount, blockSizeV1, blockSizeV1*2)
buf:=bp.Get()
defer bp.Put(buf)
//使用buf,不再举例

```

```
func mockReadFile(filepath string, b []byte) {
	f, _ := os.Open(filepath)
	for {
		n, err := io.ReadFull(f, b)
		if n == 0 || err == io.EOF {
			break
		}
	}
}

```

