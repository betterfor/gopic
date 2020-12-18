package core

import (
	"encoding/json"
	"fmt"
	"github.com/betterfor/gopic/core/plugins"
)

const (
	Github = "github"
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

// custom structure
func (c *Config) UnmarshalJSON(data []byte) error {
	fmt.Println("unmarshal this")
	tmp := struct {
		Uploaded []string
		Base     Base `json:"base"`
		Current  string
		Settings map[string]interface{}
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	fmt.Println("unmarshal success")
	var sets = make(map[string]PicUpload)
	for key, settings := range tmp.Settings {
		switch key {
		case Github:
			sets[Github] = settings.(*plugins.GithubOpts)
		}
	}
	o.Settings = sets
	o.Current = tmp.Current
	o.Base = tmp.Base
	o.Uploaded = tmp.Uploaded
	return nil
}

type Base struct {
	AutoRename bool `json:"autoRename" yaml:"autoRename"` // use timestamp rename file
}

func (c *Config) String() string {
	bts, _ := json.Marshal(c)
	return string(bts)
}
