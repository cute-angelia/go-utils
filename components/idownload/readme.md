### 下载组件

```go
idownload.Load("").Build(
    idownload.WithProxySocks5(e.config.ProxySocks5),
    idownload.WithDebug(e.config.Debug),
    idownload.WithTimeout(e.config.Timeout),
    idownload.WithReferer(e.config.Referer),
)
```

搭配协程池

https://github.com/wazsmwazsm/mortar