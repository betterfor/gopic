package cmd

import (
	"fmt"
	"github.com/betterfor/gopic/core"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

type upload struct {
	out io.Writer
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
	return cmd
}

func (u *upload) run(file string) {
	opts := uploadKind(cfg)
	fmt.Println(opts)

	var fileName string
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(u.out, "read file:%s error:%v", file, err)
		return
	}
	if cfg.Base.AutoRename {
		fileName = fmt.Sprintf("%d.%s", time.Now().UnixNano(), filepath.Ext(file))
	} else {
		fileName = filepath.Base(file)
	}
	url, err := opts.Upload(fileName, bts)
	if err != nil {
		fmt.Fprintf(u.out, "upload file:%s error:%v", file, err)
		return
	} else {
		fmt.Fprintln(u.out, url)
	}

	cfg.Uploaded = append(cfg.Uploaded, url)
}

func uploadKind(cfg *core.Config) core.PicUpload {
	switch cfg.Current {
	case core.Github:
		return &cfg.Github
	default:
		return nil
	}
}
