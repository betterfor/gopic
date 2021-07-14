package uploader

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type ImgurConfig struct {
	ClientId string `json,yaml:"clientId"`
	Proxy    string `json,yaml:"proxy"`
}

func (i *ImgurConfig) Name() string {
	return "Imgur图床"
}

func (i *ImgurConfig) Upload(img *Img) error {
	u := "https://api.imgur.com/3/image"

	buf, err := io.ReadAll(img.Reader)
	if err != nil {
		return err
	}
	postData := map[string]string{
		"image": base64.StdEncoding.EncodeToString(buf),
		"type":  "base64",
		"name":  img.FileName,
	}

	buffer := new(bytes.Buffer)
	w := multipart.NewWriter(buffer)
	for k, v := range postData {
		w.WriteField(k, v)
	}
	w.Close()
	req, err := http.NewRequest(http.MethodPost, u, buffer)
	if err != nil {
		return err
	}
	t := oauth2.Token{AccessToken: i.ClientId, TokenType: "Client-ID"}
	t.SetAuthHeader(req)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("User-Agent", "GoPic")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	if len(i.Proxy) != 0 {
		proxyUrl, err := url.Parse(i.Proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if gjson.GetBytes(bts, "success").Bool() {
			img.ImgUrl = gjson.GetBytes(bts, "data.link").String()
			return nil
		}
	}
	return fmt.Errorf("Server error,please try again ")
}
