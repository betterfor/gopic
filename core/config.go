package core

import (
	"github.com/betterfor/gopic/core/plugins"
	"github.com/betterfor/gopic/core/resize"
	"gopkg.in/yaml.v2"
)

const (
	Github = "github"
	Gitee  = "gitee"
	Smms   = "smms"
	Qiniu
	Imgur
	Aliyun
	Upyun
)

// 上传图片接口
type PicUpload interface {
	Upload(fileName string, data []byte) (string, error)
	Parse(str string) string
}

// GOPIC config
type Config struct {
	Uploaded []string // uploaded pictures
	Base     Base     // base config
	Current  string   // current use picbed
	Github   plugins.GithubOpts
	Gitee    plugins.GiteeOpts
	Smms     plugins.SmmsOpts
}

type Base struct {
	AutoRename   bool                // use timestamp rename file
	CompressType resize.CompressType // compress kind
	CompressSize int                 // compress times
}

func (c *Config) String() string {
	bts, _ := yaml.Marshal(c)
	return string(bts)
}
