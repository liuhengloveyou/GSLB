package dao

import (
	"fmt"

	"github.com/liuhengloveyou/GSLB/common"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() {
	var e error

	if db, e = sqlx.Connect("mysql", common.ServConfig.Mysql); e != nil {
		panic(e)
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	if e = db.Ping(); e != nil {
		panic(e)
	}
}

func LoadRRFromMysql() (rr []*common.RR, e error) {
	r := []*RR{}

	sql := "SELECT id, host, zone, ttl, type, record, view, policy FROM ns.rr where online=1"
	common.Logger.Debug("LoadRRFromMysql: " + sql)

	e = db.Select(&r, sql)
	common.Logger.Info(fmt.Sprintf("LoadRRFromMysql end: %v %v", r, e))
	if e != nil {
		return
	}

	for i := 0; i < len(r); i++ {
		t := &common.RR{
			ID:     r[i].ID,
			TTL:    r[i].TTL,
			Domain: r[i].Host + "." + r[i].Zone + ".",
			Type:   r[i].Type,
		}

		if r[i].Record.Valid {
			t.Record = r[i].Record.String
		}

		if r[i].View.Valid {
			t.View = r[i].View.String
		}

		rr = append(rr, t)
	}

	common.Logger.Info(fmt.Sprintf("LoadRRFromMysql ended: %#v %d\n", rr, len(rr)))
	return rr, nil
}

func SelectRRsFromMysql(d []string) (rr []*common.RR, e error) {
	r := []RR{}

	sql := "SELECT * FROM ns.rr where domain in ('" + d[0] + "'"
	for i := 1; i < len(d); i++ {
		sql = sql + ", '" + d[i] + "'"
	}
	sql = sql + ");"
	common.Logger.Debug("SelectRRsFromMysql: " + sql)

	e = db.Select(&r, sql)
	common.Logger.Info(fmt.Sprintf("SelectRRsFromMysql end: %v %v", r, e))
	if e != nil {
		return
	}

	for i := 0; i < len(r); i++ {
		t := &common.RR{
			ID:     r[i].ID,
			Domain: r[i].Host + r[i].Zone,
			TTL:    r[i].TTL,
			Type:   r[i].Type,
		}

		rr = append(rr, t)
	}

	common.Logger.Info(fmt.Sprintf("SelectRRsFromMysql ended: %#v %d\n", rr, len(rr)))
	return rr, nil
}

func SelectViewFromMysql(line, area string) (view *common.View, e error) {
	sql := "SELECT id,isp as line, province as area, view_key as view FROM gslb.viewinfo_key_mapping where isp_name='" + line + "' and province_name='" + area + "'"
	common.Logger.Debug("SelectViewFromMysql: " + sql)

	var rst View

	e = db.Get(&rst, sql)
	common.Logger.Info(fmt.Sprintf("SelectViewFromMysql end: %v %v", rst, e))
	if e != nil {
		return
	}

	view = &common.View{}
	view.ID = rst.ID
	if rst.Domain.Valid {
		view.Domain = rst.Domain.String
	}
	if rst.Line.Valid {
		view.Line = rst.Line.String
	}
	if rst.Area.Valid {
		view.Area = rst.Area.String
	}
	if rst.View.Valid {
		view.View = rst.View.String
	}

	return
}
