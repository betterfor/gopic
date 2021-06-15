package core

import (
	"github.com/betterfor/gopic/core/plugins/cloud"
	"github.com/betterfor/gopic/core/resize"
	"gopkg.in/yaml.v2"
)

type GopicType string

const (
	Github GopicType = "github"
	Gitee  GopicType = "gitee"
	Smms   GopicType = "smms"
	Qiniu  GopicType = "qiniu"
	Imgur  GopicType = "imgur"
	Aliyun GopicType = "aliyun"
	Upyun  GopicType = "upyun"
)

// 上传图片接口
type PicUpload interface {
	Upload(fileName string, data []byte) (string, error)
	Parse(str string) string
}

// Config is gopic config
type Config struct {
	Uploaded []string   // uploaded pictures
	Base     BaseConfig // base config
	Current  GopicType  // current use picbed
	Github   cloud.GithubOpts
	Gitee    cloud.GiteeOpts
	Smms     cloud.SmmsOpts
}

type BaseConfig struct {
	AutoRename    bool                // use timestamp name replace file name
	CompressType  resize.CompressType // compress kind
	CompressLevel int                 // compress level
}

func (c *Config) String() string {
	bts, _ := yaml.Marshal(c)
	return string(bts)
}
