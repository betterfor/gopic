package uploader

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GiteeConfig
// api: https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoContentsPath
type GiteeConfig struct {
	RepoName       string `json,yaml:"repoName"`
	Branch         string `json,yaml:"branch"`
	Path           string `json,yaml:"path"`
	Email          string `json,yaml:"email"`
	Password       string `json,yaml:"password"`
	ClientID       string `json,yaml:"clientId"`
	ClientSecret   string `json,yaml:"clientSecret"`
	AccessToken    string `json,yaml:"accessToken"`
	EnableTimeFile string `json,yaml:"enableTimeFile"` // use date path
}

func (g *GiteeConfig) Name() string {
	return "gitee图床"
}

// Upload POST https://gitee.com/api/v5/repos/{owner}/{repo}/contents/{path}
func (g *GiteeConfig) Upload(img *Img) error {
	err := g.check()
	if err != nil {
		return err
	}
	if g.EnableTimeFile == "true" {
		now := time.Now()
		g.Path += "/" + now.Format("2006/01/02")
	}
	u := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/contents/%s/%s", g.RepoName, g.Path, img.FileName)
	// make body
	buffer, err := io.ReadAll(img.Reader)
	if err != nil {
		return err
	}
	ret := base64.StdEncoding.EncodeToString(buffer)

	var token string
	if len(g.AccessToken) != 0 {
		token, err = g.getToken()
		if err != nil {
			return err
		}
	} else {
		token = g.AccessToken
	}
	body := giteeRequest{
		AccessToken: token,
		Message:     "upload by gopic",
		Content:     ret, // base64
		//Sha:     fmt.Sprintf("%x", md5.Sum(content)),        // md5
		Branch: g.Branch,
	}

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(body.String()))
	if err != nil {
		return err
	}
	// set header
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse result
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		img.ImgUrl = gjson.GetBytes(bts, "content.download_url").String()
		return nil
	case http.StatusUnauthorized:
		token, err = g.getToken()
		if err != nil {
			return err
		}
		viper.Set("data.gitee.accessToken", token)
		return fmt.Errorf("access token is expired, please try again")
	}
	return errors.New(resp.Status)
}

func (g *GiteeConfig) check() error {
	if len(g.RepoName) == 0 {
		return errors.New("repo name is empty")
	}
	if len(g.Branch) == 0 {
		g.Branch = "master"
	}
	if len(g.Path) == 0 {
		g.Path = "images"
	}
	//if len(g.AccessToken) == 0 {
	//	return errors.New("invalid token")
	//}
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

func (g *GiteeConfig) getToken() (string, error) {
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
