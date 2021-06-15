package cdn

import "testing"

func TestJsDelivr_Convert(t *testing.T) {
	type fields struct {
		Url          string
		Description  string
		Prefix       string
		GithubPrefix string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"test1",
			fields{
				Url:          "https://www.jsdelivr.com/",
				Description:  "a free CDN for Open Source, fast, reliable, and automated",
				Prefix:       "https://cdn.jsdelivr.net/gh/",
				GithubPrefix: "https://raw.githubusercontent.com/",
			},
			args{url: "https://raw.githubusercontent.com/betterfor/cloudImage/master/images/2021/01/08/rocketmq.png"},
			"https://cdn.jsdelivr.net/gh/betterfor/cloudImage/images/2021/01/08/rocketmq.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JsDelivr{
				Url:          tt.fields.Url,
				Description:  tt.fields.Description,
				Prefix:       tt.fields.Prefix,
				GithubPrefix: tt.fields.GithubPrefix,
			}
			if got := j.Convert(tt.args.url); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
