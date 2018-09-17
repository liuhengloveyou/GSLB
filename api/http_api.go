package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "../common"
	"../service"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

const (
	ErrServer   = 1
	ErrClient   = 2
	ErrNotFound = 3
)

type Record struct {
	Type string `json:"type"`
	Host string `json:"host"`
	TTL  uint32 `json:"ttl"`
}

type DomainRecord struct {
	N  string    `json:"n"`
	S  int       `json:"s"`
	Rs []*Record `json:"rs,omitempty"`
}

type Result struct {
	S    int             `json:"s"`
	E    string          `json:"e,omitempty"`
	U    string          `json:"u"`
	V    string          `json:"v,omitempty"`
	Data []*DomainRecord `json:"data"`
}

type D struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Images  string `json:"images" db:"images"`
	AddTime int64  `json:"add_time" db:"add_time"`
}

func InitHttpApi(addr string) error {
	http.Handle("/d", &D{})

	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (p *D) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer Logger.Sync()

	switch r.Method {
	case "GET":
		p.get(w, r)
	default:
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, 0, "")
		return
	}
}

func (p *D) get(w http.ResponseWriter, r *http.Request) {
	var dn []string
	var ip string
	var qk, qv string
	var c int

	r.ParseForm()

	for k, v := range r.Form {
		switch k {
		case "dn":
			dn = v
		case "ip":
			if len(v) > 0 {
				ip = v[0]
			}
		case "c":
			if len(v) > 0 {
				c, _ = strconv.Atoi(v[0])
			}
		default:
			qk = k
			if len(v) > 0 {
				qv = v[0]
			}
		}
	}

	rst := &Result{}

	Logger.Debug("HTTP get", zap.Strings("dn", dn), zap.String("ip", ip), zap.String("qk", qk), zap.String("qv", qv), zap.Int("c", c))

	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	if len(dn) <= 0 {
		Logger.Error("HTTP dn param nil " + ip)
		rst.S = ErrClient
		rst.E = "dn param must not null."
		httpOut(w, http.StatusOK, rst)
		return
	}

	if c < 1 {
		c = 1 // 最少返回一条记录
	}

	// 合法的域名以.结尾
	for i := 0; i < len(dn); i++ {
		if false == strings.HasSuffix(dn[i], ".") {
			dn[i] = dn[i] + "."
		}
	}

	rst.U = ip

	Logger.Info("HTTP get", zap.Strings("dn", dn), zap.String("ip", ip), zap.String("qk", qk), zap.String("qv", qv))

	// 业务定制解析; 要跟业务确定需求 @@@
	if qk != "" {
		service.ResolvLogicRule(qk, qv)
		return // 不再往下处理
	}

	qq := make(map[string]map[string][]RR)
	for _, dnn := range dn {
		qq[dnn] = make(map[string][]RR)
	}

	if err := service.ResolvDomains(ip, c, qq); err != nil {
		Logger.Error("DNS resolv ERR: " + err.Error())
		return
	}

	for k, v := range qq {
		data := &DomainRecord{
			N:  k,
			S:  ErrNotFound,
			Rs: nil,
		}

		// A记录优先于CNAME,
		if _, ok := v["A"]; ok {
			delete(v, "CNAME")
		}

		for t, r := range v {
			for _, r1 := range r {
				record := &Record{
					Type: t,
					Host: r1.Record,
					TTL:  r1.TTL,
				}
				data.S = 0
				data.Rs = append(data.Rs, record)
			}
		}

		rst.Data = append(rst.Data, data)
	}

	Logger.Info(fmt.Sprintf("DNSRootServer OK: %#v", rst))
	httpOut(w, http.StatusOK, rst)

	return
}

func httpOut(w http.ResponseWriter, statCode int, resp *Result) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statCode)

	b, _ := json.Marshal(resp)
	w.Write(b)

	w.(http.Flusher).Flush()
}
