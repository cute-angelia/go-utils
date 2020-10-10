### 使用说明

```

```


### gout 的 cookie 管理

```
jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
if err != nil {
        log.Fatal(err)
}

client := &http.Client{
        Jar: jar,
}

gout.New(client).GET(ts.URL).SetBody("hi cookie").Do()
gout.New(client).GET(ts.URL).SetBody("hi cookie2").Do() //抓下包就可以看到自动持有服务端一开始设置的cookie了。


type rspHeader struct {
    SetCookie string `header:"Set-Cookie"`
}
var header rspHeader
gout.New().GET(ts.URL).SetBody("hi cookie").BindHeader(&header).Do()
fmt.Printf("cookie value:%v\n", header)
```