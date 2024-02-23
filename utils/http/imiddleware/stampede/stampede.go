package stampede

/* 多次重复提交 防穿透 https://en.wikipedia.org/wiki/Cache_stampede#cite_note-1
多次请求返回同一结果给客户端
*/

/*
	customKeyFunc := func(r *http.Request) uint64 {
		// Read the request payload, and then setup buffer for future reader
		var buf []byte
		if r.Body != nil {
			buf, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(buf))
		}

		token := r.Header.Get("Authorization")

		return stampede.BytesToHash(
			[]byte(strings.ToLower(r.URL.Path)),
			[]byte(strings.ToLower(token)),
			buf,
		)
	}

	cached := stampede.HandlerWithKey(512, 10*time.Second, customKeyFunc)

	// 新版
	r.With(cached).Route("/v2", func(r chi.Router) {
		r.Post("/join", that.joinV2)
		r.Post("/out", that.out)
	})
*/
