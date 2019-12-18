package file

import (
	"io/ioutil"
	"os"
	"net/http"
	"crypto/tls"
	"time"
)

// Net read net file
func GetFileWithSrc(src string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: time.Second * 6}
	// set request
	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	// req.Header = args.Header
	// get response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Local read local file
func GetFileWithLocal(path string) ([]byte, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(imageFile)
}
