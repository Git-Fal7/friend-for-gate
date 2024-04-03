# friend-for-gate

A friend system plugin for the [Gate](https://golang.org/doc/install) proxy

# getting started

add the package into your gate proxy
```
go get https://github.com/git-fal7/friend-for-gate
```

then, append the plugin into your plugins from the main() func
```
func main() {
	proxy.Plugins = append(proxy.Plugins,
		// your plugins
		friendforgate.Plugin,
	)
	gate.Execute()
}
```