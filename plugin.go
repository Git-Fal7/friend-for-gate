package friendforgate

import (
	"github.com/git-fal7/friend-for-gate/internal/plugin"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "Friends",
	Init: plugin.InitPlugin,
}
