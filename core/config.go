package core

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

// GOPIC config
type Config struct {
	Uploaded []string `json:"uploaded"` // uploaded pictures
	PicBed
}

// PicBed options
type PicBed struct {
	Current string `json:"current"`
}

type BaseConfig struct {
	AutoRename bool `json:"autoRename" yaml:"autoRename"` // 时间戳重命名
}
