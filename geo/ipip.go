package geo

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

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

		start := r.Start
		end := r.End

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

func (db *IpipDB) FindIP(ip string) (*IpRecord, error) {
	return db.find(ip)
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
	r.Start = binary.BigEndian.Uint32(ipv.To4())

	ipv = net.ParseIP(fields[1])
	if ipv == nil || ipv.To4() == nil {
		return nil, ErrIPv4Format
	}
	r.End = binary.BigEndian.Uint32(ipv.To4())
	r.Country = fields[2]
	r.Province = fields[3]
	r.ISP = fields[6]
	r.Latitude, _ = strconv.ParseFloat(fields[7], 64)
	r.Longitude, _ = strconv.ParseFloat(fields[8], 64)

	return r, nil
}
