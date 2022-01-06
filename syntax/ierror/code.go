package ierror

//go:generate stringer -type Code -linecomment
const (
	// 服务级错误码
	ServerError        Code = 10101 // Internal Server Error
	TooManyRequests    Code = 10102 // Too Many Requests
	ParamBindError     Code = 10103 // 参数信息有误
	AuthorizationError Code = 10104 // 签名信息有误
	CallHTTPError      Code = 10105 // 调用第三方HTTP接口失败
	ResubmitError      Code = 10106 // ResubmitError
	ResubmitMsg        Code = 10107 // 请勿重复提交
	HashIdsDecodeError Code = 10108 // ID参数有误
	SignatureError     Code = 10109 // SignatureError

	// 业务模块级错误码
	// 用户模块
	IllegalUserName Code = 20101 // 非法用户名
	UserCreateError Code = 20102 // 创建用户失败
	UserUpdateError Code = 20103 // 更新用户失败
	UserSearchError Code = 20104 // 查询用户失败

	// 配置
	ConfigEmailError        Code = 20401 // 修改邮箱配置失败
	ConfigSaveError         Code = 20402 // 写入配置文件失败
	ConfigRedisConnectError Code = 20403 // Redis连接失败
	ConfigMySQLConnectError Code = 20404 // MySQL连接失败
	ConfigMySQLInstallError Code = 20405 // MySQL初始化数据失败
	ConfigGoVersionError    Code = 20406 // GoVersion不满足要求
)
