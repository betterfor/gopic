package uploader

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"hash"
	"io"
	"mime"
	"net/http"
	"time"
)

type AliyunConfig struct {
	Bucket          string `json,yaml:"bucket"`          // storage bucket
	Area            string `json,yaml:"area"`            // storage area
	AccessKeyId     string `json,yaml:"accessKeyId"`     // key
	AccessKeySecret string `json,yaml:"accessKeySecret"` // secret
	CustomUrl       string `json,yaml:"customUrl"`       // domain
	Path            string `json,yaml:"path"`            // storage path
}

func (a *AliyunConfig) Name() string {
	return "阿里云OSS"
}

func (a *AliyunConfig) Upload(img *Img) error {

	typ := mime.TypeByExtension(img.FileName)
	if len(typ) == 0 {
		typ = "application/octet-stream"
	}
	signature := generateSignature(a, img.FileName, typ)

	url := fmt.Sprintf("https://%s.%s.aliyuncs.com/%s/%s", a.Bucket, a.Area, a.Path, img.FileName)

	req, err := http.NewRequest(http.MethodPut, url, img.Reader)
	if err != nil {
		return err
	}
	t := oauth2.Token{AccessToken: signature, TokenType: "OSS"}
	t.SetAuthHeader(req)
	req.Header.Add("Date", time.Now().Format(http.TimeFormat))
	req.Header.Add("Content-Type", typ)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if len(a.CustomUrl) == 0 {
			img.ImgUrl = url
		} else {
			img.ImgUrl = fmt.Sprintf("%s/%s/%s", a.CustomUrl, a.Path, img.FileName)
		}
		return nil
	default:
		return fmt.Errorf("Upload failed ")
	}
}

func generateSignature(options *AliyunConfig, fileName, contentType string) string {
	now := time.Now()
	signStr := fmt.Sprintf("PUT\n\n%s\n%s\n%s", contentType, now.Format(http.TimeFormat), options.Bucket+"/"+options.Area+"/"+fileName)
	h := hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte(options.AccessKeySecret))

	io.WriteString(h, signStr)

	return fmt.Sprintf("%s:%s", options.AccessKeyId, base64.StdEncoding.EncodeToString(h.Sum(nil)))
}
