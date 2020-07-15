package checkServer

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func CheckServer(host string) error {
	timeout := time.Duration(8 * time.Second)
	// t1 := time.Now()
	// url := "www.google.com:443"
	_, err := net.DialTimeout("tcp", host, timeout)
	// fmt.Println("waist time :", time.Now().Sub(t1))
	if err != nil {
		// fmt.Println("Site unreachable, error: ", err)
		return fmt.Errorf("Site unreachable, error:  %s, %s", err.Error(), host)
	}
	return nil
}

func CheckServerHttp(uri string) error {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	// t1 := time.Now()
	resp, err := client.Get(uri)
	// fmt.Println("waist time :", time.Now().Sub(t1))
	if err != nil {
		return fmt.Errorf("Site unreachable, error:  %s, %s", err.Error(), uri)
	}
	defer resp.Body.Close()
	// fmt.Println(resp.Status)
	return nil
}
