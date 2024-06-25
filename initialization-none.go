//go:build !init

package main

import (
	_ "github.com/eolinker/apipark/resources/access"
	_ "github.com/eolinker/apipark/resources/permit"
	_ "github.com/eolinker/apipark/resources/plugin"
)

func doCheck() {

}
