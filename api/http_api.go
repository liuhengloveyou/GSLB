package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	log "github.com/sirupsen/logrus"
)

type D struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Images  string `json:"images" db:"images"`
	AddTime int64  `json:"add_time" db:"add_time"`
}

func InitHttpApi() {
	http.Handle("/d", &D{})
}

func (p *D) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		p.find(w, r)
	default:
		gocommon.HttpErr(w, http.StatusMethodNotAllowed, 0, "")
		return
	}
}

func (p *D) find(w http.ResponseWriter, r *http.Request) {

}
