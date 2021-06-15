package cloud

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const timeout = time.Second * 30

var client = &http.Client{
	Timeout: timeout,
}

// https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#create-or-update-file-contents
// github options
type GithubOpts struct {
	RepoName string `json:"repoName" yaml:"repoName"` // the name of warehouse, like: betterfor/gopic
	Branch   string `json:"branch" yaml:"branch"`     // project branch, default is master
	Token    string `json:"token" yaml:"token"`       // set github personal access tokens
	Path     string `json:"path" yaml:"path"`         // storage path in github, support variables ${time-format}, default is images
}

func (g *GithubOpts) URL() string {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/", g.RepoName)
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

func (g *GithubOpts) String() string {
	bts, _ := json.Marshal(g)
	return string(bts)
}

func (g *GithubOpts) Unmarshal(opts []byte) *GithubOpts {
	var o GithubOpts
	json.Unmarshal(opts, &o)
	return &o
}

// request of github creating or updating file contents
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

// PUT /repos/{owner}/{repo}/contents/{path}
func (g *GithubOpts) Upload(fileName string, data []byte) (string, error) {
	if err := g.check(); err != nil {
		return "", err
	}
	// make body
	now := time.Now()
	ret := base64.StdEncoding.EncodeToString(data)
	body := githubRequest{
		Message: fmt.Sprintf("upload file:%s at %s", fileName, now.String()),
		Content: ret, // base64
		//Sha:     fmt.Sprintf("%x", md5.Sum(content)),        // md5
		Branch: g.Branch,
	}

	req, err := http.NewRequest(http.MethodPut, g.URL()+fileName, strings.NewReader(body.String()))
	if err != nil {
		return "", err
	}
	// set header
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	// set token
	t := oauth2.Token{AccessToken: g.Token, TokenType: "token"}
	t.SetAuthHeader(req)
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
	default:
		return "", errors.New(resp.Status)
	}
}

func (g *GithubOpts) Parse(str string) string {
	return gjson.Get(str, "content.download_url").String()
}

func (g *GithubOpts) check() error {
	if len(g.RepoName) == 0 {
		return errors.New("repo name is empty")
	}
	if len(g.Token) == 0 {
		return errors.New("invalid token")
	}
	if len(g.Branch) == 0 {
		g.Branch = "master"
	}
	if len(g.Path) == 0 {
		g.Path = "images"
	}
	return nil
}
