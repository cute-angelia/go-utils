package apicache

import (
	"github.com/cute-angelia/go-utils/utils/encrypt/hash"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type Option func(c *Container)

type Container struct {
	config *config
	name   string
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	// 两种方式，一种是 ego 的 config 加载，一种是option with 加载
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.Timeout = timeout
	}
}
func WithOnlyToday(onlyTodya bool) Option {
	return func(c *Container) {
		c.config.OnlyToday = onlyTodya
	}
}
func WithCacheKey(cacheKey string) Option {
	return func(c *Container) {
		c.config.CacheKey = cacheKey
	}
}

func WithGenerateCacheKey(params interface{}, filtes []string) Option {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	bparms, _ := json.Marshal(params)
	m := make(map[string]interface{})
	json.Unmarshal(bparms, &m)

	for _, filte := range filtes {
		if _, ok := m[filte]; ok {
			delete(m, filte)
		}
	}

	// log.Printf("%#v", m)

	cacheKey, _ := json.Marshal(m)
	cacheKeyMd5 := hash.NewEncodeMD5(string(cacheKey))

	return func(c *Container) {
		c.config.CacheKey = cacheKeyMd5
	}
}

func WithDbName(dbName string) Option {
	return func(c *Container) {
		c.config.DbName = dbName
	}
}
func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}
func WithPrefix(prefix string) Option {
	return func(c *Container) {
		c.config.Prefix = prefix
	}
}

func WithPrefixMaxNum(prefixMaxNum int) Option {
	return func(c *Container) {
		c.config.PrefixMaxNum = prefixMaxNum
	}
}

// Build ...
func (c *Container) MustBuild(dbName string, prefix string, options ...Option) *Component {

	// Must
	c.config.DbName = dbName
	c.config.Prefix = prefix

	for _, option := range options {
		option(c)
	}
	// log.Println(PackageName, fmt.Sprintf("%+v", c.config))

	if len(c.config.Prefix) == 0 {
		c.logger.Error("cachekey prefix is empty")
		panic("cachekey prefix is empty")
	}
	if len(c.config.CacheKey) == 0 {
		c.logger.Error("cachekey is empty")
		// panic("cachekey is empty")
	}
	if len(c.config.DbName) == 0 {
		c.logger.Error("DbName is empty")
		panic("DbName is empty")
	}

	if c.config.Debug {
		c.logger.SetLevel(elog.DebugLevel)
	}

	return newComponent(c.name, c.config, c.logger)
}
