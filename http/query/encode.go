package query

import (
	"fmt"
	"net/url"
	"strings"
)

/*
	Encode Url String
	one ->
		url.QueryEscape(text)
	multi ->
		params := url.Values{}
		params.Add("q", "1 + 2")
		params.Add("s", "example for golangcode.com")
		output := params.Encode()
*/
func EncodeUrl(uri string) (string, error) {
	if strings.Contains(uri, "http") {
		if l2, err := url.Parse(uri); err != nil {
			return uri, err
		} else {
			return fmt.Sprintf("%s://%s%s?%s", l2.Scheme, l2.Host, l2.Path, l2.Query().Encode()), nil
		}
	} else {
		if l, err := url.ParseQuery(uri); err != nil {
			return uri, err
		} else {
			return l.Encode(), nil
		}
	}
}
