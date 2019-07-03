package spider

import (
	"io/ioutil"
	"net/http"
	"log"
	"fmt"
	"strings"
	"net/url"
)

func Get(client *http.Client, url string, cookie []*http.Cookie, header http.Header) (string, []*http.Cookie, error) {
	requestGet, _ := http.NewRequest(
		"GET",
		url, nil)

	requestGet.Header = header

	// add cookies
	SetCookie(cookie, requestGet)

	// log.Println(requestGet)

	//reqLocation.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//reqLocation.Header.Set("User-Agent", d.UA)

	if respGet, err := client.Do(requestGet); err != nil {
		return "", cookie, err
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
		str := string(content)

		return str, cookie, nil
	}
}

func Post(client *http.Client, url string, data url.Values, cookie []*http.Cookie, header http.Header) (string, []*http.Cookie, error) {
	req, _ := http.NewRequest(
		"POST",
		url, strings.NewReader(data.Encode()))

	req.Header = header

	// 设置 cookie
	SetCookie(cookie, req)


	// log.Println(strings.NewReader(data.Encode()))
	// log.Println(req)

	respPost, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		return "", cookie, err
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
	str := string(content)

	return str, cookie, nil
}
