package spider

import (
	"io/ioutil"
	"net/http"
	"log"
	"fmt"
	"strings"
	"net/url"
	"mime/multipart"
	"os"
	"io"
	"bytes"
)

func Get(client *http.Client, url string, cookie []*http.Cookie, header http.Header) ([]byte, []*http.Cookie, error) {
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

func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	defer w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	//req.Header = header


	log.Println("\n")
	log.Println(req)
	log.Println("\n")

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	content, _ := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	log.Println(string(content))

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}
