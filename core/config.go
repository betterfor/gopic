package core

import (
	"github.com/betterfor/gopic/core/plugins/uploader"
	"github.com/betterfor/gopic/core/resize"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// Config is gopic config
type Config struct {
	Uploaded []string           // uploaded pictures
	Base     BaseConfig         // base config
	Current  uploader.GopicType // current use picbed
	Data     map[uploader.GopicType]interface{}
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

func (c *Config) Use() uploader.Uploader {
	data := c.Data[c.Current]
	switch c.Current {
	case uploader.Github:
		var c uploader.GithubConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Gitee:
		var c uploader.GiteeConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Smms:
		var c uploader.SmmsConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Qiniu:
		var c uploader.QiniuConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Imgur:
		var c uploader.ImgurConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Aliyun:
		var c uploader.AliyunConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Upyun:
		var c uploader.UpyunConfig
		mapstructure.Decode(data, &c)
		return &c
	case uploader.Tcyun:
		var c uploader.TcyunConfig
		mapstructure.Decode(data, &c)
		return &c
	}
	return nil
}
