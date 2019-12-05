package snowflake

import (
	"github.com/sony/sonyflake"
	"io/ioutil"
	"net/http"
	"net"
	"errors"
	"fmt"
)

func NewSnowId(machineId bool) (id uint64, error error) {
	var st sonyflake.Settings
	if machineId {
		st.MachineID = MachineIP
	}
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return 0, fmt.Errorf("sonyflake not created")
	}
	id, err := sf.NextID()
	return id, err
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
