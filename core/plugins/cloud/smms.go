package cloud

import (
	"bytes"
	"errors"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// https://doc.sm.ms/#api-Image-Upload
// smms options
type SmmsOpts struct {
	Token string `json:"token"`
}

func (s *SmmsOpts) URL() string {
	return "https://sm.ms/api/v2/upload"
}

func (s *SmmsOpts) Upload(fileName string, data []byte) (string, error) {
	// read file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("smfile", fileName)
	if err != nil {
		return "", err
	}
	formFile.Write(data)

	//
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, s.URL(), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	t := oauth2.Token{AccessToken: s.Token, TokenType: "Basic"}
	t.SetAuthHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		bts, err := ioutil.ReadAll(resp.Body)
		return string(bts), err
	default:
		return "", errors.New(resp.Status)
	}
}

func (s *SmmsOpts) Parse(str string) string {
	if gjson.Get(str, "success").Bool() {
		return gjson.Get(str, "data.url").String()
	} else {
		return gjson.Get(str, "code").String() + gjson.Get(str, "message").String()
	}
}
