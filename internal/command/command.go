package command

import "go.minekube.com/gate/pkg/edition/java/proxy"

func Init(p *proxy.Proxy) {
	p.Command().Register(friendCommand(p))
}
