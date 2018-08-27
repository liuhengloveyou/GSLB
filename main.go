package main

import (
	"fmt"
	"sync"
	//	"time"

	"github.com/liuhengloveyou/GSLB/api"
	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/geo"
	"github.com/liuhengloveyou/GSLB/service"
)

type Value struct {
	t int //

}

// 解析配置缓存
var CacheTree map[string]*Value

func main() {
	defer Logger.Sync()

	go service.LoadRRCache()
	go service.LoadRuleCache()

	fmt.Println("init GEO database...")
	if err := geo.NewGeo(ServConfig.GeoDB); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	if ServConfig.HTTPApiAddr != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("HTTP %v\n", ServConfig.HTTPApiAddr)
			if err := api.InitHttpApi(ServConfig.HTTPApiAddr); err != nil {
				panic("HTTPAPI: " + err.Error())
			}
		}()
	}

	if ServConfig.DNSApiAddr != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("DNS %v\n", ServConfig.DNSApiAddr)
			if err := api.InitDnsApi(ServConfig.DNSApiAddr); err != nil {
				panic("DNSAPI: " + err.Error())
			}
		}()
	}

	wg.Wait()
}
