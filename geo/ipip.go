package geo

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type IpRecord struct {
	start        uint32
	end          uint32
	Country      string
	Province     string
	City         string
	Organization string
	ISP          string
	Latitude     string
	Longitude    string
}

type IpipDB struct {
	file *os.File

	records []IpRecord
}

var ErrIPv4Format = errors.New("ipv4 format error")
var ErrNotFound = errors.New("not found")

func newIpipDB(fn string) (geo Geo, err error) {
	db := &IpipDB{}

	if err := db.load(fn); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *IpipDB) load(fn string) (err error) {
	if db.file, err = os.Open(fn); err != nil {
		return
	}

	var i int64
	bi := bufio.NewReader(db.file)
	for {
		line, err := bi.ReadSlice('\n')
		if err != nil {
			break
		}
		i = i + 1

		var r *IpRecord
		r, err = parseLine(line)
		if err != nil {
			fmt.Println(i, err)
			continue
			// return fmt.Errorf(err.Error()+" %d", i)
		}
		db.records = append(db.records, *r)
	}

	return nil
}

func (db *IpipDB) find(ip string) (*IpRecord, error) {
	ipv := net.ParseIP(ip)
	if ipv == nil {
		return nil, ErrIPv4Format
	}
	ipiv := binary.BigEndian.Uint32(ipv.To4())

	low := 0
	mid := 0
	high := len(db.records)

	for low <= high {
		mid = int((low + high) / 2)
		r := &db.records[mid]

		start := r.start
		end := r.end

		if ipiv < start {
			high = mid - 1
		} else if ipiv > end {
			low = mid + 1
		} else {
			return r, nil
		}
	}

	return nil, nil
}

func (db *IpipDB) FindIP(ip string) (line, area string) {
	if r, err := db.find(ip); err == nil {
		if r != nil {
			if r.Country == "中国" {
				line = r.ISP
				area = r.Country + r.Province
			} else {
				line = r.Country
				area = r.Province
			}
		}
	}

	return
}

func parseLine(line []byte) (r *IpRecord, err error) {
	fields := strings.Fields(string(line))
	if len(fields) != 15 && len(fields) != 17 {
		return nil, fmt.Errorf("ipip line ERR: %d", len(fields))
	}

	r = &IpRecord{}

	ipv := net.ParseIP(fields[0])
	if ipv == nil || ipv.To4() == nil {
		return nil, ErrIPv4Format
	}
	r.start = binary.BigEndian.Uint32(ipv.To4())

	ipv = net.ParseIP(fields[1])
	if ipv == nil || ipv.To4() == nil {
		return nil, ErrIPv4Format
	}
	r.end = binary.BigEndian.Uint32(ipv.To4())
	r.Country = fields[2]
	r.Province = fields[3]
	r.City = fields[4]
	r.Organization = fields[5]
	r.ISP = fields[6]
	r.Latitude = fields[7]
	r.Longitude = fields[8]

	return r, nil
}
