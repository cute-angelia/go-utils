package redis

import (
	"context"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"strconv"
	"strings"
)

var readonlyCommand map[string]bool

type ConnGetter func(string) (redigo.Conn, error)

func init() {
	readonlyCommand = make(map[string]bool)
	readCmdlist := []string{"info", "smembers", "hlen", "hmget", "srandmember", "hvals", "randomkey", "strlen", "dbsize", "keys", "ttl", "lindex", "type", "llen", "dump", "scard", "echo", "lrange", "zcount", "exists", "sdiff", "zrange", "mget", "zrank", "get", "getbit", "getrange", "zrevrange", "zrevrangebyscore", "hexists", "object", "sinter", "zrevrank", "hget", "zscore", "hgetall", "sismember"}
	for _, cmd := range readCmdlist {
		readonlyCommand[cmd] = true
		readonlyCommand[strings.ToUpper(cmd)] = true
	}
}

type conn struct {
	redigo.Conn
	ctx context.Context
}

func (t *conn) Send(commandName string, args ...interface{}) (err error) {
	var sp opentracing.Span
	if t.ctx != nil {
		sp, _ = opentracing.StartSpanFromContext(t.ctx, "redis.send")
	} else {
		sp = opentracing.StartSpan("redis.send")
	}
	defer sp.Finish()

	ext.DBType.Set(sp, "redis")
	ext.DBStatement.Set(sp, commandString(commandName, args))

	err = t.Conn.Send(commandName, args...)
	if err != nil {
		ext.Error.Set(sp, true)
		sp.LogFields(log.Error(err))
	}
	return
}

func (t *conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	var sp opentracing.Span
	if t.ctx != nil {
		sp, _ = opentracing.StartSpanFromContext(t.ctx, "redis.do")
	} else {
		sp = opentracing.StartSpan("redis.do")
	}

	ext.DBType.Set(sp, "redis")
	ext.DBStatement.Set(sp, commandString(commandName, args))

	reply, err = t.Conn.Do(commandName, args...)
	if err != nil {
		ext.Error.Set(sp, true)
		sp.LogFields(log.Error(err))
	}
	sp.Finish()

	return
}

func commandString(command string, args []interface{}) string {
	var buffer strings.Builder
	buffer.WriteString(command)
	buffer.WriteByte(' ')
	for idx, arg := range args {
		if idx != 0 {
			buffer.WriteByte(' ')
		}
		var str string
		switch arg := arg.(type) {
		case string:
			str = arg
		case []byte:
			str = string(arg)
		case int:
			str = strconv.Itoa(arg)
		case int64:
			str = strconv.Itoa(int(arg))
		case float64:
			str = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				str = "true"
			} else {
				str = "false"
			}
		case nil:
			str = "nil"
		}
		buffer.WriteString(str)
	}

	return buffer.String()
}

type errorConnection struct {
	err error
	ctx context.Context
}

func (ec errorConnection) Err() error                    { return ec.err }
func (ec errorConnection) Close() error                  { return ec.err }
func (ec errorConnection) Flush() error                  { return ec.err }
func (ec errorConnection) Receive() (interface{}, error) { return nil, ec.err }

func (ec *errorConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if ec.ctx != nil {
		sp, _ := opentracing.StartSpanFromContext(ec.ctx, "redis.do")
		ext.DBType.Set(sp, "redis")
		ext.DBStatement.Set(sp, commandString(commandName, args))
		ext.Error.Set(sp, true)
		sp.LogFields(log.Error(ec.err))
		sp.Finish()
	}
	return nil, ec.err
}

func (ec *errorConnection) Send(commandName string, args ...interface{}) error {
	if ec.ctx != nil {
		sp, _ := opentracing.StartSpanFromContext(ec.ctx, "redis.send")
		ext.DBType.Set(sp, "redis")
		ext.DBStatement.Set(sp, commandString(commandName, args))
		ext.Error.Set(sp, true)
		sp.LogFields(log.Error(ec.err))
		sp.Finish()
	}
	return ec.err
}
