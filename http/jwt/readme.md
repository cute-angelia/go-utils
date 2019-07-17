```go

/**
	中间件 类型检查
 */
 
func beforeCheckType(ctx iris.Context) {
	authorization := ctx.GetHeader("Authorization")
	resp := middlerResp{}

	token := strings.Split(authorization, "Bearer ")[1]
	ijwt := jwt.NewJwt(" xxxx  key  xxxx")

	if ijwt.Validate(token) != nil {
		resp.Code = -1
		resp.Message = "用户校验失败, 请重新登录"
		ctx.JSON(resp)
		return
	} else {
		itype := ctx.PostValue("type")
		if jwtinfo, err := ijwt.Decode(token); err != nil {
			resp.Code = -1
			resp.Message = err.Error() + " 请重新登录"
			ctx.JSON(resp)
			return
		} else {
			if jwtUid, err := jwtinfo.Get("uid"); err != nil {
				resp.Code = -1
				resp.Message = err.Error() + " 请重新登录"
				ctx.JSON(resp)
				return
			} else {
				// check right
				uid := fmt.Sprintf("%v", jwtUid)
				uidInt, _ := strconv.Atoi(uid)
				if ok, _, _ := pan.PanCheckTypeService(uidInt, itype); ok {
					ctx.Next()
					return
				} else {
					resp.Code = -1
					resp.Message = "无权限"
					ctx.JSON(resp)
					return
				}
			}
		}
	}

}


```