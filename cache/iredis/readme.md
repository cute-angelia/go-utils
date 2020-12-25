
```
	c := iredis2.NewRedisPool()
	c.Init("main", "192.168.2.140:6379", "")
	defer c.Close()

	var ctx = context.Background()

	keyzset := "sc"
	rdb := iredis2.GetRdb("main")

	// zadd
	member := redis.Z{
		Member: 24,
		Score:  float64(time.Now().Unix()),
	}
	rdb.ZAdd(ctx, keyzset, &member)
```