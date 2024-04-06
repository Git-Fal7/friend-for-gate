package command

import (
	"strings"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Init(p *proxy.Proxy) {
	p.Command().Register(friendCommand(p))
	p.Command().Register(msgCommand(p))
}

// utils
func replaceAll(str string, replaceMap map[string]string) string {
	for key, value := range replaceMap {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}
