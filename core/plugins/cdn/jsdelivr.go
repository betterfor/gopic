package cdn

import "strings"

type JsDelivr struct {
	Url         string
	Description string
	Prefix      string
}

func NewJsDelivr() CDN {
	return &JsDelivr{
		Url:         "https://www.jsdelivr.com/",
		Description: "a free CDN for Open Source, fast,reliable, adn automated",
		Prefix:      "https://cdn.jsdelivr.net/gh/",
	}
}

func (j *JsDelivr) Convert(url string) string {
	// https://raw.githubusercontent.com/betterfor/cloudImage/master/images/2021/01/08/rocketmq.png
	// convert to
	// https://cdn.jsdelivr.net/gh/betterfor/cloudImage/images/2021/01/08/rocketmq.png
	url = strings.Replace(url, "https://raw.githubusercontent.com/", j.Prefix, 1)
	url = strings.Replace(url, "/master", "", 1)
	return url
}

func init() {
	Register(Jsdelivr, NewJsDelivr)
}
