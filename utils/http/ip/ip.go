package ip

import (
	"encoding/binary"
	"net"
	"net/http"
	"strings"
)

func GetInfo(ip string) map[string]ResultQQwry {
	ips := strings.Split(ip, ",")
	qqWry := NewQQwry()
	rs := map[string]ResultQQwry{}
	if len(ips) > 0 {
		for _, v := range ips {
			rs[v] = qqWry.Find(v)
		}
	}

	return rs
}

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
// nginx
//    location /api-auth {
//        proxy_set_header X-Forwarded-For $remote_addr; 一级代理
//    }

//     location /api-auth {
//        proxy_set_header X-Forwarded-For $http_x_forwarded_for; // 二级代理获取头部配置
//        proxy_set_header X-Real-IP $remote_addr; // 直接代理获取数据
//        proxy_pass http://127.0.0.1:9104;
//    }
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func RemoteIp2Uint(req *http.Request) uint32 {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return Ip2long(remoteAddr)
}

// Ip2long 将 IPv4 字符串形式转为 uint32
func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func Long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}
