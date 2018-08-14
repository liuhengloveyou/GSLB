package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/liuhengloveyou/GSLB/api"
	. "github.com/liuhengloveyou/GSLB/common"
)

type Value struct {
	t int //

}

// 解析配置缓存
var CacheTree map[string]*Value

func main() {
	var wg sync.WaitGroup

	if ServConfig.HTTPApiAddr != "" {
		s := &http.Server{
			Addr:           ServConfig.HTTPApiAddr,
			ReadTimeout:    10 * time.Minute,
			WriteTimeout:   10 * time.Minute,
			MaxHeaderBytes: 1 << 20,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%v HTTP %v\n", time.Now(), ServConfig.HTTPApiAddr)
			if err := s.ListenAndServe(); err != nil {
				panic("HTTPAPI: " + err.Error())
			}
		}()
	}

	if ServConfig.DNSApiAddr != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%v DNS %v\n", time.Now(), ServConfig.DNSApiAddr)
			if err := api.InitDnsApi(ServConfig.DNSApiAddr); err != nil {
				panic("DNSAPI: " + err.Error())
			}
		}()
	}

	wg.Wait()
}
