# go-utils

[![Page Views Count](https://badges.toozhao.com/badges/01EH4J7MXDPXTXC3MCMMZ407PV/green.svg)](https://badges.toozhao.com/badges/01EH4J7MXDPXTXC3MCMMZ407PV/green.svg "Get your own page views count badge on badges.toozhao.com")

go开发工具包，因很多项目都要重复写一些包，如日志，Gorm池，缓存组件

> ps: v1 版本，支持 go version 1.16 及以下

```shell
go get github.com/cute-angelia/go-utils@v1
```


包说明：

1. components (功能模块，是工具类的升级版)
2. syntax (语法类)
3. utils (工具类)
4. examples (例子)

### 主要模块

| 模块   |                                                              |
| ------ | ------------------------------------------------------------ |
| Db     | gorm                                                         |
| 缓存   | redis, buntdb(file), LRU cache                               |
| HTTP   | API,JWT,vaildation                                           |
| logger | file-rotatelogs, stdout                                      |
| syntax | file,slice,string,time,zip                                   |
| utils  | encrypt, orderid,snowflake, idownload, ip, store,task,etc... |
| limit  | retry, risk                                                  |

### Refer:

1. [goutil](https://github.com/gookit/goutil)
2. [gotools](https://github.com/asktop/gotools)
3. [golang-examples](https://github.com/SimonWaldherr/golang-examples)
4. [kit](https://github.com/ardanlabs/kit)
5. [time](https://github.com/jinzhu/now)
6. [go-extend](https://github.com/cute-angelia/go-extend)
7. [yiigo](https://github.com/shenghui0779/yiigo)
8. [downloader](https://github.com/polaris1119/downloader)
9. [写了 30 多个 Go 常用文件操作的示例](https://mp.weixin.qq.com/s/dczWeHW6JWSJMJx1nBx7rA)
10. [lancet](https://github.com/duke-git/lancet)
