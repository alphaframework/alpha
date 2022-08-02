package main

import (
	"fmt"

	_ "github.com/dartagnanli/alpha/aconfig"
	_ "github.com/dartagnanli/alpha/aerror"
	_ "github.com/dartagnanli/alpha/alog"
	_ "github.com/dartagnanli/alpha/alog/gormwrapper"
	_ "github.com/dartagnanli/alpha/autil"
	_ "github.com/dartagnanli/alpha/autil/ahttp"
	_ "github.com/dartagnanli/alpha/autil/ahttp/request"
	_ "github.com/dartagnanli/alpha/database"
	_ "github.com/dartagnanli/alpha/ginwrapper"
	_ "github.com/dartagnanli/alpha/httpclient"
	_ "github.com/dartagnanli/alpha/httpserver/rsp"
)

func main() {
	fmt.Println("Hello world")
}
