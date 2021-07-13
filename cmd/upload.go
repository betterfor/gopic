package cmd

import (
	"fmt"
	"github.com/betterfor/gopic/core/plugins/uploader"
	"github.com/betterfor/gopic/core/resize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type upload struct {
	out          io.Writer
	rename       bool
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
	f.StringVarP(&u.compress, "compress", "c", string(cfg.Base.CompressType), "choose one compress type to use")
	f.IntVarP(&u.compressSize, "size", "s", cfg.Base.CompressLevel, "must use --compress before")
	return cmd
}

func (u *upload) run(file string) {
	var err error
	now := time.Now()

	// download network file to local
	if strings.HasPrefix(file, "http") || strings.HasPrefix(file, "ftp") {
		resp, err := http.Get(file)
		if err != nil {
			fmt.Fprintf(u.out, "download file:%s error:%v\n", file, err)
			return
		}
		fi, err := os.CreateTemp(".", "*.png")
		if err != nil {
			fmt.Fprintf(u.out, "create tmp file error:%v\n", err)
			return
		}
		_, err = io.Copy(fi, resp.Body)
		if err != nil {
			fmt.Fprintf(u.out, "create tmp file error:%v\n", err)
			return
		}
		resp.Body.Close()
		fi.Close()
		file = fi.Name()
		defer func() {
			err = os.Remove(file)
			if err != nil {
				fmt.Fprintf(u.out, "remove tmp file:%s error:%v\n", file, err)
				return
			}
		}()
	}

	var fileName string
	if cfg.Base.AutoRename || u.rename {
		fileName = fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file))
	} else {
		fileName = filepath.Base(file)
	}
	if debug {
		fmt.Fprintf(u.out, "upload file:%s, rename:%s\n", file, fileName)
	}

	var reader io.Reader
	if len(u.compress) != 0 {
		// 添加压缩
		reader, err = resize.Compress(resize.CompressType(u.compress), u.compressSize, file)
		if err != nil {
			fmt.Fprintf(u.out, "compress file:%s error:%v\n", file, err)
			return
		}
	} else {
		reader, err = os.Open(file)
		if err != nil {
			fmt.Fprintf(u.out, "read file:%s error:%v\n", file, err)
			return
		}
	}

	img := &uploader.Img{
		FileName: fileName,
		Reader:   reader,
	}
	up := cfg.Use()
	err = up.Upload(img)
	if err != nil {
		fmt.Fprintf(u.out, "upload failed"+err.Error())
		return
	}

	fn, err := os.Open(fileName)
	ioutil.ReadAll(fn)

	imgUrl := img.ImgUrl
	if debug {
		fmt.Fprintf(u.out, "consume:%v, upload file response %s\n", time.Since(now).String(), imgUrl)
	}
	fmt.Fprintln(u.out, imgUrl)

	ret := viper.GetStringSlice("uploaded")
	viper.Set("uploaded", append(ret, imgUrl))
}
