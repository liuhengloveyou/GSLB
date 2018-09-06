package common

type RR struct {
	ID     int
	Domain string
	Ttl    uint32
	Type   uint16
	Class  uint16
	Data   string
	Group  string

	Weight        int32
	CurrentWeight int32
}

type Rule struct {
	ID     int
	Domain string
	Line   string
	Area   string
	Zone   string
	Group  string
}

type Group struct {
	ID     int
	Domain string
	Name   string
	Policy int
}
