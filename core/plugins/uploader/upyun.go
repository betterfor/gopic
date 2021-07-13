package uploader

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type UpyunConfig struct {
	Bucket   string `json:"bucket" yaml:"bucket"`
	Operator string `json:"operator" yaml:"operator"`
	Password string `json:"password" yaml:"password"`
	Url      string `json:"url" yaml:"url"`
	Options  string `json:"options" yaml:"options"`
	Path     string `json:"path" yaml:"path"`
}

func (u *UpyunConfig) Name() string {
	return "又拍云图床"
}

func (u *UpyunConfig) Upload(img *Img) error {
	sign := u.generateSignature(img.FileName)

	url := fmt.Sprintf("https://v0.api.upyun.com/%s/%s/%s", u.Bucket, u.Password, img.FileName)
	req, err := http.NewRequest(http.MethodPut, url, img.Reader)
	if err != nil {
		return err
	}
	t := oauth2.Token{AccessToken: u.Operator + ":" + sign, TokenType: "UPYUN"}
	t.SetAuthHeader(req)
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		img.ImgUrl = fmt.Sprintf("%s/%s/%s%s", u.Url, u.Password, img.FileName, u.Options)
		return nil
	default:
		return fmt.Errorf("uploaded failed")
	}
}

func (u UpyunConfig) generateSignature(filename string) string {
	// generate token: https://help.upyun.com/knowledge-base/object_storage_authorization/#token-e8aea4e8af81
	h := md5.New()
	h.Write([]byte(u.Password))
	md5Password := fmt.Sprintf("%x", h.Sum(nil))
	uri := fmt.Sprintf("%s/%s/%s", u.Bucket, u.Path, filename)
	now := time.Now()
	value := fmt.Sprintf("PUT&%s&%d", uri, now.Unix())
	sign := calHMACDigest(md5Password, value)
	return base64.StdEncoding.EncodeToString([]byte(sign))
}
