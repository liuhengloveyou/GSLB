package common

type RR struct {
	ID     int
	TTL    uint32
	Type   uint16
	Domain string
	Record string
	View   string

	Weight        int32
	CurrentWeight int32

	Latitude  float64
	Longitude float64
	Distance  float64
}

type View struct {
	ID     int
	Domain string
	Line   string
	Area   string
	View   string
}

type Group struct {
	ID     int
	Domain string
	Name   string
	Policy int
}
