package main

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/liuhengloveyou/HTTPDNS/common"
	gocommon "github.com/liuhengloveyou/go-common"
)

type Value struct {
	t int // 

}

// 解析配置缓存
var CacheTree map[string]*Value

func main() {
	if e := gocommon.LoadJsonConfig("./app.conf", &ServConfig); e != nil {
		panic(e)
	}

	s := &http.Server{
		Addr:           ServConfig.Listen,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("%v GO %v\n", time.Now(), ServConfig.Listen)
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
