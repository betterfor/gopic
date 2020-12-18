package core

import "encoding/json"

type picturesOpt int

const (
	Github picturesOpt = iota // github
	Smms
	tcyun
	qiniu
	imgur
	aliyun
	upyun
)

// 上传图片接口
type PicUpload interface {
	Upload(fileName string, content []byte) (string, error)
	String() string
}

// GOPIC config
type Config struct {
	Uploaded []string `json:"uploaded"` // uploaded pictures
	Base     Base     `json:"base"`     // base config
	PicBed
}

// PicBed options
type PicBed struct {
	Current  string               `json:"current"`  // current use picbed
	Settings map[string]PicUpload `json:"settings"` // pic bed details
}

type Base struct {
	AutoRename bool `json:"autoRename" yaml:"autoRename"` // use timestamp rename file
}

func (c *Config) String() string {
	bts, _ := json.Marshal(c)
	return string(bts)
}
