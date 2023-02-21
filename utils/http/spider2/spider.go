package spider2

import (
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"log"
)

type Spider struct {
	SpiderOptions *SpiderOptions
}

func NewSpider(SpiderOptions *SpiderOptions) *Spider {
	return &Spider{SpiderOptions: SpiderOptions}
}

func (self Spider) getOut() *dataflow.Gout {
	g := gout.New(self.SpiderOptions.Client)

	if len(self.SpiderOptions.SOCKS5) > 0 {
		g.SetSOCKS5(self.SpiderOptions.SOCKS5)
	}

	if self.SpiderOptions.Debug {
		g.Debug(true)
	}

	g.SetCookies(self.SpiderOptions.Cookies...)

	return g
}

func (self Spider) Get(url string, header gout.H) ([]byte, error) {
	g := self.getOut().GET(url)
	s := []byte{}
	err := g.SetHeader(header).BindBody(&s).Do()

	if err != nil {
		log.Println("spider get error : -> ", err)
		return nil, err
	}
	return s, nil
}

func (self Spider) Post(url string, header gout.H, data gout.H) ([]byte, error) {
	s := []byte{}
	err := self.getOut().POST(url).SetHeader(header).SetWWWForm(data).BindBody(&s).Do()
	if err != nil {
		log.Println("spider get error : -> ", err)
		return nil, err
	}
	return s, nil
}
