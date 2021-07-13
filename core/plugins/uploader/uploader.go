package uploader

import (
	"io"
	"net/http"
	"time"
)

// GopicType is uploader name
type GopicType string

const (
	Unknown GopicType = ""
	Github  GopicType = "github"
	Gitee   GopicType = "gitee"
	Smms    GopicType = "smms"
	Qiniu   GopicType = "qiniu"
	Imgur   GopicType = "imgur"
	Aliyun  GopicType = "aliyun"
	Upyun   GopicType = "upyun"
	Tcyun   GopicType = "tcyun"
)

// Uploader upload image
type Uploader interface {
	Name() string
	Upload(img *Img) error
}

type Img struct {
	FileName string
	Reader   io.Reader
	ImgUrl   string
}

const timeout = time.Second * 30

var client = &http.Client{
	Timeout: timeout,
}

func Convert(v string) GopicType {
	switch v {
	case "github":
		return Github
	case "gitee":
		return Gitee
	case "smms":
		return Smms
	case "qiniu":
		return Qiniu
	case "imgur":
		return Imgur
	case "aliyun":
		return Aliyun
	case "upyun":
		return Upyun
	case "tcyun":
		return Tcyun
	default:
		return Unknown
	}
}
