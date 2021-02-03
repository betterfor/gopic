/**
 *Created by XieJian on 2020/12/24 14:08
 *@Desc:
 */
package plugins

import (
	"io/ioutil"
	"testing"
)

func TestSmmsOpts_Upload(t *testing.T) {
	opts := &SmmsOpts{Token: "xxx"}
	bts, _ := ioutil.ReadFile("./testdata/helloworld.png")
	url, err := opts.Upload("test1.png", bts)
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
