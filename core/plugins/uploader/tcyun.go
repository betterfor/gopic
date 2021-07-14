package uploader

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/tidwall/gjson"
	"hash"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"
)

type TcyunConfig struct {
	SecretId  string `json,yaml:"secretId"`
	SecretKey string `json,yaml:"secretKey"`
	Bucket    string `json,yaml:"bucket"`
	AppId     string `json,yaml:"appId"`
	Area      string `json,yaml:"area"`
	Path      string `json,yaml:"path"`
	CustomUrl string `json,yaml:"customUrl"`
}

func (t *TcyunConfig) Name() string {
	return "腾讯云COS"
}

func (t *TcyunConfig) Upload(img *Img) error {
	var req *http.Request

	url := fmt.Sprintf("http://%s.cos.%s.myqcloud.com/%s/%s", t.Bucket, t.Area, t.Path, img.FileName)
	req, err := http.NewRequest(http.MethodPut, url, img.Reader)
	if err != nil {
		return err
	}

	typ := mime.TypeByExtension(img.FileName)

	// sign: https://cloud.tencent.com/document/product/436/7778
	signTime, sign := t.generateSignature(img.FileName)
	auth := strings.Join([]string{
		"q-sign-algorithm=sha1",
		"q-ak=" + t.SecretId,
		"q-sign-time=" + signTime,
		"q-key-time=" + signTime,
		"q-header-list=host",
		"q-url-param-list=",
		"q-signature=" + sign,
	}, "&")
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", typ)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if len(t.CustomUrl) != 0 {
			img.ImgUrl = fmt.Sprintf("%s/%s/%s", t.CustomUrl, t.Path, img.FileName)
		} else {
			img.ImgUrl = fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s/%s", t.Bucket, t.Area, t.Path, img.FileName)
		}
		return nil
	default:
		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(gjson.GetBytes(bts, "msg").String())
	}
}

func (t *TcyunConfig) generateSignature(filename string) (string, string) {
	now := time.Now()
	start := now.Unix()
	end := start + 86400
	signTime := fmt.Sprintf("%d:%d", start, end)

	signKey := calHMACDigest(t.SecretKey, signTime)
	httpString := fmt.Sprintf("put\n/%s/%s\n\nhost=%s.cos.%s.myqcloud.com\n", t.Path, filename, t.Bucket, t.Area)
	sha1edHttpString := calHMACDigest("", httpString)
	stringToSign := fmt.Sprintf("sha1\n%s\n%s\n", signTime, sha1edHttpString)
	signature := calHMACDigest(signKey, stringToSign)
	return signTime, signature
}

func calHMACDigest(key, msg string) string {
	h := hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte(key))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))
}
