package geo

type Geo interface {
	FindIP(string) (line, area string)
}

var defaultGeo Geo

func NewGeo(fn string) (err error) {
	defaultGeo, err = newIpipDB(fn)

	return
}

//// public interface
func FindIP(ip string) (line, area string) {
	return defaultGeo.FindIP(ip)
}
