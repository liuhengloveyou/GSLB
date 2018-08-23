package dao

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/GSLB/common"
)

func TestSelectRulesFromMysql(t *testing.T) {

	rules, e := SelectRulesFromMysql([]string{"0x7c00.net."})

	fmt.Println("TestSelectRulesFromMysql: ", rules, e)
}
