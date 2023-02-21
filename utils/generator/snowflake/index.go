package snowflake

import (
	"github.com/sony/sonyflake"
	"io/ioutil"
	"net/http"
	"net"
	"errors"
	"fmt"
	"strconv"
)

var st sonyflake.Settings
var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(st)
}

type SnowRest struct {
	Id uint64
}

// Int64 returns an int64 of the snowflake ID
func (f SnowRest) Int64() int64 {
	return int64(f.Id)
}

// String returns a string of the snowflake ID
func (f SnowRest) String() string {
	return strconv.FormatInt(int64(f.Id), 10)
}

func NewSnowId(mode int64) (SnowRest, error) {
	z := SnowRest{}

	if sf == nil {
		return z, fmt.Errorf("sonyflake not created")
	}

	id, err := sf.NextID()
	z.Id = id
	return z, err
}

func NewSnowIdWithSetting(st sonyflake.Settings) (SnowRest, error) {
	sf := sonyflake.NewSonyflake(st)
	z := SnowRest{}

	if sf == nil {
		return z, fmt.Errorf("sonyflake not created")
	}
	id, err := sf.NextID()
	z.Id = id
	return z, err
}

func EC2PrivateIPv4() (net.IP, error) {
	res, err := http.Get("https://ident.me")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(body))
	if ip == nil {
		return nil, errors.New("invalid ip address")
	}
	return ip.To4(), nil
}

func MachineIP() (uint16, error) {
	ip, err := EC2PrivateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}
