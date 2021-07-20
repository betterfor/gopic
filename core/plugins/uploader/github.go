package uploader

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
	"time"
)

// GithubConfig github options
// api: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#create-or-update-file-contents
type GithubConfig struct {
	RepoName       string `json,yaml:"repoName"`       // the name of warehouse, like: betterfor/gopic
	Branch         string `json,yaml:"branch"`         // project branch, default is master
	Token          string `json,yaml:"token"`          // set github personal access tokens
	Path           string `json,yaml:"path"`           // storage path in github, default is images
	CustomUrl      string `json,yaml:"customUrl"`      // convert url to custom url
	EnableTimeFile string `json,yaml:"enableTimeFile"` // use date path
}

func (g *GithubConfig) Name() string {
	return "Github图床"
}

// Upload upload image to github
func (g *GithubConfig) Upload(img *Img) error {
	buf, err := io.ReadAll(img.Reader)
	if err != nil {
		return err
	}
	ret := base64.StdEncoding.EncodeToString(buf)
	body := githubRequest{
		Message: "upload by gopic",
		Content: ret, // base64
		//Sha:     fmt.Sprintf("%x", md5.Sum(content)),        // md5
		Branch: g.Branch,
	}

	// PUT /repos/{owner}/{repo}/contents/{path}
	if g.EnableTimeFile == "true" {
		now := time.Now()
		g.Path += "/" + now.Format("2006/01/02")
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s/%s", g.RepoName, g.Path, img.FileName)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body.String()))
	if err != nil {
		return err
	}
	// set header
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	// set token
	t := oauth2.Token{AccessToken: g.Token, TokenType: "token"}
	t.SetAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse result
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		result, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if len(g.CustomUrl) == 0 {
			img.ImgUrl = gjson.GetBytes(result, "content.download_url").String()
		} else {
			img.ImgUrl = fmt.Sprintf("%s/%s/%s/%s", g.CustomUrl, g.RepoName, g.Path, img.FileName)
		}
		return nil
	default:
		return fmt.Errorf("Server error: %s, please try again ", resp.Status)
	}
}

// githubRequest request of github creating or updating file contents
type githubRequest struct {
	Message string `json:"message"` // The commit message.
	Content string `json:"content"` // The new file content, using Base64 encoding.
	Sha     string `json:"sha"`     // Required if you are updating a file. The blob SHA of the file being replaced.
	Branch  string `json:"branch"`  // The branch name. Default: the repository’s default branch (usually master)
}

func (g *githubRequest) String() string {
	bts, _ := json.Marshal(g)
	return string(bts)
}
