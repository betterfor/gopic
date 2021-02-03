package plugins

import (
	"io/ioutil"
	"testing"
)

func TestURL(t *testing.T) {
	// basic
	opts := &GithubOpts{RepoName: "betterfor/cloudImage", Path: "images"}
	t.Log(opts.URL())

	// use time format
	opts.Path = "images/${2006-01-02}"
	t.Log(opts.URL())

	opts.Path = "images/${2006/01/02}"
	t.Log(opts.URL())

	opts.Path = "images/${2006-01-02T15:04:05Z}"
	t.Log(opts.URL())

	opts.Path = "images/${2006-01-02 15:04:05}"
	t.Log(opts.URL())

	// not use time format
	opts.Path = "images/${2006-01-02"
	t.Log(opts.URL())

	opts.Path = "images/$2006-01-02"
	t.Log(opts.URL())
}

func TestExampleGithub(t *testing.T) {
	opts := &GithubOpts{
		RepoName: "betterfor/cloudImage",
		Branch:   "master",
		Token:    "xxx",
		Path:     "images/test",
	}

	t.Log(opts.URL())

	bts, _ := ioutil.ReadFile("./testdata/helloworld.png")
	results, err := opts.Upload("test1.png", bts)
	if err != nil {
		t.Fatal("upload error:", err)
		return
	}
	t.Log(results)
}
