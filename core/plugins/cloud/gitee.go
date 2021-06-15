package cloud

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// api: https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoContentsPath
type GiteeOpts struct {
	RepoName     string `json:"repoName" yaml:"repoName"`
	Branch       string `json:"branch" yaml:"branch"`
	Path         string `json:"path" yaml:"path"`
	Email        string `json:"email" yaml:"email"`
	Password     string `json:"password" yaml:"password"`
	ClientID     string `json:"clientId" yaml:"clientId"`
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
	AccessToken  string `json:"accessToken" yaml:"accessToken"`
}

func (g *GiteeOpts) URL() string {
	baseURL := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/contents/", g.RepoName)
	if len(g.Path) == 0 {
		return baseURL
	}
	// check if the path contains variables
	for {
		t0 := strings.Index(g.Path, "{")
		t1 := strings.Index(g.Path, "}")
		if t1-t0 > 0 {
			now := time.Now()
			g.Path = strings.Replace(g.Path, g.Path[t0:t1+1], now.Format(g.Path[t0+1:t1]), 1)
		} else {
			break
		}
	}

	return baseURL + g.Path + "/"
}

// request of github creating or updating file contents
type giteeRequest struct {
	AccessToken string `json:"access_token"` //
	Message     string `json:"message"`      // The commit message.
	Content     string `json:"content"`      // The new file content, using Base64 encoding.
	Branch      string `json:"branch"`       // The branch name. Default: the repository’s default branch (usually master)
}

func (g *giteeRequest) String() string {
	bts, _ := json.Marshal(g)
	return string(bts)
}

// POST https://gitee.com/api/v5/repos/{owner}/{repo}/contents/{path}
func (g *GiteeOpts) Upload(fileName string, data []byte) (string, error) {
	err := g.check()
	if err != nil {
		return "", err
	}
	// make body
	now := time.Now()
	ret := base64.StdEncoding.EncodeToString(data)

	var token string
	if len(g.AccessToken) != 0 {
		token, err = g.getToken()
		if err != nil {
			return "", err
		}
	} else {
		token = g.AccessToken
	}
	body := giteeRequest{
		AccessToken: token,
		Message:     fmt.Sprintf("upload file:%s at %s", fileName, now.String()),
		Content:     ret, // base64
		//Sha:     fmt.Sprintf("%x", md5.Sum(content)),        // md5
		Branch: g.Branch,
	}

	req, err := http.NewRequest(http.MethodPost, g.URL()+fileName, strings.NewReader(body.String()))
	if err != nil {
		return "", err
	}
	// set header
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// parse result
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		bts, err := ioutil.ReadAll(resp.Body)
		return string(bts), err
	case http.StatusUnauthorized:
		token, err = g.getToken()
		if err != nil {
			return "", err
		}
		viper.Set("gitee.accessToken", token)
	}
	return "", errors.New(resp.Status)
}

func (g *GiteeOpts) Parse(str string) string {
	return gjson.Get(str, "content.download_url").String()
}

func (g *GiteeOpts) check() error {
	if len(g.RepoName) == 0 {
		return errors.New("repo name is empty")
	}
	if len(g.Branch) == 0 {
		g.Branch = "master"
	}
	if len(g.Path) == 0 {
		g.Path = "images"
	}
	if len(g.AccessToken) == 0 {
		return errors.New("invalid token")
	}
	return nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

func (g *GiteeOpts) getToken() (string, error) {
	var m = map[string]string{
		"grant_type":    "password",
		"username":      g.Email,
		"password":      g.Password,
		"client_id":     g.ClientID,
		"client_secret": g.ClientSecret,
		"scope":         "projects",
	}
	v := url.Values{}
	for key, val := range m {
		v.Add(key, val)
	}
	resp, err := http.Post("https://gitee.com/oauth/token", "application/x-www-form-urlencoded", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
