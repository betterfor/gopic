package core

import (
	"github.com/betterfor/gopic/core/plugins"
	"gopkg.in/yaml.v2"
)

const (
	Github = "github"
	Smms   = "smms"
	Tcyun
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
	Smms     plugins.SmmsOpts
}

type Base struct {
	AutoRename bool // use timestamp rename file
}

func (c *Config) String() string {
	bts, _ := yaml.Marshal(c)
	return string(bts)
}
