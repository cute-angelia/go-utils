package redis

import (
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"context"
	"log"
)

var RedisCmder *RedisCmd

// 初始化
func InitRedis(name string, ip string, port string, auth string) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	RedisCmder = &RedisCmd{
		Name: name,
		Pool: New(name, MakeConfig(addr, auth)),
	}
}

// --- CMD ---
type RedisCmd struct {
	Name string
	Pool Pool
	ctx  context.Context
}

func (self *RedisCmd) SetCtx(ctx context.Context) {
	self.ctx = ctx
}

// 获得连接
func (self *RedisCmd) GetConn() redigo.Conn {
	var conn redigo.Conn
	if self.ctx != nil {
		conn = self.getConn(WithName(self.Name), WithContext(self.ctx))
	} else {
		conn = self.getConn(WithName(self.Name))
	}
	return conn
}

// 获得pool conn
func (self *RedisCmd) getConn(opt ...CmdOption) redigo.Conn {
	cmdOpt := newCmdOption(opt...)
	var conn redigo.Conn
	if cmdOpt.Ctx != nil {
		conn = self.Pool.GetWithCtx(cmdOpt.Ctx, cmdOpt.Name)
	} else {
		conn = self.Pool.Get(cmdOpt.Name)
	}
	return conn
}

func (self *RedisCmd) Get(key string) (string, error) {
	conn := self.GetConn()
	return redigo.String(conn.Do("GET", key))
}

func (self *RedisCmd) Hgetall(key string) (map[string]string, error) {
	conn := self.GetConn()

	log.Println("conn", conn)

	defer conn.Close()

	var data map[string]string
	data, err := redigo.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func (self *RedisCmd) Set(key string, value string) error {
	conn := self.GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}
	return err
}

func (self *RedisCmd) Del(key string) error {
	conn := self.GetConn()
	defer conn.Close()

	_, err := conn.Do("Del", key)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}
	return err
}

func (self *RedisCmd) Setex(key string, time int, value string) error {
	conn := self.GetConn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, time, value)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}
	return err
}

func (self *RedisCmd) Expire(key string, value int) error {
	conn := self.GetConn()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, value)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}
	return err
}

func (self *RedisCmd) Hmset(key string, some interface{}) error {
	conn := self.GetConn()
	defer conn.Close()

	if _, err := conn.Do("HMSET", redigo.Args{}.Add(key).AddFlat(&some)...); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error setting key %s to %v", key, err)
	}

	return nil
}

func (self *RedisCmd) Exists(key string) (bool, error) {
	conn := self.GetConn()
	defer conn.Close()

	ok, err := redigo.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func (self *RedisCmd) GetKeys(pattern string) ([]string, error) {
	conn := self.GetConn()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redigo.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redigo.Int(arr[0], nil)
		k, _ := redigo.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (self *RedisCmd) Incr(counterKey string) (int, error) {
	conn := self.GetConn()
	defer conn.Close()

	return redigo.Int(conn.Do("INCR", counterKey))
}
