package service

import (
	"fmt"
	"testing"

	. "../common"
)

func TestWRR(t *testing.T) {
	groupCache = make(map[string]map[string][]*RR)

	groupCache["t.com"] = make(map[string][]*RR)
	groupCache["t.com"]["t1"] = []*RR{
		&RR{
			Data:   "1",
			Weight: 1,
		},
		&RR{
			Data:   "2",
			Weight: 2,
		},
		&RR{
			Data:   "3",
			Weight: 3,
		},
		&RR{
			Data:   "4",
			Weight: 4,
		},
	}

	lb := &LBWRR{}

	for i := 0; i < 20; i++ {
		rr := lb.Get("t.com", "t1")
		fmt.Printf(">>>%#v\n", rr)
	}

	lb1 := LBRR{}
	for i := 0; i < 20; i++ {
		fmt.Printf(">>>1 %#v\n", lb1.Get("t.com", "t1"))
	}
}
