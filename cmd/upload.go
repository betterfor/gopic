package cmd

import (
	"fmt"
	"github.com/betterfor/gopic/core"
	"github.com/betterfor/gopic/core/resize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

type upload struct {
	out          io.Writer
	rename       bool
	kind         string
	compress     string
	compressSize int
}

func newUploadCmd(out io.Writer) *cobra.Command {
	u := &upload{out: out}
	cmd := &cobra.Command{
		Use:     "upload",
		Short:   "upload local file to remote before completing setting",
		Example: "gopic upload testdata/helloworld.png",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				u.run(arg)
			}
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&u.rename, "rename", "r", false, "rename upload file to timestamp")
	f.StringVarP(&u.kind, "source", "k", "", "select one way to upload")
	f.StringVarP(&u.compress, "compress", "c", string(cfg.Base.CompressType), "choose one compress type to use")
	f.IntVarP(&u.compressSize, "size", "s", cfg.Base.CompressSize, "must use --compress before")
	return cmd
}

func (u *upload) run(file string) {
	var err error
	now := time.Now()
	if len(u.kind) != 0 {
		cfg.Current = u.kind
	}
	opts := uploadKind(cfg)

	var fileName string
	if cfg.Base.AutoRename || u.rename {
		fileName = fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file))
	} else {
		fileName = filepath.Base(file)
	}
	if debug {
		fmt.Fprintf(u.out, "upload file:%s, rename:%s\n", file, fileName)
	}

	var bts []byte
	if len(u.compress) != 0 {
		// 添加压缩
		bts, err = resize.Compress(resize.CompressType(u.compress), u.compressSize, file)
		if err != nil {
			fmt.Fprintf(u.out, "compress file:%s error:%v\n", file, err)
			return
		}
	} else {
		bts, err = ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(u.out, "read file:%s error:%v\n", file, err)
			return
		}
	}

	result, err := opts.Upload(fileName, bts)
	if err != nil {
		fmt.Fprintf(u.out, "upload file:%s error:%v\n", file, err)
		return
	}
	if debug {
		fmt.Fprintf(u.out, "consume:%v, upload file response %s\n", time.Since(now).String(), result)
	}
	url := opts.Parse(result)
	fmt.Fprintln(u.out, url)

	ret := viper.GetStringSlice("uploaded")
	viper.Set("uploaded", append(ret, url))
}

func uploadKind(cfg *core.Config) core.PicUpload {
	switch cfg.Current {
	case core.Github:
		return &cfg.Github
	case core.Smms:
		return &cfg.Smms
	case core.Gitee:
		return &cfg.Gitee
	default:
		return nil
	}
}
