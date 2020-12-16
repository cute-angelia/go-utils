package umysql

import "time"

type GormOptions struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Dbname       string
	Dsn          string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
	LogDebug     bool
}

type GormOption func(options *GormOptions)

func NewGormOpts(opts ...GormOption) GormOptions {
	var sopt GormOptions
	for _, opt := range opts {
		opt(&sopt)
	}

	if sopt.Port == "" {
		sopt.Port = "3306"
	}

	if sopt.MaxIdleConns == 0 {
		sopt.MaxIdleConns = int(10)
	}

	sopt.LogDebug = true

	return sopt
}

func WithGormOptHost(host string) GormOption {
	return func(options *GormOptions) {
		options.Host = host
	}
}

func WithGormOptPort(port string) GormOption {
	return func(options *GormOptions) {
		options.Port = port
	}
}

func WithGormOptUserName(username string) GormOption {
	return func(options *GormOptions) {
		options.Username = username
	}
}

func WithGormOptPassword(pwd string) GormOption {
	return func(options *GormOptions) {
		options.Password = pwd
	}
}

func WithGormOptDbname(dbname string) GormOption {
	return func(options *GormOptions) {
		options.Dbname = dbname
	}
}

func WithGormOptMaxIdleConns(connmax int) GormOption {
	return func(options *GormOptions) {
		options.MaxIdleConns = connmax
	}
}

func WithGormOptMaxOpenConns(connmax int) GormOption {
	return func(options *GormOptions) {
		options.MaxOpenConns = connmax
	}
}

func WithGormLogDebug(logdebug bool) GormOption {
	return func(options *GormOptions) {
		options.LogDebug = logdebug
	}
}

func WithGormDsn(dsn string) GormOption {
	return func(options *GormOptions) {
		options.Dsn = dsn
	}
}

func WithMaxLifeTime(maxLifeTime time.Duration) GormOption {
	return func(options *GormOptions) {
		options.MaxLifetime = maxLifeTime
	}
}
