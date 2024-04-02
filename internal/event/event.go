package event

import (
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Init(p *proxy.Proxy) {
	event.Subscribe(p.Event(), 0, loginEvent(p))
}
