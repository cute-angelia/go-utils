package spider2

import (
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
	"net/http/cookiejar"
)

type SpiderOptions struct {
	Name    string
	Client  *http.Client
	Cookies []*http.Cookie
	SOCKS5  string
	Debug   bool
}

type SpiderOption func(*SpiderOptions)

func NewSpiderOption(opts ...SpiderOption) *SpiderOptions {
	var sopt SpiderOptions
	for _, opt := range opts {
		opt(&sopt)
	}
	if sopt.Name == "" {
		sopt.Name = "default"
	}

	if sopt.Client == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			log.Fatal(err)
		}
		sopt.Client = &http.Client{
			Jar: jar,
		}
	}

	return &sopt
}

func WithName(name string) SpiderOption {
	return func(options *SpiderOptions) {
		options.Name = name
	}
}

func WithClient(Client *http.Client) SpiderOption {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	if Client == nil {
		Client = &http.Client{
			Jar: jar,
		}
	} else {
		Client.Jar = jar
	}

	return func(options *SpiderOptions) {
		options.Client = Client
	}
}

func WithCookies(Cookies []*http.Cookie) SpiderOption {
	return func(options *SpiderOptions) {
		options.Cookies = Cookies
	}
}

func WithSock5(SOCKS5 string) SpiderOption {
	return func(options *SpiderOptions) {
		options.SOCKS5 = SOCKS5
	}
}

func WithDebug(Debug bool) SpiderOption {
	return func(options *SpiderOptions) {
		options.Debug = Debug
	}
}
