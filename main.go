package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/liuhengloveyou/GSLB/api"
	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
	"github.com/liuhengloveyou/GSLB/geo"
	"github.com/liuhengloveyou/GSLB/service"

	gocommon "github.com/liuhengloveyou/go-common"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if e := gocommon.LoadTomlConfig("./app.conf.toml", &ServConfig); e != nil {
		panic(e)
	}

	InitLogger()
	defer Logger.Sync()

	fmt.Println("init mysql...")
	dao.InitDB()

	go loadCaches()

	fmt.Println("init GEO database...")
	geo.InitGEO()

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

	/*
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
	*/

	wg.Wait()
}

func loadCaches() {
	for {
		sleep := ServConfig.CacheTTL

		/*
			s, err := service.LoadGroupCache()
			if err != nil {
				Logger.Error("LoadGroupCache ERR: " + err.Error())
			}
			if err == nil && s < sleep {
				sleep = s
			}
		*/

		s, err := service.LoadRRCache()
		if err != nil {
			Logger.Error("LoadRRCache ERR: " + err.Error())
		}
		if err == nil && s < sleep {
			sleep = s
		}

		Logger.Sync()

		if sleep < 3 {
			sleep = 3
		}
		time.Sleep(time.Second * time.Duration(sleep))
	}
}
