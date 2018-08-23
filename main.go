package main

import (
	"fmt"
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
	defer Logger.Sync()

	var wg sync.WaitGroup

	if ServConfig.HTTPApiAddr != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%v HTTP %v\n", time.Now(), ServConfig.HTTPApiAddr)
			if err := api.InitHttpApi(ServConfig.HTTPApiAddr); err != nil {
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
