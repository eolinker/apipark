package main

import (
	_ "github.com/eolinker/ap-account/plugin"
	_ "github.com/eolinker/apipark/frontend"
	_ "github.com/eolinker/apipark/gateway/apinto"
	_ "github.com/eolinker/apipark/plugins/core"
	_ "github.com/eolinker/apipark/plugins/permit"
	_ "github.com/eolinker/apipark/plugins/publish_flow"
	_ "github.com/eolinker/go-common/cache/cache_redis"
	_ "github.com/eolinker/go-common/log-init"
	_ "github.com/eolinker/go-common/store/store_mysql"
)
