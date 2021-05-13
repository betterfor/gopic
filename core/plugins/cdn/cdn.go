package cdn

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	Jsdelivr = "jsdelivr"
)

type CDN interface {
	Convert(url string) string
}

type Instance func() CDN

var adapters = make(map[string]Instance)

func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cdn: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cdn: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewCDN(adapterName string) (adapter CDN, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		return nil, errors.WithStack(fmt.Errorf("cdn: unknown adapter name %s ", adapterName))
	}
	adapter = instanceFunc()
	return adapter, nil
}
