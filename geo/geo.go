package geo

import (
	"fmt"

	"../common"
)

type Geo interface {
	FindIP(string) (line, area string)
}

type GeoType func(fn string) (Geo, error)

var geos map[string]GeoType = make(map[string]GeoType)

var defaultGeo Geo

func RegisterGeo(name string, newFunc GeoType) {
	if newFunc == nil {
		panic("Register GEO nil.")
	}

	if _, ok := geos[name]; ok {
		panic("Register GEO duplicate for " + name)
	}

	geos[name] = newFunc
}

func NewGeo(t, fn string) (geo Geo, err error) {

	if newFunc, ok := geos[t]; ok {
		return newFunc(fn)
	}

	return nil, fmt.Errorf("No GEO types " + t)
}

func InitGEO() {
	// 现在只有ipip
	RegisterGeo("ipip", newIpipDB)

	var err error
	if defaultGeo, err = NewGeo(common.ServConfig.GeoFmt, common.ServConfig.GeoDB); err != nil {
		panic(err)
	}
}

//// public interface
func FindIP(ip string) (line, area string) {
	return defaultGeo.FindIP(ip)
}
