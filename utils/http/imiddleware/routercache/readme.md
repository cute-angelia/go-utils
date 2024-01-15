## 路由缓存中间件

### usage

常规用法：

```go

package demo

// 定义一个列表缓存
func getCacheList() *routercache.Component {
	return routercache.New(
		routercache.WithStore(ibunt.GetComponent("cache")),
		routercache.WithTtl(time.Minute*20),
		routercache.WithCustomKey("xxxxx"), // 特殊需求 - 如：后台清理特定缓存
	)
}

// 使用缓存
func (rs Posts) Routes() chi.Router {
    r := chi.NewRouter()
	// 使用中间件
    r.With(getCacheList().NewMiddleware).Post("/listByPageOneContent", rs.listByPageOneContent)
}


// 清理缓存
head -> 带刷新key

// 清理特定缓存
routercache.New(routercache.WithStore(ibunt.GetComponent("cache"))).DeleteCustomKey("xxxxx")

```