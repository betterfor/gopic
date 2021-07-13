package uploader

import (
	"bytes"
	"errors"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// SmmsConfig sm.ms options
// api: https://doc.sm.ms/#api-Image-Upload
type SmmsConfig struct {
	Token string `json:"token"`
}

func (s *SmmsConfig) Name() string {
	return "SM.MS图床"
}

func (s *SmmsConfig) Upload(img *Img) error {
	url := "https://sm.ms/api/v2/upload"
	// read file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("smfile", img.FileName)
	if err != nil {
		return err
	}
	io.Copy(formFile, img.Reader)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	t := oauth2.Token{AccessToken: s.Token, TokenType: "Basic"}
	t.SetAuthHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		code := gjson.GetBytes(bts, "code").String()
		if code == "success" {
			img.ImgUrl = gjson.GetBytes(bts, "data.url").String()
		} else if code == "image_repeated" {
			img.ImgUrl = gjson.GetBytes(bts, "images").String()
		} else {
			return errors.New(gjson.GetBytes(bts, "message").String())
		}
		return nil
	default:
		return errors.New(resp.Status)
	}
}

func (s *SmmsConfig) Parse(str string) string {
	if gjson.Get(str, "success").Bool() {
		return gjson.Get(str, "data.url").String()
	} else {
		return gjson.Get(str, "code").String() + gjson.Get(str, "message").String()
	}
}
