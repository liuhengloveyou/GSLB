package dao

import (
	"database/sql"
)

// SELECT id,host,zone, type, ttl,record,view FROM gslb.dnsinfo where status=1;
type RR struct {
	ID     int            `db:"id"`
	TTL    uint32         `db:"ttl"`
	Host   string         `db:"host"`
	Zone   string         `db:"zone"`
	Type   string         `db:"type"`
	Record sql.NullString `db:"record"`
	View   sql.NullString `db:"view"`
}

type View struct {
	ID     int            `db:"id"`
	Domain sql.NullString `db:"domain"`
	Line   sql.NullString `db:"line"`
	Area   sql.NullString `db:"area"`
	View   sql.NullString `db:"view"`
}

// ip库记录
type IpRecord struct {
	ID        int     `db:"id"`
	Start     string  `db:"ip_start"`
	End       string  `db:"ip_end"`
	Country   string  `db:"country"`
	Province  string  `db:"province"`
	ISP       string  `db:"isp"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}
