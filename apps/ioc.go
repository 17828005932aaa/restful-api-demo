package app

import "restful-api-demo/apps/host"

// 我们能不能把这些注册的信息 都放在公共库里
var (
	Host host.ServiceServer
)
