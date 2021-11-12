package spider

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Get(client *http.Client, url string, cookie []*http.Cookie, header http.Header) ([]byte, []*http.Cookie, error) {
	if requestGet, err := http.NewRequest(
		"GET",
		url, nil); err != nil {
		log.Println("Get:", err)

		return []byte{}, nil, err
	} else {

		requestGet.Header = header

		// add cookies
		SetCookie(cookie, requestGet)

		// log.Println(requestGet)

		//reqLocation.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		//reqLocation.Header.Set("User-Agent", d.UA)

		if respGet, err := client.Do(requestGet); err != nil {
			return nil, cookie, err
		} else {
			defer respGet.Body.Close()

			// merger cookie
			if respGet != nil && len(respGet.Cookies()) > 0 {
				for _, z := range respGet.Cookies() {
					log.Println("新增 cookie", z.Name+"="+z.Value)
				}
				// merge
				cookie = MergeCookie(cookie, respGet.Cookies())

				// log
				LogCookie(fmt.Sprintf("url: %s \n", url), cookie)
			}

			// location
			if respGet != nil {
				defer func() {
					// location
					location := respGet.Header.Get("Location")
					if len(location) > 0 {
						log.Println(fmt.Sprintf("url: %s \n location: %s \n", url, location))
						reqLocation, _ := http.NewRequest(
							"GET",
							location, nil)

						// add cookies
						SetCookie(cookie, reqLocation)
						reqLocation.Header = header
						respLocation, _ := client.Do(reqLocation)
						defer respLocation.Body.Close()
						cookie = MergeCookie(cookie, respLocation.Cookies())
					}
				}()
			}

			content, _ := ioutil.ReadAll(respGet.Body)

			return content, cookie, nil
		}
	}
}

func Post(client *http.Client, url string, data url.Values, cookie []*http.Cookie, header http.Header) ([]byte, []*http.Cookie, error) {
	req, _ := http.NewRequest(
		"POST",
		url, strings.NewReader(data.Encode()))

	if len(header.Get("Content-Type")) == 0 {
		header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}

	req.Header = header

	// 设置 cookie
	SetCookie(cookie, req)

	respPost, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		return nil, cookie, err
	}

	defer respPost.Body.Close()

	// cookies
	cookie = MergeCookie(cookie, respPost.Cookies())

	// login
	// LogCookie(fmt.Sprintf("POST -> url: %s \n cookie:", url), cookie)

	// location
	if respPost != nil {
		defer func() {
			// location
			location := respPost.Header.Get("Location")
			if len(location) > 0 {
				log.Println(fmt.Sprintf("url: %s \n location: %s \n", url, location))
				reqLocation, _ := http.NewRequest(
					"GET",
					location, nil)

				// add cookies
				SetCookie(cookie, reqLocation)
				reqLocation.Header = header
				respLocation, _ := client.Do(reqLocation)
				defer respLocation.Body.Close()
				cookie = MergeCookie(cookie, respLocation.Cookies())
			}
		}()
	}

	content, _ := ioutil.ReadAll(respPost.Body)

	return content, cookie, nil
}

func Upload(client *http.Client, url string, params map[string]string, paramName, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// add params
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	if req, err := http.NewRequest("POST", url, body); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// log.Println("req", req)

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return nil, err
		} else {
			if content, err := ioutil.ReadAll(resp.Body); err != nil {
				log.Println(err)
				return nil, err
			} else {
				resp.Body.Close()
				if resp.StatusCode == 200 {
					log.Println(body)
					return content, nil
				} else {
					return nil, fmt.Errorf("错误:%d", resp.Status)
				}
			}
			//fmt.Println(resp.StatusCode)
			//fmt.Println(resp.Header)
		}
	}
}
