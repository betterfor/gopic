package uploader

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
)

type QiniuConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
	Url       string `json:"url"`
	Area      string `json:"area"`
	Option    string `json:"option"`
	Path      string `json:"path"`
}

func (q *QiniuConfig) Name() string {
	return "七牛图床"
}

func (q *QiniuConfig) Upload(img *Img) error {
	area := getArea(q.Area)
	base64FileName := base64.StdEncoding.EncodeToString([]byte(q.Path + img.FileName))
	url := fmt.Sprintf("http://upload%s.qiniu.com/putb64/-1/key/%s", area, base64FileName)
	buf, err := io.ReadAll(img.Reader)
	if err != nil {
		return err
	}
	base64File := base64.StdEncoding.EncodeToString(buf)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(base64File))
	if err != nil {
		return err
	}
	t := oauth2.Token{AccessToken: getToken(q.Bucket, q.AccessKey, q.SecretKey), TokenType: "UpToken"}
	t.SetAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:

		img.ImgUrl = q.Url + "/" + gjson.GetBytes(bts, "key").String() + q.Option
		return nil
	default:
		return fmt.Errorf(gjson.GetBytes(bts, "msg").String())
	}
}

func getArea(area string) string {
	if len(area) == 0 {
		area = "z0"
	} else {
		area = "-" + area
	}
	return area
}

func getToken(bucket, key, secret string) string {
	putPolicy := storage.PutPolicy{Scope: bucket}
	mac := auth.New(key, secret)
	return putPolicy.UploadToken(mac)
}
