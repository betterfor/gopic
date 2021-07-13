package resize

import (
	"fmt"
	"image"
	"io"
	"os"
	"testing"
)

func Test_getReadSizeFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantImg image.Image
		wantErr bool
	}{
		{
			name: "test0",
			args: struct{ file string }{file: "../plugins/testdata/helloworld.png"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := getReadSizeFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("getReadSizeFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getReadSizeFile() gotX = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompress(t *testing.T) {
	type args struct {
		t    CompressType
		base int
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "compress0",
			args: struct {
				t    CompressType
				base int
				file string
			}{t: CompressAuto, base: 10240, file: "../plugins/testdata/helloworld.png"},
		},
		{
			name: "compress1",
			args: struct {
				t    CompressType
				base int
				file string
			}{t: CompressCustom, base: 2, file: "../plugins/testdata/helloworld.png"},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.args.t, tt.args.base, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			f, _ := os.OpenFile(fmt.Sprintf("%d.jpg", i), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			_, err = io.Copy(f, got)
			if err != nil {
				t.Errorf("Write file error %v", err)
			}
		})
	}
}
