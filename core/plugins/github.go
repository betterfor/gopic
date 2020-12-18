package plugins

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
	"time"
)

const timeout = time.Second * 10

// github options
type GithubOpts struct {
	RepoName  string `json:"repoName"`  // the name of warehouse, like: betterfor/gopic
	Branch    string `json:"branch"`    // project branch, like: master
	Token     string `json:"token"`     // set github personal access tokens
	Path      string `json:"path"`      // storage path in github, support variables ${time-format}
	CustomURL string `json:"customUrl"` // custom domain name,like: https://xxx.com
}

func (g *GithubOpts) URL() string {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/", g.RepoName)
	if len(g.Path) == 0 {
		return baseURL
	}
	// check if the path contains variables
	for {
		t0 := strings.Index(g.Path, "${")
		t1 := strings.Index(g.Path, "}")
		if t1-t0 > 0 {
			now := time.Now()
			g.Path = strings.Replace(g.Path, g.Path[t0:t1+1], now.Format(g.Path[t0+2:t1]), 1)
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

// request of github creating or updating file contents
type githubRequest struct {
	Message string `json:"message"` // The commit message.
	Content string `json:"content"` // The new file content, using Base64 encoding.
	Sha     string `json:"sha"`     // Required if you are updating a file. The blob SHA of the file being replaced.
	Branch  string `json:"branch"`  // The branch name. Default: the repository’s default branch (usually master)
	//Committer Committer	`json:"committer"` // The person that committed the file. Default: the authenticated user.
}

type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (g *githubRequest) String() string {
	bts, _ := json.Marshal(g)
	return string(bts)
}

// https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#create-or-update-file-contents
// PUT /repos/{owner}/{repo}/contents/{path}
func (g *GithubOpts) Upload(fileName string, content []byte) (string, error) {
	// make body
	now := time.Now()

	ret := base64.StdEncoding.EncodeToString(content)
	body := githubRequest{
		Message: fmt.Sprintf("upload file:%s at %s", fileName, now.String()),
		Content: ret, // base64
		//Sha:     fmt.Sprintf("%x", md5.Sum(content)),        // md5
		Branch: g.Branch,
	}

	var client = &http.Client{
		//Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodPut, g.URL()+fileName, strings.NewReader(body.String()))
	if err != nil {
		return "", err
	}
	// set header
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	// set token
	t := oauth2.Token{
		AccessToken: g.Token,
		TokenType:   "token",
	}
	t.SetAuthHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// parse result
	switch resp.StatusCode {
	case 200:
		fmt.Println("Created!")
		return parseData(resp.Body)
	case 201:
		fmt.Println("Updated!")
		return parseData(resp.Body)
	default:
		return "", errors.New(resp.Status)
	}
}

func parseData(result io.Reader) (string, error) {
	var m map[string]map[string]interface{}
	err := json.NewDecoder(result).Decode(&m)
	if err != nil {
		return "", err
	}
	return m["content"]["download_url"].(string), nil
}
